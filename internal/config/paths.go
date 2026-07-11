package config

import (
	"os"
	"path/filepath"
)

func ConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "iris", "config.toml"), nil
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

func CachePath() (string, error) {
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
