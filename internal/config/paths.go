package config

import (
	"os"
	"path/filepath"
)

func ConfigPath() (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome != "" {
		return filepath.Join(configHome, "iris", "config.toml"), nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "iris", "config.toml"), nil
}

func ThemePath() (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome != "" {
		return filepath.Join(configHome, "iris", "themes"), nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "iris", "themes"), nil
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
