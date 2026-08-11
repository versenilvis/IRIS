package root

import (
	"fmt"
	"os"
	"testing"
)

func TestRelayWatchdogCWDAppliesLatestOfMultipleWrites(t *testing.T) {
	firstDir := t.TempDir()
	secondDir := t.TempDir()
	t.Chdir(t.TempDir())

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fmt.Fprintf(w, "%s\x00%s\x00", firstDir, secondDir); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	relayWatchdogCWD(r)

	if got := currentDir(t); got != wantDir(t, secondDir) {
		t.Fatalf("relayWatchdogCWD left cwd at %q, want the last write %q", got, wantDir(t, secondDir))
	}
}
