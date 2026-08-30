package config

import (
	"time"
)

func DefaultConfig() *Config {
	return &Config{
		Core: CoreConfig{
			Version:           1,
			Shell:             "",
			ShellLogin:        false,
			Mode:              "last",
			Debug:             false,
			ExpandAlias:       true,
			AutoExecute:       false,
			CobraProbeEnabled: true,
			NavigateClosed:    "history",
		},
		UI: UIConfig{
			Style:           "modern",
			GhostText:       GhostTextOn,
			ShowHiddenFiles: false,
			MaxSuggestions:  100,
			MaxHeight:       6,
			MaxWidth:        Width{}, // unset; the overlay falls back to its own default width
			NerdFonts:       true,
		},
		Git: GitConfig{
			FilterActiveBranch:  true,
			DeduplicateBranches: true,
		},
		Updater: UpdaterConfig{
			CheckOnStartup: true,
			Channel:        "stable",
			CheckInterval:  Duration(24 * time.Hour),
			AutoUpdate:     0,
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
		Zoxide: ZoxideConfig{
			ExtendCd: false,
		},
		Keybindings: KeybindingsConfig{
			ToggleMode:       "ctrl+r",
			ToggleMenu:       "shift+tab",
			SelectSuggestion: "tab",
			NavigateUp:       "up",
			NavigateDown:     "down",
		},
	}
}

func DefaultState() *State {
	return &State{
		LastMode: "spec",
		Updater: UpdaterState{
			LastCheckTime:     time.Time{},
			SeenVersion:       "",
			AutoUpdateTarget:  "",
			AutoUpdateAttempt: 0,
			DeclinedVersion:   "",
		},
	}
}
