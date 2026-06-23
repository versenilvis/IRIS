package tests

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/versenilvis/iris/config"
	"github.com/versenilvis/iris/root"
)

func TestIsNewer(t *testing.T) {
	tests := []struct {
		current string
		latest  string
		want    bool
	}{
		{"v1.0.0", "v1.0.1", true},
		{"v1.0.1", "v1.0.0", false},
		{"v1.0.0", "v1.0.0", false},
		{"v1.2.3", "v1.2.4", true},
		{"v1.2.0", "v1.1.9", false},
		{"dev", "v1.0.0", false}, // dev never updates
		{"v1.0.0", "dev", false},
		{"", "v1.0.0", false},
		{"v1.0.0", "v1.1.0-nightly.8cb1f47", false}, // nightly never triggers update
		{"v1.1.0-nightly.abc", "v1.2.0", true},      // but if you are on nightly, you can update to stable
	}

	for _, tt := range tests {
		if got := root.IsNewer(tt.current, tt.latest); got != tt.want {
			t.Errorf("IsNewer(%q, %q) = %v; want %v", tt.current, tt.latest, got, tt.want)
		}
	}
}

func TestUpdateState(t *testing.T) {
	// Use a temporary directory for the state file
	tmpDir, err := os.MkdirTemp("", "iris-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Override home dir for testing
	homeBackup := os.Getenv("HOME")
	err = os.Setenv("HOME", tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Setenv("HOME", homeBackup) }()

	xdgBackup := os.Getenv("XDG_DATA_HOME")
	err = os.Setenv("XDG_DATA_HOME", filepath.Join(tmpDir, ".local", "share"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Setenv("XDG_DATA_HOME", xdgBackup) }()

	state := config.LoadState()
	state.Updater.SeenVersion = "v1.0.0"
	state.Updater.LastCheckTime = time.Unix(123456789, 0)

	err = config.SaveState(state)
	if err != nil {
		t.Fatalf("failed to save state: %v", err)
	}

	loaded := config.LoadState()
	if loaded.Updater.SeenVersion != state.Updater.SeenVersion {
		t.Errorf("Expected SeenVersion %q, got %q", state.Updater.SeenVersion, loaded.Updater.SeenVersion)
	}
	if loaded.Updater.LastCheckTime.Unix() != state.Updater.LastCheckTime.Unix() {
		t.Errorf("Expected LastCheck %v, got %v", state.Updater.LastCheckTime, loaded.Updater.LastCheckTime)
	}
}
