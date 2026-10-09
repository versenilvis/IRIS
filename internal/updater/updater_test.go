package updater

import (
	"testing"

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
		{"dev", "v1.0.0", "stable", false},
		{"v1.0.0", "dev", "stable", false},
		{"", "v1.0.0", "stable", false},
		{"v1.0.0", "v1.1.0-nightly.8cb1f47", "stable", false},
		{"v1.1.0-nightly.abc", "v1.2.0", "stable", true},
		{"v1.1.0-nightly.abc", "v1.1.0-nightly.def", "nightly", true},
		{"v1.1.0-nightly.abc", "v1.1.0-nightly.abc", "nightly", false},
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
