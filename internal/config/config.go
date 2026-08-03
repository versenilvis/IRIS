package config

import (
	"encoding"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/BurntSushi/toml"
)

type Duration time.Duration

var (
	_ encoding.TextUnmarshaler = (*Duration)(nil)
	_ encoding.TextMarshaler   = (*Duration)(nil)
)

func (d *Duration) UnmarshalText(text []byte) error {
	dur, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}
	*d = Duration(dur)
	return nil
}

func (d Duration) MarshalText() ([]byte, error) {
	return []byte(time.Duration(d).String()), nil
}

type CoreConfig struct {
	Version     int    `toml:"version"`
	Shell       string `toml:"shell"`
	ShellLogin  bool   `toml:"shell-login"`
	Mode        string `toml:"mode"`
	Debug       bool   `toml:"debug"`
	ExpandAlias bool   `toml:"expand-alias"`
	AutoExecute bool   `toml:"auto-execute"`
}

type UIConfig struct {
	Style           string `toml:"style"`
	GhostText       bool   `toml:"ghost-text"`
	ShowHiddenFiles bool   `toml:"hidden-files"`
	MaxSuggestions  int    `toml:"max-suggestions"`
	MaxHeight       int    `toml:"max-height"`
	MaxWidth        int    `toml:"max-width"`
	NerdFonts       bool   `toml:"nerd-fonts"`
}

type GitConfig struct {
	FilterActiveBranch  bool `toml:"filter-active-branch"`
	DeduplicateBranches bool `toml:"deduplicate-branches"`
}

type UpdaterConfig struct {
	CheckOnStartup bool     `toml:"check-on-startup"`
	Channel        string   `toml:"channel"`
	CheckInterval  Duration `toml:"check-interval"`
}

type KeybindingsConfig struct {
	ToggleMode       string `toml:"toggle-mode"`
	ToggleMenu       string `toml:"toggle-menu"`
	SelectSuggestion string `toml:"select"`
	NavigateUp       string `toml:"navigate-up"`
	NavigateDown     string `toml:"navigate-down"`
}

type SuggestOnEmptyConfig struct {
	Enabled       bool `toml:"enabled"`
	DebounceMS    int  `toml:"debounce_ms"`
	MinIntervalMS int  `toml:"min_interval_ms"`
}

type ProviderConfig struct {
	InheritedFrom    string         `toml:"inherited_from"`
	Endpoint         string         `toml:"endpoint"`
	APIKey           string         `toml:"api_key"`
	APIKeyEnv        string         `toml:"api_key_env"`
	Model            string         `toml:"model"`
	TimeoutMS        int            `toml:"timeout_ms"`
	ExtraRequestBody map[string]any `toml:"extra_request_body"`
}

type AIConfig struct {
	Enabled        bool                      `toml:"enabled"`
	Provider       string                    `toml:"provider"`
	DebounceMS     int                       `toml:"debounce_ms"`
	MinIntervalMS  int                       `toml:"min_interval_ms"`
	Providers      map[string]ProviderConfig `toml:"providers"`
	SuggestOnEmpty SuggestOnEmptyConfig      `toml:"suggest_on_empty"`
}

func (c *AIConfig) GetActiveProvider() (ProviderConfig, bool) {
	if c.Providers == nil {
		return ProviderConfig{}, false
	}
	p, ok := c.Providers[c.Provider]
	return p, ok
}

func (p *ProviderConfig) GetAPIKey() string {
	if p.APIKey != "" {
		return p.APIKey
	}
	if p.APIKeyEnv != "" {
		return os.Getenv(p.APIKeyEnv)
	}
	return ""
}

type Config struct {
	Core        CoreConfig        `toml:"core"`
	UI          UIConfig          `toml:"ui"`
	Git         GitConfig         `toml:"git"`
	Updater     UpdaterConfig     `toml:"updater"`
	AI          AIConfig          `toml:"ai"`
	Keybindings KeybindingsConfig `toml:"keybindings"`
}

var (
	activeConfig atomic.Pointer[Config]
	once         sync.Once
)

func Get() *Config {
	once.Do(func() {
		if activeConfig.Load() == nil {
			activeConfig.Store(DefaultConfig())
		}
	})
	return activeConfig.Load()
}

func Init(cfg *Config) {
	activeConfig.Store(cfg)
	once.Do(func() {})
}

func AutoDetectConfigChange(onReload func(cfg *Config)) {
	cfgPath, err := ConfigPath()
	if err != nil {
		return
	}
	themePath, _ := ThemePath()

	go func() {
		statMod := func(p string) time.Time {
			if info, err := os.Stat(p); err == nil {
				return info.ModTime()
			}
			return time.Time{}
		}
		cfgLast := statMod(cfgPath)
		themeLast := statMod(themePath)
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			changed := false
			cfgMod := statMod(cfgPath)
			if !cfgLast.IsZero() && cfgMod.After(cfgLast) {
				cfgLast = cfgMod
				changed = true
			} else if cfgLast.IsZero() {
				cfgLast = cfgMod
			}
			if themePath != "" {
				themeMod := statMod(themePath)
				if !themeLast.IsZero() && (themeMod.After(themeLast) || themeMod.IsZero()) {
					themeLast = themeMod
					changed = true
				} else if themeLast.IsZero() {
					themeLast = themeMod
				}
			}
			if changed {
				if newCfg, err := Load(); err == nil {
					Init(newCfg)
					if onReload != nil {
						onReload(newCfg)
					}
				}
			}
		}
	}()
}

func Load() (*Config, error) {
	cfg := DefaultConfig()

	path, err := ConfigPath()
	if err == nil {
		if _, statErr := os.Stat(path); statErr == nil {
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				return cfg, fmt.Errorf("config: read %s: %w", path, readErr)
			}
			if _, decodeErr := toml.Decode(string(data), cfg); decodeErr != nil {
				return cfg, fmt.Errorf("config: parse %s: %w", path, decodeErr)
			}
		}
	}

	if themePath, err := ThemePath(); err == nil {
		LoadTheme(themePath)
	}

	applyEnv(cfg)

	// fallback for empty keybindings
	if cfg.Keybindings.ToggleMode == "" {
		cfg.Keybindings.ToggleMode = "ctrl+r"
	}
	if cfg.Keybindings.ToggleMenu == "" {
		cfg.Keybindings.ToggleMenu = "shift+tab"
	}
	if cfg.Keybindings.SelectSuggestion == "" {
		cfg.Keybindings.SelectSuggestion = "tab"
	}
	if cfg.Keybindings.NavigateUp == "" {
		cfg.Keybindings.NavigateUp = "up"
	}
	if cfg.Keybindings.NavigateDown == "" {
		cfg.Keybindings.NavigateDown = "down"
	}

	if err := validate(cfg); err != nil {
		return cfg, fmt.Errorf("config: invalid value: %w", err)
	}

	return cfg, nil
}

func Save(cfg *Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := toml.NewEncoder(file)
	if err := enc.Encode(cfg); err != nil {
		return err
	}

	return nil
}

func validate(cfg *Config) error {
	validModes := map[string]bool{"last": true, "spec": true, "history": true}
	if cfg.Core.Mode != "" && !validModes[cfg.Core.Mode] {
		return fmt.Errorf("core.mode: invalid value %q (want: last|spec|history)", cfg.Core.Mode)
	}

	validShells := map[string]bool{"": true, "bash": true, "zsh": true, "fish": true}
	if !validShells[cfg.Core.Shell] {
		return fmt.Errorf("core.shell: invalid value %q (want: bash|zsh|fish)", cfg.Core.Shell)
	}

	validChannels := map[string]bool{"stable": true, "nightly": true}
	if !validChannels[cfg.Updater.Channel] {
		return fmt.Errorf("updater.channel: invalid value %q (want: stable|nightly)", cfg.Updater.Channel)
	}

	if cfg.UI.MaxSuggestions < 1 || cfg.UI.MaxSuggestions > 500 {
		return fmt.Errorf("ui.max-suggestions: must be between 1 and 500")
	}

	if cfg.UI.MaxHeight < 3 || cfg.UI.MaxHeight > 50 {
		return fmt.Errorf("ui.max-height: must be between 3 and 50")
	}

	return nil
}
