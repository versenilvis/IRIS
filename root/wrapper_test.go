package root

import (
	"os"
	"path/filepath"
	"testing"
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
