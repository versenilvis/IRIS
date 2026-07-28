package root

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"github.com/versenilvis/iris/internal/config"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "manage iris configuration",
}

var ConfigInitCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize default configuration file with comments",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := config.ConfigPath()
		if err != nil {
			fmt.Printf("failed to get config path: %v\n", err)
			return
		}

		if _, statErr := os.Stat(path); statErr == nil {
			fmt.Printf("config file already exists at %s\n", path)
			return
		}

		_ = os.MkdirAll(filepath.Dir(path), 0755)

		defaultContent := `# ~/.config/iris/config.toml
# iris configuration file

[core]
# schema version
# do not edit this field manually
version = 1

# override shell: "bash", "zsh", "fish", keep empty for auto detection
shell = ""

# startup mode: "last", "spec", "history"
# "last" = remember last mode used
mode = "last"

# enable debug logging
debug = false

[ui]
# visual style: "modern" (icons, category pills, shortcut footer) or "classic" (minimalist, centered number, no icons)
style = "modern"

# enable Nerd Fonts icons in overlay menu
nerd-fonts = true

# enable inline ghost text
ghost-text = true

# maximum suggestions to display
max-suggestions = 100

# maximum height of the overlay
max-height = 15

[git]
# hide current branch in checkout/switch list
filter-active-branch = true

# merge remote and local branches with same name
deduplicate-branches = true

[updater]
# check for updates on startup
check-on-startup = true

# update channel: "stable", "nightly"
channel = "stable"

# interval between update checks, e.g. "24h", "6h", "30m"
check-interval = "24h"

# ── AI Suggestions ───────────────────────────────────────────
# IRIS supports two API protocols:
#   "openai"    — OpenAI-compatible chat completions (OpenRouter, Groq, Ollama,
#                 DeepSeek, LM Studio, Together, Fireworks, xAI, etc.)
#   "anthropic" — Anthropic native Messages API (Claude models)
#
# Set inherited_from to the protocol your provider uses.
# Use api_key_env (recommended) to avoid hardcoding keys in plaintext.
#
# Pick ONE provider by uncommenting and setting ai.provider to its name.

[ai]
# enable AI-powered command suggestions
enabled = false

# which provider block to use (must match a name under [ai.providers.*])
provider = ""

# delay (ms) before firing an AI request after the user stops typing
debounce_ms = 500

# minimum interval (ms) between successive AI requests
min_interval_ms = 1000

# ── Provider: OpenRouter (any model, OpenAI-compatible) ──────
# [ai.providers.openrouter]
# inherited_from = "openai"
# endpoint = "https://openrouter.ai/api/v1"
# api_key_env = "OPENROUTER_API_KEY"
# model = "google/gemini-2.0-flash-001"
# timeout_ms = 3000

# ── Provider: Groq (fast, free tier) ─────────────────────────
# [ai.providers.groq]
# inherited_from = "openai"
# endpoint = "https://api.groq.com/openai/v1"
# api_key_env = "GROQ_API_KEY"
# model = "llama-3.3-70b-versatile"
# timeout_ms = 3000

# ── Provider: Ollama (local, no API key) ─────────────────────
# [ai.providers.ollama]
# inherited_from = "openai"
# endpoint = "http://localhost:11434/v1"
# model = "qwen2.5-coder"
# timeout_ms = 5000

# ── Provider: Anthropic (native Claude API) ──────────────────
# [ai.providers.anthropic]
# inherited_from = "anthropic"
# endpoint = "https://api.anthropic.com/v1/messages"
# api_key_env = "ANTHROPIC_API_KEY"
# model = "claude-3-5-haiku-20241022"
# timeout_ms = 3000

# ── Provider: DeepSeek ───────────────────────────────────────
# [ai.providers.deepseek]
# inherited_from = "openai"
# endpoint = "https://api.deepseek.com/v1"
# api_key_env = "DEEPSEEK_API_KEY"
# model = "deepseek-chat"
# timeout_ms = 3000

# ── Provider: LM Studio (local) ──────────────────────────────
# [ai.providers.lmstudio]
# inherited_from = "openai"
# endpoint = "http://localhost:1234/v1"
# model = "local-model"
# timeout_ms = 5000

# ── Provider: OpenAI ─────────────────────────────────────────
# [ai.providers.openai]
# inherited_from = "openai"
# endpoint = "https://api.openai.com/v1"
# api_key_env = "OPENAI_API_KEY"
# model = "gpt-4o-mini"
# timeout_ms = 3000

# ── Suggest-on-empty (contextual hints on blank prompt) ──────
# [ai.suggest_on_empty]
# enabled = false
# debounce_ms = 800
# min_interval_ms = 5000
`
		err = os.WriteFile(path, []byte(defaultContent), 0644)
		if err != nil {
			fmt.Printf("failed to write config file: %v\n", err)
			return
		}
		fmt.Printf("initialized config file at %s\n", path)
	},
}

var ConfigShowCmd = &cobra.Command{
	Use:   "show",
	Short: "show the resolved configuration",
	Run: func(cmd *cobra.Command, args []string) {
		enc := toml.NewEncoder(cmd.OutOrStdout())
		if err := enc.Encode(config.Get()); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "failed to encode config: %v\n", err)
		}
	},
}

func init() {
	ConfigCmd.AddCommand(ConfigInitCmd)
	ConfigCmd.AddCommand(ConfigShowCmd)
	rootCmd.AddCommand(ConfigCmd)
}
