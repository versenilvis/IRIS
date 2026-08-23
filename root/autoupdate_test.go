package root

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/versenilvis/iris/internal/config"
)

// captureStdout redirects os.Stdout for the duration of fn and returns
// everything written to it. writeStdout writes to os.Stdout directly, so
// this is the only way to observe it without a real terminal.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	original := os.Stdout
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = original
	out, _ := io.ReadAll(r)
	return string(out)
}

func TestDecideAutoUpdateActionOff(t *testing.T) {
	st := config.UpdaterState{AutoUpdateTarget: "v1.1.0", AutoUpdateAttempt: 5}
	if got := decideAutoUpdateAction(0, "v1.1.0", st); got != autoUpdateNotifyOnly {
		t.Errorf("mode 0: expected notifyOnly, got %v", got)
	}
}

func TestDecideAutoUpdateActionDeclinedVersionNeverRePrompted(t *testing.T) {
	st := config.UpdaterState{DeclinedVersion: "v1.1.0"}
	if got := decideAutoUpdateAction(1, "v1.1.0", st); got != autoUpdateNotifyOnly {
		t.Errorf("mode 1, declined: expected notifyOnly, got %v", got)
	}
	if got := decideAutoUpdateAction(2, "v1.1.0", st); got != autoUpdateNotifyOnly {
		t.Errorf("mode 2, declined: expected notifyOnly, got %v", got)
	}
}

func TestDecideAutoUpdateActionModeTwoAlwaysConfirms(t *testing.T) {
	cases := []config.UpdaterState{
		{},
		{AutoUpdateTarget: "v1.1.0", AutoUpdateAttempt: 0},
		{AutoUpdateTarget: "v1.1.0", AutoUpdateAttempt: 1},
		{AutoUpdateTarget: "v1.1.0", AutoUpdateAttempt: 3},
	}
	for _, st := range cases {
		if got := decideAutoUpdateAction(2, "v1.1.0", st); got != autoUpdateConfirm {
			t.Errorf("mode 2, state %+v: expected confirm, got %v", st, got)
		}
	}
}

func TestDecideAutoUpdateActionModeOneEscalationLadder(t *testing.T) {
	// a different or absent target: this is a fresh version, try #1
	if got := decideAutoUpdateAction(1, "v1.1.0", config.UpdaterState{}); got != autoUpdateInstallSilently {
		t.Errorf("fresh target: expected installSilently, got %v", got)
	}
	if got := decideAutoUpdateAction(1, "v1.1.0", config.UpdaterState{AutoUpdateTarget: "v1.0.5", AutoUpdateAttempt: 2}); got != autoUpdateInstallSilently {
		t.Errorf("different target: expected installSilently (attempt resets), got %v", got)
	}

	// same target, one attempt already made: this is try #2
	st := config.UpdaterState{AutoUpdateTarget: "v1.1.0", AutoUpdateAttempt: 1}
	if got := decideAutoUpdateAction(1, "v1.1.0", st); got != autoUpdateConfirm {
		t.Errorf("attempt 1 stored: expected confirm, got %v", got)
	}

	// same target, two attempts already made: this is try #3, give up
	st = config.UpdaterState{AutoUpdateTarget: "v1.1.0", AutoUpdateAttempt: 2}
	if got := decideAutoUpdateAction(1, "v1.1.0", st); got != autoUpdateGiveUp {
		t.Errorf("attempt 2 stored: expected giveUp, got %v", got)
	}

	// already given up: stays given up, doesn't loop back to installing
	st = config.UpdaterState{AutoUpdateTarget: "v1.1.0", AutoUpdateAttempt: 3}
	if got := decideAutoUpdateAction(1, "v1.1.0", st); got != autoUpdateGiveUp {
		t.Errorf("attempt 3 stored: expected giveUp (stays given up), got %v", got)
	}
}

func TestDecideAutoUpdateActionNewerVersionResetsLadder(t *testing.T) {
	// gave up on v1.1.0, but v1.2.0 is a different target entirely - fresh start
	st := config.UpdaterState{AutoUpdateTarget: "v1.1.0", AutoUpdateAttempt: 3}
	if got := decideAutoUpdateAction(1, "v1.2.0", st); got != autoUpdateInstallSilently {
		t.Errorf("newer target after giving up on an older one: expected installSilently, got %v", got)
	}
}

func TestPrintAutoUpdateInstalledNotice(t *testing.T) {
	originalVersion := Version
	Version = "v1.0.0"
	t.Cleanup(func() { Version = originalVersion })

	out := captureStdout(t, func() { printAutoUpdateInstalledNotice("v1.1.0") })

	if !strings.Contains(out, "v1.0.0") || !strings.Contains(out, "v1.1.0") {
		t.Errorf("expected the version transition, got %q", out)
	}
	if !strings.Contains(out, "restart your terminal") {
		t.Errorf("expected a restart notice, got %q", out)
	}
}

func TestPrintAutoUpdateConfirmPrompt(t *testing.T) {
	originalVersion := Version
	Version = "v1.0.0"
	t.Cleanup(func() { Version = originalVersion })

	out := captureStdout(t, func() {
		printAutoUpdateConfirmPrompt("v1.1.0", "## Changelog\n### Bug fixes\n* abc1234  fix something\n")
	})

	if !strings.Contains(out, "v1.1.0") {
		t.Errorf("expected the target version, got %q", out)
	}
	if !strings.Contains(out, "fix something") {
		t.Errorf("expected the changelog summary, got %q", out)
	}
	if !strings.Contains(out, "install now?") {
		t.Errorf("expected the confirm prompt text, got %q", out)
	}
}

func TestPrintAutoUpdateGiveUpNotice(t *testing.T) {
	out := captureStdout(t, func() { printAutoUpdateGiveUpNotice("v1.1.0") })

	if !strings.Contains(out, "v1.1.0") {
		t.Errorf("expected the target version, got %q", out)
	}
	if !strings.Contains(out, "iris update") {
		t.Errorf("expected a pointer to manual iris update, got %q", out)
	}
}

func TestDeclineAutoUpdateRecordsState(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("XDG_DATA_HOME", tmpDir+"/.local/share")

	state := config.LoadState()
	state.Updater.AutoUpdateTarget = "v1.1.0"
	state.Updater.AutoUpdateAttempt = 2
	if err := config.SaveState(state); err != nil {
		t.Fatalf("SaveState: %v", err)
	}

	captureStdout(t, func() { declineAutoUpdate("v1.1.0") })

	got := config.LoadState()
	if got.Updater.DeclinedVersion != "v1.1.0" {
		t.Errorf("expected DeclinedVersion v1.1.0, got %q", got.Updater.DeclinedVersion)
	}
	if got.Updater.AutoUpdateTarget != "" || got.Updater.AutoUpdateAttempt != 0 {
		t.Errorf("expected the escalation ladder to be reset, got target=%q attempt=%d",
			got.Updater.AutoUpdateTarget, got.Updater.AutoUpdateAttempt)
	}
}

func TestHandleAutoUpdateConfirmKeyNotArmedIsNoop(t *testing.T) {
	pendingAutoUpdateConfirm.Store(false)

	if handleAutoUpdateConfirmKey('y') {
		t.Error("expected handleAutoUpdateConfirmKey to report unhandled when nothing is pending")
	}
}

func TestHandleAutoUpdateConfirmKeyDeclineFlow(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("XDG_DATA_HOME", tmpDir+"/.local/share")
	t.Cleanup(func() { pendingAutoUpdateConfirm.Store(false) })

	armAutoUpdateConfirm("v1.1.0")

	captureStdout(t, func() {
		if !handleAutoUpdateConfirmKey('n') {
			t.Error("expected 'n' to be reported as handled while a confirm is pending")
		}
	})

	if pendingAutoUpdateConfirm.Load() {
		t.Error("expected the pending confirm flag to clear after declining")
	}
	if got := config.LoadState().Updater.DeclinedVersion; got != "v1.1.0" {
		t.Errorf("expected the decline to be recorded, got %q", got)
	}
}

func TestHandleAutoUpdateConfirmKeySwallowsUnrecognizedBytes(t *testing.T) {
	t.Cleanup(func() { pendingAutoUpdateConfirm.Store(false) })
	armAutoUpdateConfirm("v1.1.0")

	if !handleAutoUpdateConfirmKey('x') {
		t.Error("expected an unrecognized byte to still be reported as consumed while pending")
	}
	if !pendingAutoUpdateConfirm.Load() {
		t.Error("expected the prompt to stay armed until y/n/Esc/Ctrl+C")
	}
}
