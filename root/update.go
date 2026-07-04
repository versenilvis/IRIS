package root

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/versenilvis/iris/config"
)

// updateResult is passed from the async checker to the main loop
type updateResult struct {
	latestVersion string
	hasUpdate     bool
}

// pendingUpdate is set by the background goroutine and consumed once after the first IRIS_CMD_STOP
var pendingUpdate chan updateResult

// FetchLatestVersion hits the GitHub Releases API and returns the latest tag name
func FetchLatestVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	endpoint := os.Getenv("IRIS_UPDATE_URL")
	if endpoint == "" {
		if config.Get().Updater.Channel == "nightly" {
			endpoint = "https://api.github.com/repos/versenilvis/iris/releases"
		} else {
			endpoint = "https://api.github.com/repos/versenilvis/iris/releases/latest"
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if config.Get().Updater.Channel == "nightly" && os.Getenv("IRIS_UPDATE_URL") == "" {
		var releases []struct {
			TagName string `json:"tag_name"`
		}
		if err := json.Unmarshal(body, &releases); err != nil {
			return "", err
		}
		if len(releases) == 0 {
			return "", fmt.Errorf("no releases found")
		}
		return releases[0].TagName, nil
	}

	var result struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if result.TagName == "" {
		return "", fmt.Errorf("no tag_name in response")
	}
	return result.TagName, nil
}

// IsNewer returns true if latest is a newer semantic version than current.
// it supports basic vX.Y.Z formats.
func IsNewer(current, latest string) bool {
	c := strings.TrimPrefix(current, "v")
	l := strings.TrimPrefix(latest, "v")

	// dev builds or empty versions never trigger an update
	if c == "" || c == "dev" || l == "" || l == "dev" {
		return false
	}

	// nightly builds are never shown as stable update targets
	if config.Get().Updater.Channel != "nightly" && strings.Contains(l, "-nightly.") {
		return false
	}

	if c == l {
		return false
	}

	cParts := strings.Split(c, ".")
	lParts := strings.Split(l, ".")

	// compare major.minor.patch
	for i := 0; i < len(cParts) && i < len(lParts); i++ {
		// strip pre-release tags like -beta or -rc for numeric comparison
		cClean := strings.Split(cParts[i], "-")[0]
		lClean := strings.Split(lParts[i], "-")[0]

		cv, _ := strconv.Atoi(cClean)
		lv, _ := strconv.Atoi(lClean)
		if lv > cv {
			return true
		}
		if lv < cv {
			return false
		}
	}

	// if all parts are equal, the one with more parts is newer (e.g. 1.0.1 > 1.0)
	return len(lParts) > len(cParts)
}

// startBackgroundUpdateCheck runs a non-blocking goroutine to check for updates.
// it sends a result on the returned channel exactly once, then closes it
//
// for testing without a real release, set IRIS_MOCK_LATEST_VERSION=v1.99.0
func startBackgroundUpdateCheck() chan updateResult {
	ch := make(chan updateResult, 1)

	if !config.Get().Updater.CheckOnStartup {
		close(ch)
		return ch
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				WriteCrashLog(r)
				restoreTerminal()
				printCrashNotice()
				startRescueShell()
				os.Exit(2)
			}
		}()
		defer close(ch)

		// debug override: skip network entirely, resolve immediately
		if mock := os.Getenv("IRIS_MOCK_LATEST_VERSION"); mock != "" {
			if IsNewer(Version, mock) {
				ch <- updateResult{latestVersion: mock, hasUpdate: true}
			}
			return
		}

		state := config.LoadState()

		// only check once every configured check-interval to avoid hammering the API
		if time.Since(state.Updater.LastCheckTime) < time.Duration(config.Get().Updater.CheckInterval) {
			// already checked recently; still notify if we have a cached pending update
			if state.Updater.SeenVersion != "" && IsNewer(Version, state.Updater.SeenVersion) {
				ch <- updateResult{latestVersion: state.Updater.SeenVersion, hasUpdate: true}
			}
			return
		}

		latest, err := FetchLatestVersion()
		if err != nil {
			// no network or API error: silently do nothing
			return
		}

		// update the last check time regardless of result
		state.Updater.LastCheckTime = time.Now()

		if IsNewer(Version, latest) {
			// only notify if user hasn't already seen this specific version notification
			if state.Updater.SeenVersion != latest {
				ch <- updateResult{latestVersion: latest, hasUpdate: true}
			}
			// save the latest as seen_version so future sessions don't re-notify
			// unless a NEWER version comes out (different tag)
			state.Updater.SeenVersion = latest
		} else {
			// up to date: clear the seen_version flag so the next update triggers a fresh notification
			state.Updater.SeenVersion = ""
		}

		_ = config.SaveState(state)
	}()

	return ch
}

// printUpdateNotice writes the one-time update message to stdout
func printUpdateNotice(latest string) {
	fmt.Printf(
		"\r\033[K\033[33m[IRIS] new version %s → %s available, run \033[1miris update\033[0m\033[33m to upgrade\033[0m\n",
		Version, latest,
	)
}

func init() {
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current Iris version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("iris %s\n", Version)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Iris to the latest release",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("checking for updates (current: %s)...\n", Version)

		latest, err := FetchLatestVersion()
		if err != nil {
			fmt.Printf("\033[31m[IRIS] could not reach update server: %v\033[0m\n", err)
			return
		}

		if Version != "dev" && Version != "" && !IsNewer(Version, latest) {
			fmt.Printf("\033[32m[IRIS] already up to date (%s)\033[0m\n", Version)
			// clear seen_version so the notification doesn't show again
			state := config.LoadState()
			state.Updater.SeenVersion = ""
			_ = config.SaveState(state)
			return
		}

		fmt.Printf("\033[36m[IRIS] updating %s → %s\033[0m\n", Version, latest)

		// download and replace the binary using the install script
		installScript := "https://raw.githubusercontent.com/versenilvis/iris/main/scripts/install.sh"
		fmt.Printf("running: curl -sSL %s | sh\n\n", installScript)

		cmdRun := exec.Command("sh", "-c", "curl -sSL "+installScript+" | sh")
		cmdRun.Stdout = os.Stdout
		cmdRun.Stderr = os.Stderr
		cmdRun.Stdin = os.Stdin
		if err := cmdRun.Run(); err != nil {
			fmt.Printf("\n\033[31m[IRIS] update failed: %v\033[0m\n", err)
			return
		}

		// after a successful update, mark as seen so no more notifications
		state := config.LoadState()
		state.Updater.SeenVersion = ""
		_ = config.SaveState(state)

		fmt.Printf("\n\033[32m[IRIS] restart your terminal to use the new version\033[0m\n")
	},
}
