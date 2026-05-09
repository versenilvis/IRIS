package tests

import (
	"os"
	"path/filepath"
	"testing"

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
	defer os.RemoveAll(tmpDir)

	// Override home dir for testing
	homeBackup := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", homeBackup)

	// Ensure .iris directory exists
	_ = os.MkdirAll(filepath.Join(tmpDir, ".iris"), 0755)

	state := root.LoadUpdateState()
	state.SeenVersion = "v1.0.0"
	state.LastCheck = 123456789

	root.SaveUpdateState(state)

	loaded := root.LoadUpdateState()
	if loaded.SeenVersion != state.SeenVersion {
		t.Errorf("Expected SeenVersion %q, got %q", state.SeenVersion, loaded.SeenVersion)
	}
	if loaded.LastCheck != state.LastCheck {
		t.Errorf("Expected LastCheck %d, got %d", state.LastCheck, loaded.LastCheck)
	}
}
