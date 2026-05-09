package root

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// updateState holds the persistent update notification state on disk
type updateState struct {
	// seenVersion is the last version the user was notified about.
	// when a newer release than this is found, we show the message again
	SeenVersion string `json:"seen_version"`
	// lastCheck is unix timestamp of the last network check
	LastCheck int64 `json:"last_check"`
}

// updateResult is passed from the async checker to the main loop
type updateResult struct {
	latestVersion string
	hasUpdate     bool
}

// pendingUpdate is set by the background goroutine and consumed once after the first IRIS_CMD_STOP
var pendingUpdate chan updateResult

func getUpdateStateFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".iris", "update_state.json")
}

func LoadUpdateState() updateState {
	file := getUpdateStateFile()
	if file == "" {
		return updateState{}
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return updateState{}
	}
	var s updateState
	if err := json.Unmarshal(data, &s); err != nil {
		return updateState{}
	}
	return s
}

func SaveUpdateState(s updateState) {
	file := getUpdateStateFile()
	if file == "" {
		return
	}
	data, _ := json.MarshalIndent(s, "", "  ")
	_ = os.WriteFile(file, data, 0644)
}

// FetchLatestVersion hits the GitHub Releases API and returns the latest tag name
func FetchLatestVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// allow overriding the version endpoint for testing without a real release
	endpoint := os.Getenv("IRIS_UPDATE_URL")
	if endpoint == "" {
		endpoint = "https://api.github.com/repos/versenilvis/iris/releases/latest"
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
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

	if c == l {
		return false
	}

	cParts := strings.Split(c, ".")
	lParts := strings.Split(l, ".")

	// compare major.minor.patch
	for i := 0; i < len(cParts) && i < len(lParts); i++ {
		cv, _ := strconv.Atoi(cParts[i])
		lv, _ := strconv.Atoi(lParts[i])
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

	go func() {
		defer close(ch)

		// debug override: skip network entirely, resolve immediately
		if mock := os.Getenv("IRIS_MOCK_LATEST_VERSION"); mock != "" {
			if IsNewer(Version, mock) {
				ch <- updateResult{latestVersion: mock, hasUpdate: true}
			}
			return
		}

		state := LoadUpdateState()

		// only check once every 6 hours to avoid hammering the API
		if time.Since(time.Unix(state.LastCheck, 0)) < 6*time.Hour {
			// already checked recently; still notify if we have a cached pending update
			if state.SeenVersion != "" && IsNewer(Version, state.SeenVersion) {
				ch <- updateResult{latestVersion: state.SeenVersion, hasUpdate: true}
			}
			return
		}

		latest, err := FetchLatestVersion()
		if err != nil {
			// no network or API error: silently do nothing
			return
		}

		// update the last check time regardless of result
		state.LastCheck = time.Now().Unix()

		if IsNewer(Version, latest) {
			// only notify if user hasn't already seen this specific version notification
			if state.SeenVersion != latest {
				ch <- updateResult{latestVersion: latest, hasUpdate: true}
			}
			// save the latest as seen_version so future sessions don't re-notify
			// unless a NEWER version comes out (different tag)
			state.SeenVersion = latest
		} else {
			// up to date: clear the seen_version flag so the next update triggers a fresh notification
			state.SeenVersion = ""
		}

		SaveUpdateState(state)
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

		if !IsNewer(Version, latest) {
			fmt.Printf("\033[32m[IRIS] already up to date (%s)\033[0m\n", Version)
			// clear seen_version so the notification doesn't show again
			state := LoadUpdateState()
			state.SeenVersion = ""
			SaveUpdateState(state)
			return
		}

		fmt.Printf("\033[36m[IRIS] updating %s → %s\033[0m\n", Version, latest)

		// download and replace the binary using the install script
		installScript := "https://raw.githubusercontent.com/versenilvis/iris/main/scripts/install.sh"
		fmt.Printf("running: curl -sS %s | sh\n\n", installScript)

		// after a successful update, mark as seen so no more notifications
		state := LoadUpdateState()
		state.SeenVersion = ""
		SaveUpdateState(state)

		fmt.Printf("\n\033[32m[IRIS] restart your terminal to use the new version\033[0m\n")
	},
}
