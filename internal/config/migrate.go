package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type legacyState struct {
	Mode string `json:"mode"`
}

type legacyUpdateState struct {
	SeenVersion string `json:"seen_version"`
	LastCheck   int64  `json:"last_check"`
}

func MigrateFromLegacyJSON() error {
	statePath, err := StatePath()
	if err != nil {
		return err
	}

	if _, statErr := os.Stat(statePath); statErr == nil {
		return nil
	}

	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return homeErr
	}

	legacyDir := filepath.Join(home, ".iris")
	legacyStatePath := filepath.Join(legacyDir, "state.json")
	legacyUpdatePath := filepath.Join(legacyDir, "update_state.json")

	hasLegacyState := false
	if _, err := os.Stat(legacyStatePath); err == nil {
		hasLegacyState = true
	}
	hasLegacyUpdate := false
	if _, err := os.Stat(legacyUpdatePath); err == nil {
		hasLegacyUpdate = true
	}

	if !hasLegacyState && !hasLegacyUpdate {
		return nil
	}

	state := DefaultState()

	if hasLegacyState {
		data, err := os.ReadFile(legacyStatePath)
		if err == nil {
			var ls legacyState
			if err := json.Unmarshal(data, &ls); err == nil {
				if ls.Mode == "history" || ls.Mode == "spec" {
					state.LastMode = ls.Mode
				}
			}
		}
	}

	if hasLegacyUpdate {
		data, err := os.ReadFile(legacyUpdatePath)
		if err == nil {
			var lu legacyUpdateState
			if err := json.Unmarshal(data, &lu); err == nil {
				state.Updater.SeenVersion = lu.SeenVersion
				if lu.LastCheck > 0 {
					state.Updater.LastCheckTime = time.Unix(lu.LastCheck, 0)
				}
			}
		}
	}

	if err := SaveState(state); err != nil {
		return err
	}

	if hasLegacyState {
		_ = os.Rename(legacyStatePath, legacyStatePath+".bak")
	}
	if hasLegacyUpdate {
		_ = os.Rename(legacyUpdatePath, legacyUpdatePath+".bak")
	}

	return nil
}
