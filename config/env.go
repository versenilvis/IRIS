package config

import (
	"os"
	"strconv"
	"time"
)

func applyEnv(cfg *Config) {
	if val := os.Getenv("IRIS_CORE_DEBUG"); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			cfg.Core.Debug = b
		}
	}
	if val := os.Getenv("IRIS_CORE_SHELL"); val != "" {
		cfg.Core.Shell = val
	}
	if val := os.Getenv("IRIS_CORE_MODE"); val != "" {
		cfg.Core.Mode = val
	}
	if val := os.Getenv("IRIS_UI_GHOST_TEXT"); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			cfg.UI.GhostText = b
		}
	}
	if val := os.Getenv("IRIS_UI_MAX_SUGGESTIONS"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			cfg.UI.MaxSuggestions = i
		}
	}
	if val := os.Getenv("IRIS_UI_MAX_HEIGHT"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			cfg.UI.MaxHeight = i
		}
	}
	if val := os.Getenv("IRIS_UPDATER_CHANNEL"); val != "" {
		cfg.Updater.Channel = val
	}
	if val := os.Getenv("IRIS_UPDATER_INTERVAL"); val != "" {
		if dur, err := time.ParseDuration(val); err == nil {
			cfg.Updater.CheckInterval = Duration(dur)
		}
	}
	if val := os.Getenv("IRIS_UPDATER_CHECK_ON_STARTUP"); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			cfg.Updater.CheckOnStartup = b
		}
	}
	if val := os.Getenv("IRIS_AI_ENABLED"); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			cfg.AI.Enabled = b
		}
	}
	if val := os.Getenv("IRIS_AI_PROVIDER"); val != "" {
		cfg.AI.Provider = val
	}
}
