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

func TestNullTokenSplit(t *testing.T) {
	tests := []struct {
		name        string // subtest name
		data        string // input data
		atEOF       bool   // is at eof
		wantAdvance int    // expected bytes consumed
		wantToken   string // expected token returned
		wantMore    bool   // wants more data before producing a token
	}{
		{"empty at eof", "", true, 0, "", true},
		{"complete token", "/some/dir\x00rest", false, len("/some/dir\x00"), "/some/dir", false},
		{"partial token, not at eof", "/some/dir", false, 0, "", true},
		{"partial token at eof is returned as final token", "/some/dir", true, len("/some/dir"), "/some/dir", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			advance, token, err := nullTokenSplit([]byte(tt.data), tt.atEOF)
			if err != nil {
				t.Fatalf("nullTokenSplit(%q, %v) returned error: %v", tt.data, tt.atEOF, err)
			}
			if advance != tt.wantAdvance {
				t.Errorf("advance = %d, want %d", advance, tt.wantAdvance)
			}
			if tt.wantMore {
				if token != nil {
					t.Errorf("token = %q, want nil (more data needed)", token)
				}
				return
			}
			if string(token) != tt.wantToken {
				t.Errorf("token = %q, want %q", token, tt.wantToken)
			}
		})
	}
}
