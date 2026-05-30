package config

import (
	"encoding"
	"fmt"
	"os"
	"path/filepath"
	"sync"
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
	Version int    `toml:"version"`
	Shell   string `toml:"shell"`
	Mode    string `toml:"mode"`
	Debug   bool   `toml:"debug"`
}

type UIConfig struct {
	GhostText      bool `toml:"ghost-text"`
	MaxSuggestions int  `toml:"max-suggestions"`
	MaxHeight      int  `toml:"max-height"`
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

type Config struct {
	Core    CoreConfig    `toml:"core"`
	UI      UIConfig      `toml:"ui"`
	Git     GitConfig     `toml:"git"`
	Updater UpdaterConfig `toml:"updater"`
}

var (
	activeConfig *Config
	once         sync.Once
)

func Get() *Config {
	once.Do(func() {
		if activeConfig == nil {
			activeConfig = DefaultConfig()
		}
	})
	return activeConfig
}

func Init(cfg *Config) {
	activeConfig = cfg
	once.Do(func() {})
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

	applyEnv(cfg)

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

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
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
