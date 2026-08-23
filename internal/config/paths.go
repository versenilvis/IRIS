package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// ConfigDirEnv overrides where config.toml and theme.toml are read from.
// It is an environment variable rather than a plain flag because the shell rc
// starts the real session with a bare `exec iris`: arguments are lost at that
// boundary, the environment is not.
const ConfigDirEnv = "IRIS_CONFIG_DIR"

// ErrConfigDir reports an unusable IRIS_CONFIG_DIR. An explicit override that
// points nowhere is a typo or a broken generated path, and silently falling
// back to the default location turns that into "my theme stopped applying".
type ErrConfigDir struct {
	Dir    string
	Reason string
	// Source names what set Dir, so the message points at the flag the user
	// actually typed rather than always blaming the environment variable.
	Source string
}

func (e *ErrConfigDir) Error() string {
	source := e.Source
	if source == "" {
		source = ConfigDirEnv
	}
	return fmt.Sprintf("%s %q %s", source, e.Dir, e.Reason)
}

// configDir resolves the directory holding config.toml and theme.toml.
// overridden reports whether it came from ConfigDirEnv rather than XDG.
func configDir() (dir string, overridden bool, err error) {
	if custom := os.Getenv(ConfigDirEnv); custom != "" {
		if !filepath.IsAbs(custom) {
			abs, absErr := filepath.Abs(custom)
			if absErr != nil {
				return "", true, &ErrConfigDir{Dir: custom, Reason: "cannot be resolved to an absolute path"}
			}
			custom = abs
		}
		info, statErr := os.Stat(custom)
		if statErr != nil {
			return "", true, &ErrConfigDir{Dir: custom, Reason: "does not exist"}
		}
		if !info.IsDir() {
			return "", true, &ErrConfigDir{Dir: custom, Reason: "is not a directory"}
		}
		return custom, true, nil
	}

	if configHome := os.Getenv("XDG_CONFIG_HOME"); configHome != "" {
		return filepath.Join(configHome, "iris"), false, nil
	}
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", false, homeErr
	}
	return filepath.Join(home, ".config", "iris"), false, nil
}

// ConfigDir returns the directory config.toml and theme.toml live in.
func ConfigDir() (string, error) {
	dir, _, err := configDir()
	return dir, err
}

func ConfigPath() (string, error) {
	dir, _, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.toml"), nil
}

func ThemePath() (string, error) {
	dir, _, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "theme.toml"), nil
}

func StatePath() (string, error) {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dataHome = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(dataHome, "iris", "state.toml"), nil
}

func HistoryDBPath() (string, error) {
	statePath, err := StatePath()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(statePath), "history.db"), nil
}

func AtuinDBPath() (string, error) {
	if cfg := Get(); cfg != nil && cfg.Core.AtuinDBPath != "" {
		return cfg.Core.AtuinDBPath, nil
	}
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dataHome = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(dataHome, "atuin", "history.db"), nil
}

func CachePath() (string, error) {
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome != "" {
		return filepath.Join(cacheHome, "iris"), nil
	}
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "iris"), nil
}

func CrashDir() (string, error) {
	cache, err := CachePath()
	if err != nil {
		return "", err
	}
	return filepath.Join(cache, "crashes"), nil
}
