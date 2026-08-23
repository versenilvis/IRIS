package root

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/versenilvis/iris/internal/config"
)

type updateResultKind int

const (
	updateResultNotify updateResultKind = iota
	updateResultAutoInstalled
	updateResultConfirm
	updateResultGiveUp
)

// updateResult is passed from the async checker to the main loop
type updateResult struct {
	kind          updateResultKind
	latestVersion string
	notes         string
	hasUpdate     bool
}

// pendingUpdate is set by the background goroutine and consumed once after the first IRIS_CMD_STOP
var pendingUpdate chan updateResult

// ErrRateLimited is returned when the GitHub API responds 403/429.
var ErrRateLimited = errors.New("rate limited by GitHub API")

func newGitHubRequestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func fetchGitHubBody(ctx context.Context, endpoint string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimited
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// FetchLatestRelease hits the GitHub Releases API and returns the latest release
func FetchLatestRelease() (Release, error) {
	ctx, cancel := newGitHubRequestContext()
	defer cancel()

	endpoint := os.Getenv("IRIS_UPDATE_URL")
	if endpoint == "" {
		if config.Get().Updater.Channel == "nightly" {
			endpoint = "https://api.github.com/repos/versenilvis/iris/releases"
		} else {
			endpoint = "https://api.github.com/repos/versenilvis/iris/releases/latest"
		}
	}

	body, err := fetchGitHubBody(ctx, endpoint)
	if err != nil {
		return Release{}, err
	}

	if config.Get().Updater.Channel == "nightly" && os.Getenv("IRIS_UPDATE_URL") == "" {
		var releases []Release
		if err := json.Unmarshal(body, &releases); err != nil {
			return Release{}, err
		}
		if len(releases) == 0 {
			return Release{}, fmt.Errorf("no releases found")
		}
		return releases[0], nil
	}

	var release Release
	if err := json.Unmarshal(body, &release); err != nil {
		return Release{}, err
	}
	if release.TagName == "" {
		return Release{}, fmt.Errorf("no tag_name in response")
	}
	return release, nil
}

// FetchLatestVersion hits the GitHub Releases API and returns the latest tag name
func FetchLatestVersion() (string, error) {
	release, err := FetchLatestRelease()
	if err != nil {
		return "", err
	}
	return release.TagName, nil
}

// IsNewer returns true if latest is a newer semantic version than current.
// it supports basic vX.Y.Z formats.
func IsNewer(current, latest string) bool {
	c := strings.TrimPrefix(current, "v")
	l := strings.TrimPrefix(latest, "v")
	channel := config.Get().Updater.Channel
	// dev builds or empty versions never trigger an update
	if c == "" || c == "dev" || l == "" || l == "dev" {
		return false
	}

	// nightly builds are never shown as stable update targets
	if channel != "nightly" && strings.Contains(l, "-nightly.") {
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

	if channel == "nightly" && strings.Contains(l, "-nightly.") && c != l {
		return true
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

		release, err := FetchLatestRelease()
		if err != nil {
			// no network or API error: silently do nothing
			return
		}
		latest := release.TagName

		// update the last check time regardless of result
		state.Updater.LastCheckTime = time.Now()

		if IsNewer(Version, latest) {
			mode := config.Get().Updater.AutoUpdate
			switch decideAutoUpdateAction(mode, latest, state.Updater) {
			case autoUpdateNotifyOnly:
				// only notify if user hasn't already seen this specific version notification
				if state.Updater.SeenVersion != latest {
					ch <- updateResult{kind: updateResultNotify, latestVersion: latest, notes: release.Body, hasUpdate: true}
				}
				// save the latest as seen_version so future sessions don't re-notify
				// unless a NEWER version comes out (different tag)
				state.Updater.SeenVersion = latest
				_ = config.SaveState(state)
				return

			case autoUpdateInstallSilently:
				state.Updater.AutoUpdateTarget = latest
				state.Updater.AutoUpdateAttempt = 1
				// write before installing so a crash still counts as an attempt
				_ = config.SaveState(state)
				if _, installErr := performUpdate(latest, false); installErr == nil {
					state.Updater.AutoUpdateTarget = ""
					state.Updater.AutoUpdateAttempt = 0
					state.Updater.SeenVersion = ""
					_ = config.SaveState(state)
					ch <- updateResult{kind: updateResultAutoInstalled, latestVersion: latest, notes: release.Body, hasUpdate: true}
				}
				// on failure: state already recorded attempt 1 for this
				// target, so the next check escalates to a confirm prompt
				return

			case autoUpdateConfirm:
				nextAttempt := 1
				if state.Updater.AutoUpdateTarget == latest {
					nextAttempt = state.Updater.AutoUpdateAttempt + 1
				}
				state.Updater.AutoUpdateTarget = latest
				state.Updater.AutoUpdateAttempt = nextAttempt
				_ = config.SaveState(state)
				ch <- updateResult{kind: updateResultConfirm, latestVersion: latest, notes: release.Body, hasUpdate: true}
				return

			case autoUpdateGiveUp:
				// only announce the exact transition into giving up, not
				// every cycle after
				announce := state.Updater.AutoUpdateAttempt == 2
				state.Updater.AutoUpdateAttempt = 3
				_ = config.SaveState(state)
				if announce {
					ch <- updateResult{kind: updateResultGiveUp, latestVersion: latest, hasUpdate: true}
				}
				return
			}
		}

		// up to date, or nothing left to do: clear update-related state so
		// the next detected version starts a fresh notify/attempt cycle
		state.Updater.SeenVersion = ""
		state.Updater.AutoUpdateTarget = ""
		state.Updater.AutoUpdateAttempt = 0
		_ = config.SaveState(state)
	}()

	return ch
}

func changelogSummaryLines(body string, max int) []string {
	var lines []string
	for line := range strings.SplitSeq(body, "\n") {
		line = strings.TrimRight(line, "\r")
		if !strings.HasPrefix(line, "* ") {
			continue
		}
		entry := strings.TrimPrefix(line, "* ")
		_, msg, ok := strings.Cut(entry, "  ")
		if !ok {
			msg = entry
		}
		lines = append(lines, msg)
		if len(lines) >= max {
			break
		}
	}
	return lines
}

func printUpdateNotice(latest, notes string) {
	var b strings.Builder
	fmt.Fprintf(&b,
		"\r\033[K\033[33m[IRIS] new version %s → %s available, run \033[1miris update\033[0m\033[33m to upgrade\033[0m\n",
		Version, latest,
	)
	for _, line := range changelogSummaryLines(notes, 2) {
		fmt.Fprintf(&b, "\033[33m  - %s\033[0m\n", line)
	}
	writeStdout([]byte(b.String()))
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

		runningPrefix := ""
		if config.Get().Updater.Channel == "nightly" {
			runningPrefix = fmt.Sprintf("IRIS_RELEASE_TAG=%s ", latest)
		}
		fmt.Printf("running: %scurl -sSL %s | sh\n\n", runningPrefix, resolveInstallScriptURL())

		if _, err := performUpdate(latest, true); err != nil {
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
