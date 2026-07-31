package config

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

func DefaultConfig() *Config {
	return &Config{
		Core: CoreConfig{
			Version:     1,
			Shell:       "",
			ShellLogin:  false,
			Mode:        "last",
			Debug:       false,
			ExpandAlias: true,
			AutoExecute: false,
		},
		UI: UIConfig{
			Style:           "modern",
			GhostText:       true,
			ShowHiddenFiles: false,
			MaxSuggestions:  100,
			MaxHeight:       15,
			MaxWidth:        0, // 0 means no limit, fallback to terminal width
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
		Keybindings: KeybindingsConfig{
			ToggleMode:       "ctrl+r",
			ToggleMenu:       "shift+tab",
			SelectSuggestion: "tab",
			NavigateUp:       "up",
			NavigateDown:     "down",
		},
		Theme: ThemeConfig{
			Border:     lipgloss.Color("#a277ff"),
			Accent:     lipgloss.Color("#61ffca"),
			Muted:      lipgloss.Color("#6d6a7f"),
			Text:       lipgloss.Color("#edecee"),
			TextSel:    lipgloss.Color("#ffffff"),
			Match:      lipgloss.Color("#61ffca"),
			Desc:       lipgloss.Color("#9692a8"),
			DescSel:    lipgloss.Color("#edecee"),
			SelBg:      lipgloss.Color("#3d375e"),
			ScrollInfo: lipgloss.Color("#a277ff"),
			GhostText:  lipgloss.Color("#4B4A4C"),
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
