package root

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"

	"github.com/versenilvis/iris/internal/config"
)

const installScriptURL = "https://raw.githubusercontent.com/versenilvis/iris/main/scripts/install.sh"

// resolveInstallScriptURL returns IRIS_INSTALL_URL when set (used by
// just debug-autoupdate to point at a local mock instead of the real
// script), falling back to the real install script otherwise.
func resolveInstallScriptURL() string {
	if url := os.Getenv("IRIS_INSTALL_URL"); url != "" {
		return url
	}
	return installScriptURL
}

// performUpdate downloads and installs the given version via the install
// script. Interactive mode (iris update) streams output to the real
// terminal as before, with no timeout - the user is watching and can
// Ctrl+C. Non-interactive mode (the background auto-updater) captures
// output instead, since the wrapper's terminal is in raw mode for the
// whole session and streaming installer output there would corrupt the
// display, and bounds it with a timeout so a hung install can't block
// future checks indefinitely.
func performUpdate(latest string, interactive bool) (output string, err error) {
	ctx := context.Background()
	if !interactive {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
	}

	cmdRun := exec.CommandContext(ctx, "sh", "-c", "curl -sSL "+resolveInstallScriptURL()+" | sh")
	if config.Get().Updater.Channel == "nightly" {
		cmdRun.Env = append(os.Environ(), "IRIS_RELEASE_TAG="+latest)
	}

	if interactive {
		cmdRun.Stdout = os.Stdout
		cmdRun.Stderr = os.Stderr
		cmdRun.Stdin = os.Stdin
		return "", cmdRun.Run()
	}

	out, runErr := cmdRun.CombinedOutput()
	return string(out), runErr
}

type autoUpdateAction int

const (
	autoUpdateNotifyOnly autoUpdateAction = iota
	autoUpdateInstallSilently
	autoUpdateConfirm
	autoUpdateGiveUp
)

// decideAutoUpdateAction is a pure function over the current updater state,
// so the attempt/escalate/give-up ladder can be unit tested without
// touching the network or filesystem.
//
// mode 0 (off) and a version the user already declined both fall back to a
// plain notify. mode 2 (always confirm) never auto-installs silently and
// never gives up - it just keeps asking. mode 1 (auto) escalates per
// target version: attempt 1 installs silently, attempt 2 asks first,
// attempt 3+ gives up and tells the user to run `iris update` themselves.
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

// pendingAutoUpdateConfirm/-Version coordinate the confirm prompt between
// the background checker goroutine (which prints the prompt) and the
// keystroke-reading goroutine (which intercepts the y/n response). Neither
// side blocks on the other.
var (
	pendingAutoUpdateConfirm atomic.Bool
	pendingAutoUpdateVersion atomic.Value
)

func armAutoUpdateConfirm(version string) {
	pendingAutoUpdateVersion.Store(version)
	pendingAutoUpdateConfirm.Store(true)
}

// handleAutoUpdateConfirmKey consumes one byte while a confirm prompt is
// pending. It reports whether the byte was consumed - callers must not
// fall through to normal key handling when true.
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

	if _, err := performUpdate(version, false); err != nil {
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
