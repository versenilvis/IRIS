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
			GhostText:      true,
			MaxSuggestions: 100,
			MaxHeight:      15,
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
