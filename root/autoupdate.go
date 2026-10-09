package root

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/versenilvis/iris/internal/config"
	"github.com/versenilvis/iris/internal/updater"
)

type autoUpdateAction int

const (
	autoUpdateNotifyOnly autoUpdateAction = iota
	autoUpdateInstallSilently
	autoUpdateConfirm
	autoUpdateGiveUp
)

// mode 1 escalates per target version: attempt 1 installs silently, attempt
// 2 asks first, attempt 3+ gives up. mode 2 always asks, never gives up
func decideAutoUpdateAction(mode int, latest string, st config.UpdaterState) autoUpdateAction {
	if mode == 0 {
		return autoUpdateNotifyOnly
	}
	if latest == st.DeclinedVersion {
		return autoUpdateNotifyOnly
	}
	if mode == 2 {
		return autoUpdateConfirm
	}

	attempt := st.AutoUpdateAttempt
	if st.AutoUpdateTarget != latest {
		attempt = 0
	}
	switch attempt {
	case 0:
		return autoUpdateInstallSilently
	case 1:
		return autoUpdateConfirm
	default:
		return autoUpdateGiveUp
	}
}

var (
	pendingAutoUpdateConfirm atomic.Bool
	pendingAutoUpdateVersion atomic.Value
)

func armAutoUpdateConfirm(version string) {
	pendingAutoUpdateVersion.Store(version)
	pendingAutoUpdateConfirm.Store(true)
}

// reports whether the byte was consumed - callers must not fall through
// to normal key handling when true
func handleAutoUpdateConfirmKey(b byte) bool {
	if !pendingAutoUpdateConfirm.Load() {
		return false
	}

	switch b {
	case 'y', 'Y':
		pendingAutoUpdateConfirm.Store(false)
		version, _ := pendingAutoUpdateVersion.Load().(string)
		go runConfirmedAutoUpdate(version)
	case 'n', 'N', 0x1b, 0x03: // n, N, Esc, Ctrl+C
		pendingAutoUpdateConfirm.Store(false)
		version, _ := pendingAutoUpdateVersion.Load().(string)
		declineAutoUpdate(version)
	}
	return true
}

func runConfirmedAutoUpdate(version string) {
	if version == "" {
		return
	}
	writeStdout([]byte("\r\033[K\033[36m[IRIS] updating...\033[0m\n"))

	if _, err := updater.PerformUpdate(version, false); err != nil {
		writeStdout(fmt.Appendf(nil, "\033[31m[IRIS] update failed: %v\033[0m\n", err))
		return
	}

	state := config.LoadState()
	state.Updater.AutoUpdateTarget = ""
	state.Updater.AutoUpdateAttempt = 0
	state.Updater.SeenVersion = ""
	_ = config.SaveState(state)

	writeStdout(fmt.Appendf(nil,
		"\033[32m[IRIS] updated to %s, restart your terminal to use it\033[0m\n", version,
	))
}

func declineAutoUpdate(version string) {
	if version == "" {
		return
	}
	state := config.LoadState()
	state.Updater.DeclinedVersion = version
	state.Updater.AutoUpdateTarget = ""
	state.Updater.AutoUpdateAttempt = 0
	_ = config.SaveState(state)

	writeStdout([]byte("\r\033[K\033[33m[IRIS] update declined, run `iris update` any time to install manually\033[0m\n"))
}

func printAutoUpdateInstalledNotice(latest string) {
	writeStdout(fmt.Appendf(nil,
		"\r\033[K\033[32m[IRIS] auto-updated %s → %s\033[0m\n\033[32m[IRIS] restart your terminal to use the new version\033[0m\n",
		Version, latest,
	))
}

func printAutoUpdateConfirmPrompt(latest, notes string) {
	var b strings.Builder
	fmt.Fprintf(&b, "\r\033[K\033[33m[IRIS] new version %s → %s available\033[0m\n", Version, latest)
	for _, line := range changelogSummaryLines(notes, 2) {
		fmt.Fprintf(&b, "\033[33m  - %s\033[0m\n", line)
	}
	fmt.Fprint(&b, "\033[33minstall now? [y/N] \033[0m")
	writeStdout([]byte(b.String()))
}

func printAutoUpdateGiveUpNotice(latest string) {
	writeStdout(fmt.Appendf(nil,
		"\r\033[K\033[33m[IRIS] could not auto-update to %s after repeated attempts, run \033[1miris update\033[0m\033[33m to install manually\033[0m\n",
		latest,
	))
}
