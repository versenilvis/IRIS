package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

type UpdaterState struct {
	LastCheckTime time.Time `toml:"last-check-time"`
	SeenVersion   string    `toml:"seen-version"`
}

type State struct {
	LastMode string       `toml:"last-mode"`
	Updater  UpdaterState `toml:"updater"`
}

func LoadState() *State {
	s := DefaultState()

	path, err := StatePath()
	if err != nil {
		return s
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return s
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return s
	}

	if _, err := toml.Decode(string(data), s); err != nil {
		return s
	}

	return s
}

func SaveState(s *State) error {
	path, err := StatePath()
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
	if err := enc.Encode(s); err != nil {
		return err
	}

	return nil
}
