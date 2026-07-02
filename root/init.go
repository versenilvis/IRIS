package root

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/versenilvis/iris/config"
)

var initCmd = &cobra.Command{
	Use:   "init [bash|zsh|fish]",
	Short: "Generate the autostart script for your shell",
	Long: `Add the output of this command to your shell's configuration file to start Iris automatically.
For example, add this to your ~/.zshrc:
  eval "$(iris init zsh)"`,
	ValidArgs: []string{"bash", "zsh", "fish"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		shell := args[0]
		switch shell {
		case "zsh":
			fmt.Printf(`
# Iris Autostart Hook
if [ -n "$TMUX" ] && [ -n "$IRIS_PID" ]; then
    if ps -o comm= -p $PPID 2>/dev/null | grep -q "tmux"; then
        unset IRIS_PID IRIS_IS_CHILD IRIS_FD
    fi
fi

if [ -z "$IRIS_PID" ]; then
    export IRIS_ACTIVE_SHELL="zsh"
    exec iris
fi

# Iris Autocomplete Hook
if [ -n "$IRIS_PID" ] && [ -n "$IRIS_FD" ]; then
  _iris_send_lbuffer() {
    print -u $IRIS_FD -N -r -- "$LBUFFER" 2>/dev/null
  }

  _iris_precmd() {
    print -u $IRIS_FD -N -r -- "IRIS_CMD_STOP" 2>/dev/null
  }

  autoload -Uz add-zle-hook-widget
  autoload -Uz add-zsh-hook
  
  add-zle-hook-widget line-pre-redraw _iris_send_lbuffer
  add-zsh-hook precmd _iris_precmd
fi
`)
		case "bash":
			fmt.Printf(`
# Iris Autostart Hook
if [ -n "$TMUX" ] && [ -n "$IRIS_PID" ]; then
    if ps -o comm= -p $PPID 2>/dev/null | grep -q "tmux"; then
        unset IRIS_PID IRIS_IS_CHILD IRIS_FD
    fi
fi

if [ -z "$IRIS_PID" ]; then
    export IRIS_ACTIVE_SHELL="bash"
    exec iris
fi
`)
		case "fish":
			fmt.Printf(`
# Iris Autostart Hook
if set -q TMUX; and set -q IRIS_PID
    if ps -o comm= -p $PPID 2>/dev/null | grep -q "tmux"
        set -e IRIS_PID
        set -e IRIS_IS_CHILD
        set -e IRIS_FD
    end
end

if not set -q IRIS_PID
    set -gx IRIS_ACTIVE_SHELL "fish"
    exec iris
end
`)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(setupCmd)
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Automatically setup iris shell integration and install binary",
	Run: func(cmd *cobra.Command, args []string) {
		home, _ := os.UserHomeDir()

		localBin := filepath.Join(home, ".local", "bin")
		_ = os.MkdirAll(localBin, 0755)

		exe, _ := os.Executable()
		targetExe := filepath.Join(localBin, "iris")

		fmt.Printf("Installing iris to %s...\n", targetExe)
		input, err := os.ReadFile(exe)
		if err != nil {
			fmt.Printf("Failed to read current executable: %v\n", err)
			return
		}

		_ = os.Remove(targetExe)
		err = os.WriteFile(targetExe, input, 0755)
		if err != nil {
			fmt.Printf("Failed to write to %s: %v\n", targetExe, err)
			return
		}

		shellPath := os.Getenv("SHELL")
		shellName := filepath.Base(shellPath)
		var configFile string
		var evalCmd string

		switch shellName {
		case "zsh":
			configFile = filepath.Join(home, ".zshrc")
			evalCmd = `eval "$(iris init zsh)"`
		case "bash":
			configFile = filepath.Join(home, ".bashrc")
			evalCmd = `eval "$(iris init bash)"`
		case "fish":
			configFile = filepath.Join(home, ".config", "fish", "config.fish")
			evalCmd = `iris init fish | source`
		default:
			fmt.Printf("Unsupported shell: %s. Please add iris init manually.\n", shellName)
			return
		}

		content, _ := os.ReadFile(configFile)
		if strings.Contains(string(content), "iris init") {
			fmt.Printf("Iris is already configured in %s\n", configFile)
		} else {
			f, err := os.OpenFile(configFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				fmt.Printf("Failed to update %s: %v\n", configFile, err)
				return
			}
			defer func() { _ = f.Close() }()

			_, _ = f.WriteString("\n# Iris Autocomplete\n" + evalCmd + "\n")
			fmt.Printf("✓ Added iris integration to %s\n", configFile)
		}

		// initialize default config file if it does not exist
		if path, err := config.ConfigPath(); err == nil {
			if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
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
`
				if errWrite := os.WriteFile(path, []byte(defaultContent), 0644); errWrite == nil {
					fmt.Printf("✓ Initialized default config file at %s\n", path)
				}
			}
		}

		fmt.Println("\nSetup complete! Please restart your terminal or run:")
		fmt.Printf("  \033[32msource %s\033[0m\n", configFile)
	},
}
