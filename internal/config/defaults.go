package config

import "time"

func DefaultConfig() *Config {
	return &Config{
		Core: CoreConfig{
			Version: 1,
			Shell:   "",
			Mode:    "last",
			Debug:   false,
		},
		UI: UIConfig{
			Style:          "modern",
			GhostText:      true,
			MaxSuggestions: 100,
			MaxHeight:      15,
			NerdFonts:      true,
		},
		Git: GitConfig{
			FilterActiveBranch:  true,
			DeduplicateBranches: true,
		},
		Updater: UpdaterConfig{
			CheckOnStartup: true,
			Channel:        "stable",
			CheckInterval:  Duration(24 * time.Hour),
		},
		AI: AIConfig{
			Enabled:       false,
			Provider:      "",
			DebounceMS:    500,
			MinIntervalMS: 1000,
			Providers:     nil,
			SuggestOnEmpty: SuggestOnEmptyConfig{
				Enabled:       false,
				DebounceMS:    800,
				MinIntervalMS: 5000,
			},
		},
	}
}

func DefaultState() *State {
	return &State{
		LastMode: "spec",
		Updater: UpdaterState{
			LastCheckTime: time.Time{},
			SeenVersion:   "",
		},
	}
}
