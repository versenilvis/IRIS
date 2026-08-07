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
	// auto-update loop guard: which version the last auto-install attempt
	// targeted, how many consecutive attempts have been made against it
	// (used to escalate to a confirm prompt, then give up), and any version
	// the user explicitly declined so it's never re-prompted
	AutoUpdateTarget  string `toml:"auto-update-target"`
	AutoUpdateAttempt int    `toml:"auto-update-attempt"`
	DeclinedVersion   string `toml:"declined-version"`
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

	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
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
	if err := enc.Encode(s); err != nil {
		return err
	}

	return nil
}
