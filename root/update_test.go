package root

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/versenilvis/iris/internal/config"
)

func TestIsNewer(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() {
		config.Init(originalConfig)
	})

	tests := []struct {
		current string
		latest  string
		channel string
		want    bool
	}{
		{"v1.0.0", "v1.0.1", "stable", true},
		{"v1.0.1", "v1.0.0", "stable", false},
		{"v1.0.0", "v1.0.0", "stable", false},
		{"v1.2.3", "v1.2.4", "stable", true},
		{"v1.2.0", "v1.1.9", "stable", false},
		{"dev", "v1.0.0", "stable", false}, // dev never updates
		{"v1.0.0", "dev", "stable", false},
		{"", "v1.0.0", "stable", false},
		{"v1.0.0", "v1.1.0-nightly.8cb1f47", "stable", false},          // nightly never triggers update
		{"v1.1.0-nightly.abc", "v1.2.0", "stable", true},               // but if you are on nightly, you can update to stable
		{"v1.1.0-nightly.abc", "v1.1.0-nightly.def", "nightly", true},  // nightly can update to newer nightly
		{"v1.1.0-nightly.abc", "v1.1.0-nightly.abc", "nightly", false}, // same nightly is not newer
	}

	for _, tt := range tests {
		t.Run(tt.current+"_"+tt.latest+"_"+tt.channel, func(t *testing.T) {
			cfg := config.DefaultConfig()
			cfg.Updater.Channel = tt.channel
			config.Init(cfg)

			if got := IsNewer(tt.current, tt.latest); got != tt.want {
				t.Errorf("IsNewer(%q, %q, %q) = %v; want %v", tt.current, tt.latest, tt.channel, got, tt.want)
			}
		})
	}
}

func TestUpdateState(t *testing.T) {
	// Use a temporary directory for the state file
	tmpDir, err := os.MkdirTemp("", "iris-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	t.Setenv("HOME", tmpDir)
	t.Setenv("XDG_DATA_HOME", filepath.Join(tmpDir, ".local", "share"))

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
