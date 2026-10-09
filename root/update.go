package root

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/versenilvis/iris/internal/config"
	"github.com/versenilvis/iris/internal/updater"
)

type updateResultKind int

const (
	updateResultNotify updateResultKind = iota
	updateResultAutoInstalled
	updateResultConfirm
	updateResultGiveUp
)

type updateResult struct {
	kind          updateResultKind
	latestVersion string
	notes         string
	hasUpdate     bool
}

var pendingUpdate chan updateResult

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

		if mock := os.Getenv("IRIS_MOCK_LATEST_VERSION"); mock != "" {
			if updater.IsNewer(Version, mock) {
				ch <- updateResult{latestVersion: mock, hasUpdate: true}
			}
			return
		}

		state := config.LoadState()

		if time.Since(state.Updater.LastCheckTime) < time.Duration(config.Get().Updater.CheckInterval) {
			if state.Updater.SeenVersion != "" && updater.IsNewer(Version, state.Updater.SeenVersion) {
				ch <- updateResult{latestVersion: state.Updater.SeenVersion, hasUpdate: true}
			}
			return
		}

		release, err := updater.FetchLatestRelease()
		if err != nil {
			return
		}
		latest := release.TagName

		state.Updater.LastCheckTime = time.Now()

		if updater.IsNewer(Version, latest) {
			mode := config.Get().Updater.AutoUpdate
			switch decideAutoUpdateAction(mode, latest, state.Updater) {
			case autoUpdateNotifyOnly:
				if state.Updater.SeenVersion != latest {
					ch <- updateResult{kind: updateResultNotify, latestVersion: latest, notes: release.Body, hasUpdate: true}
				}
				state.Updater.SeenVersion = latest
				_ = config.SaveState(state)
				return

			case autoUpdateInstallSilently:
				state.Updater.AutoUpdateTarget = latest
				state.Updater.AutoUpdateAttempt = 1
				_ = config.SaveState(state)
				if _, installErr := updater.PerformUpdate(latest, false); installErr == nil {
					state.Updater.AutoUpdateTarget = ""
					state.Updater.AutoUpdateAttempt = 0
					state.Updater.SeenVersion = ""
					_ = config.SaveState(state)
					ch <- updateResult{kind: updateResultAutoInstalled, latestVersion: latest, notes: release.Body, hasUpdate: true}
				}
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
				announce := state.Updater.AutoUpdateAttempt == 2
				state.Updater.AutoUpdateAttempt = 3
				_ = config.SaveState(state)
				if announce {
					ch <- updateResult{kind: updateResultGiveUp, latestVersion: latest, hasUpdate: true}
				}
				return
			}
		}

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
