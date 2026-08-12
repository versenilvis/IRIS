package root

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/term"
)

func wantDir(t *testing.T, path string) string {
	t.Helper()
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatal(err)
	}
	return resolved
}

func currentDir(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return wantDir(t, cwd)
}

func TestRestoreStdoutTTYRedirectsFileStdoutToTerminal(t *testing.T) {
	tty, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	if err != nil {
		t.Skip("no controlling terminal available")
	}
	_ = tty.Close()

	// stand in for the temp file powerlevel10k's instant prompt leaves on stdout
	redirected, err := os.CreateTemp(t.TempDir(), "instant-prompt-output")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = redirected.Close() }()

	original := os.Stdout
	os.Stdout = redirected
	t.Cleanup(func() { os.Stdout = original })

	restoreStdoutTTY()

	if os.Stdout == redirected {
		t.Fatal("restoreStdoutTTY left stdout pointing at the redirected file")
	}
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		t.Fatal("restoreStdoutTTY did not leave stdout on a terminal")
	}
	_ = os.Stdout.Close()
}

func TestRestoreStdoutTTYKeepsExistingTerminal(t *testing.T) {
	tty, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	if err != nil {
		t.Skip("no controlling terminal available")
	}

	original := os.Stdout
	os.Stdout = tty
	t.Cleanup(func() {
		os.Stdout = original
		_ = tty.Close()
	})

	restoreStdoutTTY()

	if os.Stdout != tty {
		t.Fatal("restoreStdoutTTY reopened stdout that was already a terminal")
	}
}

func TestSyncProcessCWDFollowsShell(t *testing.T) {
	shellDir := t.TempDir()
	t.Chdir(t.TempDir())

	syncProcessCWD(shellDir)

	if got := currentDir(t); got != wantDir(t, shellDir) {
		t.Fatalf("os.Getwd() = %q, want %q", got, wantDir(t, shellDir))
	}
}

func TestSyncProcessCWDKeepsDirectoryOnBadPath(t *testing.T) {
	launcherDir := t.TempDir()
	t.Chdir(launcherDir)

	for _, cwd := range []string{"", "relative/path", filepath.Join(launcherDir, "does-not-exist")} {
		syncProcessCWD(cwd)

		if got := currentDir(t); got != wantDir(t, launcherDir) {
			t.Fatalf("syncProcessCWD(%q) moved the process to %q, want %q", cwd, got, wantDir(t, launcherDir))
		}
	}
}
