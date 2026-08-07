package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultConfigAndState(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Core.Version != 1 {
		t.Errorf("expected version 1, got %d", cfg.Core.Version)
	}
	if cfg.Core.ShellLogin {
		t.Errorf("expected login shell to be disabled by default")
	}
	if cfg.UI.MaxSuggestions != 100 {
		t.Errorf("expected suggestions 100, got %d", cfg.UI.MaxSuggestions)
	}
	if cfg.AI.Enabled {
		t.Errorf("expected AI to be disabled by default")
	}
	if cfg.AI.Provider != "" {
		t.Errorf("expected default provider to be empty, got %q", cfg.AI.Provider)
	}
	if cfg.AI.Providers != nil {
		t.Errorf("expected default providers map to be nil, got %v", cfg.AI.Providers)
	}

	// test manual provider registration
	cfg.AI.Provider = "custom"
	cfg.AI.Providers = map[string]ProviderConfig{
		"custom": {
			InheritedFrom: "openai",
			Endpoint:      "https://custom-api.com/v1",
			APIKey:        "test-key",
			Model:         "test-model",
			TimeoutMS:     1000,
		},
	}
	p, ok := cfg.AI.GetActiveProvider()
	if !ok {
		t.Fatalf("expected custom provider to exist")
	}
	if p.InheritedFrom != "openai" {
		t.Errorf("expected inherited_from openai, got %q", p.InheritedFrom)
	}
	if p.GetAPIKey() != "test-key" {
		t.Errorf("expected api key test-key, got %q", p.GetAPIKey())
	}
	if cfg.AI.SuggestOnEmpty.DebounceMS != 800 {
		t.Errorf("expected debounce 800, got %d", cfg.AI.SuggestOnEmpty.DebounceMS)
	}
	if cfg.AI.SuggestOnEmpty.MinIntervalMS != 5000 {
		t.Errorf("expected min interval 5000, got %d", cfg.AI.SuggestOnEmpty.MinIntervalMS)
	}

	state := DefaultState()
	if state.LastMode != "spec" {
		t.Errorf("expected last mode spec, got %q", state.LastMode)
	}
}

func TestCustomDuration(t *testing.T) {
	var dur Duration
	err := dur.UnmarshalText([]byte("6h"))
	if err != nil {
		t.Fatalf("unexpected error unmarshalling duration: %v", err)
	}
	if time.Duration(dur) != 6*time.Hour {
		t.Errorf("expected 6 hours, got %v", time.Duration(dur))
	}

	b, err := dur.MarshalText()
	if err != nil {
		t.Fatalf("unexpected error marshaling duration: %v", err)
	}
	if string(b) != "6h0m0s" {
		t.Errorf("expected 6h0m0s, got %q", string(b))
	}

	err = dur.UnmarshalText([]byte("invalid"))
	if err == nil {
		t.Errorf("expected error for invalid duration")
	}
}

func TestValidationAndEnvironmentOverrides(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "iris-config-env-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// UserConfigDir follows HOME on macOS and XDG_CONFIG_HOME on Unix.
	// Override both so the test cannot read or write the user's real config.
	t.Setenv("HOME", tmpDir)
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	configPath, err := ConfigPath()
	if err != nil {
		t.Fatalf("failed to get config path: %v", err)
	}
	configDir := filepath.Dir(configPath)
	if mkErr := os.MkdirAll(configDir, 0755); mkErr != nil {
		t.Fatalf("failed to create config dir: %v", mkErr)
	}
	tomlContent := `
[core]
shell-login = true

[ai]
enabled = true
provider = "groq"

[ai.providers.groq]
inherited_from = "openai"
endpoint = "https://api.groq.com/openai/v1"
api_key_env = "GROQ_API_KEY"
model = "qwen-2.5-coder-32b"
`
	if wrErr := os.WriteFile(configPath, []byte(tomlContent), 0644); wrErr != nil {
		t.Fatalf("failed to write config file: %v", wrErr)
	}

	t.Setenv("IRIS_CORE_DEBUG", "true")
	t.Setenv("IRIS_CORE_SHELL", "fish")
	t.Setenv("IRIS_CORE_MODE", "history")
	t.Setenv("IRIS_UI_GHOST_TEXT", "false")
	t.Setenv("IRIS_UI_MAX_SUGGESTIONS", "250")
	t.Setenv("IRIS_UI_MAX_HEIGHT", "25")
	t.Setenv("IRIS_UPDATER_CHANNEL", "nightly")
	t.Setenv("IRIS_UPDATER_INTERVAL", "12h")
	t.Setenv("IRIS_UPDATER_CHECK_ON_STARTUP", "false")
	t.Setenv("IRIS_AI_PROVIDER", "ollama")
	t.Setenv("GROQ_API_KEY", "gsk_test_123")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if !cfg.Core.Debug {
		t.Errorf("expected debug to be true")
	}
	if !cfg.Core.ShellLogin {
		t.Errorf("expected login shell to be enabled from TOML")
	}
	if cfg.Core.Shell != "fish" {
		t.Errorf("expected shell fish, got %q", cfg.Core.Shell)
	}
	if cfg.Core.Mode != "history" {
		t.Errorf("expected mode history, got %q", cfg.Core.Mode)
	}
	if cfg.UI.GhostText {
		t.Errorf("expected ghost text to be false")
	}
	if cfg.UI.MaxSuggestions != 250 {
		t.Errorf("expected max suggestions 250, got %d", cfg.UI.MaxSuggestions)
	}
	if cfg.UI.MaxHeight != 25 {
		t.Errorf("expected max height 25, got %d", cfg.UI.MaxHeight)
	}
	if cfg.Updater.Channel != "nightly" {
		t.Errorf("expected channel nightly, got %q", cfg.Updater.Channel)
	}
	if time.Duration(cfg.Updater.CheckInterval) != 12*time.Hour {
		t.Errorf("expected 12h, got %v", time.Duration(cfg.Updater.CheckInterval))
	}
	if cfg.Updater.CheckOnStartup {
		t.Errorf("expected check on startup to be false")
	}
	if cfg.AI.Provider != "ollama" {
		t.Errorf("expected provider ollama from env, got %q", cfg.AI.Provider)
	}
	groqCfg := cfg.AI.Providers["groq"]
	if groqCfg.GetAPIKey() != "gsk_test_123" {
		t.Errorf("expected groq api key gsk_test_123 from env, got %q", groqCfg.GetAPIKey())
	}

	t.Setenv("IRIS_CORE_MODE", "invalid")
	_, err = Load()
	if err == nil {
		t.Errorf("expected validation error for invalid mode in env")
	}
}

func TestValidateAutoUpdateRange(t *testing.T) {
	cfg := DefaultConfig()

	for _, valid := range []int{0, 1, 2} {
		cfg.Updater.AutoUpdate = valid
		if err := validate(cfg); err != nil {
			t.Errorf("expected auto-update=%d to be valid, got error: %v", valid, err)
		}
	}

	for _, invalid := range []int{-1, 3} {
		cfg.Updater.AutoUpdate = invalid
		if err := validate(cfg); err == nil {
			t.Errorf("expected auto-update=%d to be rejected", invalid)
		}
	}
}

func TestLoadSave(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "iris-config-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// UserConfigDir follows HOME on macOS and XDG_CONFIG_HOME on Unix.
	// Override both so the test cannot read or write the user's real config.
	t.Setenv("HOME", tmpDir)
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	cfg.Core.Shell = "zsh"
	cfg.Core.ShellLogin = true
	cfg.UI.MaxHeight = 20

	err = Save(cfg)
	if err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("failed to load after save: %v", err)
	}

	if loaded.Core.Shell != "zsh" {
		t.Errorf("expected loaded shell to be zsh, got %q", loaded.Core.Shell)
	}
	if !loaded.Core.ShellLogin {
		t.Errorf("expected loaded login shell setting to be true")
	}
	if loaded.UI.MaxHeight != 20 {
		t.Errorf("expected loaded height to be 20, got %d", loaded.UI.MaxHeight)
	}
}

func TestMigration(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "iris-migrate-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	t.Setenv("HOME", tmpDir)
	t.Setenv("XDG_DATA_HOME", filepath.Join(tmpDir, ".local", "share"))

	legacyDir := filepath.Join(tmpDir, ".iris")
	if errMkdir := os.MkdirAll(legacyDir, 0755); errMkdir != nil {
		t.Fatalf("failed to create legacy dir: %v", errMkdir)
	}

	legacyStateJson := `{"mode": "history"}`
	_ = os.WriteFile(filepath.Join(legacyDir, "state.json"), []byte(legacyStateJson), 0644)

	legacyUpdateJson := `{"seen_version": "v1.2.3", "last_check": 1234567890}`
	_ = os.WriteFile(filepath.Join(legacyDir, "update_state.json"), []byte(legacyUpdateJson), 0644)

	err = MigrateFromLegacyJSON()
	if err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	state := LoadState()
	if state.LastMode != "history" {
		t.Errorf("expected migrated last mode 'history', got %q", state.LastMode)
	}
	if state.Updater.SeenVersion != "v1.2.3" {
		t.Errorf("expected migrated seen version 'v1.2.3', got %q", state.Updater.SeenVersion)
	}
	if state.Updater.LastCheckTime.Unix() != 1234567890 {
		t.Errorf("expected migrated check time 1234567890, got %v", state.Updater.LastCheckTime.Unix())
	}

	if _, err := os.Stat(filepath.Join(legacyDir, "state.json.bak")); err != nil {
		t.Errorf("expected backup file state.json.bak to exist")
	}
	if _, err := os.Stat(filepath.Join(legacyDir, "update_state.json.bak")); err != nil {
		t.Errorf("expected backup file update_state.json.bak to exist")
	}
}

func TestMatchKey(t *testing.T) {
	tests := []struct {
		input    []byte
		expected string
		matched  bool
		consumed int
	}{
		{[]byte{0x19}, "ctrl+y", true, 1},
		{[]byte{0x0b}, "ctrl+k", true, 1},
		{[]byte{0x0a}, "ctrl+j", true, 1},
		{[]byte{0x12}, "ctrl+r", true, 1},
		{[]byte{'j'}, "j", true, 1},
		{[]byte{'k'}, "k", true, 1},
		{[]byte{0x09}, "tab", true, 1},
		{[]byte{0x0d}, "enter", true, 1},
		{[]byte{0x0d}, "ctrl+r", false, 0},
		{[]byte("\x1b[106;4u"), "ctrl+j", true, 8},
		{[]byte("\x1b[106;5u"), "ctrl+j", true, 8},
		{[]byte("\x1b[107;4u"), "ctrl+k", true, 8},
		{[]byte("\x1b[107;12u"), "ctrl+k", true, 9},
		{[]byte("\x1b[106;1u"), "ctrl+j", false, 0},
		{[]byte("\x1b[97;4u"), "ctrl+a", true, 7},
		{[]byte("\x1b[106;4U"), "ctrl+j", false, 0},
		{[]byte("\x1b[106u"), "ctrl+j", false, 0},
	}

	for _, tt := range tests {
		m, c := MatchKey(tt.input, tt.expected)
		if m != tt.matched || c != tt.consumed {
			t.Errorf("MatchKey(%v, %q) = (%v, %d); want (%v, %d)", tt.input, tt.expected, m, c, tt.matched, tt.consumed)
		}
	}
}

// TestMatchKey_EnterReserved verifies that the Enter key (0x0d '\\r') can never
// be claimed by another keybinding. In a raw terminal Ctrl+M and the
// Enter/Return key are byte-identical (both 0x0d), so a "ctrl+m" keybinding
// must not shadow Enter, otherwise line submission (the Enter key) breaks.
func TestMatchKey_EnterReserved(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
		matched  bool
		consumed int
	}{
		{"ctrl+m must not match the Enter byte", []byte{0x0d}, "ctrl+m", false, 0},
		{"ctrl+m must not match Ctrl+M as a generic binding", []byte{0x0d}, "ctrl+m", false, 0},
		// other ctrl keys are still distinguishable from Enter and keep working
		{"ctrl+n unaffected", []byte{0x0e}, "ctrl+n", true, 1},
		{"ctrl+p unaffected", []byte{0x10}, "ctrl+p", true, 1},
		{"ctrl+j (0x0a) remains a distinct binding", []byte{0x0a}, "ctrl+j", true, 1},
	}

	for _, tt := range tests {
		m, c := MatchKey(tt.input, tt.expected)
		if m != tt.matched || c != tt.consumed {
			t.Errorf("%s: MatchKey(%v, %q) = (%v, %d); want (%v, %d)", tt.name, tt.input, tt.expected, m, c, tt.matched, tt.consumed)
		}
	}
}

// TestMatchKey_NavKeybindingsNoLongerHijackEnter is a regression guard: with a
// user-configured navigation (or any) keybinding set to "ctrl+m", pressing the
// Enter key must NOT be swallowed by that keybinding check.
func TestMatchKey_NavKeybindingsNoLongerHijackEnter(t *testing.T) {
	kb := []string{"ctrl+m", "ctrl+m", "<ctrl-m>", "CTRL+M"}
	for _, expected := range kb {
		if m, _ := MatchKey([]byte{0x0d}, expected); m {
			t.Errorf("MatchKey(enter{0x0d}, %q) matched; Enter must remain reserved", expected)
		}
	}
}
