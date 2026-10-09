package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/versenilvis/iris/integration/shell"
	"github.com/versenilvis/iris/internal/config"
)

var InitCmd = &cobra.Command{
	Use:   "init [bash|zsh|fish]",
	Short: "Generate the autostart script for your shell",
	Long: `Add the output of this command to your shell's configuration file to start Iris automatically.
For example, add this to your ~/.zshrc:
  eval "$(iris init zsh)"`,
	ValidArgs: []string{"bash", "zsh", "fish"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		shellName := args[0]
		switch shellName {
		case "zsh":
			fmt.Printf(`
# Iris Autostart Hook
# a multiplexer pane inherits IRIS_* but runs on its own tty, so those vars
# point at an iris that is not driving this terminal
if [ -n "$IRIS_PID" ] && [ "$IRIS_PID" != "$PPID" ] && [ "${TTY:-$(tty 2>/dev/null)}" != "$IRIS_TTY" ]; then
    unset IRIS_PID IRIS_IS_CHILD IRIS_FD IRIS_TTY
fi

# a non-interactive shell (tool runners sourcing rc files, scripts) has no
# prompt to complete, and exec'ing here would seize the tty from the real iris
if [[ -o interactive ]] && [ -t 0 ] && [ -z "$IRIS_PID" ] && [ -z "$IRIS_RESCUE" ] && [ -z "$ZSH_EXECUTION_STRING" ]; then
    export IRIS_ACTIVE_SHELL="zsh"
    if [ -n "$FNM_MULTISHELL_PATH" ]; then
        path=(${path:#*fnm_multishells*})
        unset FNM_MULTISHELL_PATH
    fi
    exec iris
fi

# Iris Autocomplete Hook
if [ -n "$IRIS_PID" ] && [ -n "$IRIS_FD" ]; then
  # the cursor comes with the line: sending only $LBUFFER made iris drop
  # everything to the right of the cursor and believe it sat at the end
  _iris_send_lbuffer() {
    print -u $IRIS_FD -N -r -- "IRIS_LINE:${#LBUFFER}:$BUFFER" 2>/dev/null
  }

  _iris_sync_cwd() {
    print -u $IRIS_FD -N -r -- "IRIS_CWD:$PWD" 2>/dev/null
  }

  _iris_send_aliases() {
    local a
    a="$(alias -L 2>/dev/null)"
    if [[ "$a" != "$_iris_last_aliases" ]]; then
      _iris_last_aliases="$a"
      print -u $IRIS_FD -N -r -- "IRIS_ALIASES:$a" 2>/dev/null
    fi
  }

  _iris_precmd() {
    local iris_exit_code=$?
    _iris_sync_cwd
    _iris_send_aliases
    print -u $IRIS_FD -N -r -- "IRIS_CMD_STOP:$iris_exit_code" 2>/dev/null
  }

  _iris_preexec() {
    print -u $IRIS_FD -N -r -- "IRIS_CMD_START" 2>/dev/null
  }

  autoload -Uz add-zle-hook-widget
  autoload -Uz add-zsh-hook

  add-zle-hook-widget line-pre-redraw _iris_send_lbuffer
  add-zsh-hook precmd _iris_precmd
  add-zsh-hook preexec _iris_preexec
  add-zsh-hook chpwd _iris_sync_cwd
fi
`)
		case "bash":
			fmt.Printf(`
# Iris Autostart Hook
# a multiplexer pane inherits IRIS_* but runs on its own tty, so those vars
# point at an iris that is not driving this terminal
if [ -n "$IRIS_PID" ] && [ "$IRIS_PID" != "$PPID" ] && [ "$(tty 2>/dev/null)" != "$IRIS_TTY" ]; then
    unset IRIS_PID IRIS_IS_CHILD IRIS_FD IRIS_TTY
fi

# a non-interactive shell (tool runners sourcing rc files, scripts) has no
# prompt to complete, and exec'ing here would seize the tty from the real iris
if [[ $- == *i* ]] && [ -t 0 ] && [ -z "$IRIS_PID" ] && [ -z "$IRIS_RESCUE" ] && [ -z "$BASH_EXECUTION_STRING" ]; then
    export IRIS_ACTIVE_SHELL="bash"
    if [ -n "$FNM_MULTISHELL_PATH" ]; then
        PATH="${PATH//:$FNM_MULTISHELL_PATH\/bin:/:}"
        PATH="${PATH//:$FNM_MULTISHELL_PATH:/:}"
        PATH="${PATH#:}"
        PATH="${PATH%%:}"
        export PATH
        unset FNM_MULTISHELL_PATH
    fi
    exec iris
fi

# Iris Autocomplete Hook
if [ -n "$IRIS_PID" ] && [ -n "$IRIS_FD" ]; then
  _iris_last_line=""
  _iris_last_aliases=""

  _iris_send_aliases() {
    local a
    a="$(alias -p 2>/dev/null)"
    if [[ "$a" != "$_iris_last_aliases" ]]; then
      _iris_last_aliases="$a"
      printf 'IRIS_ALIASES:%%s\0' "$a" >&$IRIS_FD 2>/dev/null
    fi
  }

  _iris_bash_cmd_start() {
    printf 'IRIS_CMD_START\0' >&$IRIS_FD 2>/dev/null
  }

  _iris_send_line() {
    if [[ "$READLINE_LINE" != "$_iris_last_line" ]]; then
      _iris_last_line="$READLINE_LINE"
      printf 'IRIS_LINE:%%d:%%s\0' "$READLINE_POINT" "$READLINE_LINE" >&$IRIS_FD 2>/dev/null
    fi
  }

  _iris_precmd() {
    local iris_exit_code=$?
    printf 'IRIS_CWD:%%s\0' "$PWD" >&$IRIS_FD 2>/dev/null
    _iris_send_aliases
    printf 'IRIS_CMD_STOP:%%d\0' "$iris_exit_code" >&$IRIS_FD 2>/dev/null
  }

  trap '_iris_bash_cmd_start' DEBUG
  PROMPT_COMMAND="_iris_precmd${PROMPT_COMMAND:+;$PROMPT_COMMAND}"
  bind -x '"\e[200~": _iris_send_line' 2>/dev/null
fi
`)
		case "fish":
			disableFishAutosuggest := ""
			if config.Get().UI.GhostText != config.GhostTextOff {
				disableFishAutosuggest = `
    set -g fish_autosuggestion_enabled 0
`
			}
			fmt.Printf(`
# Iris Autostart Hook
# a multiplexer pane inherits IRIS_* but runs on its own tty, so those vars
# point at an iris that is not driving this terminal
if set -q IRIS_PID; and test "$IRIS_PID" != "$fish_pid"; and test (tty 2>/dev/null) != "$IRIS_TTY"
    set -e IRIS_PID IRIS_IS_CHILD IRIS_FD IRIS_TTY
end

# a non-interactive shell (tool runners sourcing rc files, scripts) has no
# prompt to complete, and exec'ing here would seize the tty from the real iris
if status is-interactive; and test -t 0; and not set -q IRIS_PID; and not set -q IRIS_RESCUE
    set -gx IRIS_ACTIVE_SHELL "fish"
    if set -q FNM_MULTISHELL_PATH
        set -l fnm_idx (contains -i $FNM_MULTISHELL_PATH $fish_user_paths)
        if test -n "$fnm_idx"
            set -e fish_user_paths[$fnm_idx]
        end
        set -e FNM_MULTISHELL_PATH
    end
    exec iris
end

# Iris Autocomplete Hook
if set -q IRIS_PID; and set -q IRIS_FD
%s
    set -g _iris_last_line ""
    set -g _iris_last_aliases ""
    set -g _iris_last_abbrs ""
    set -g _iris_last_func_count 0

    function _iris_postexec --on-event fish_postexec
        printf "IRIS_CMD_STOP:%%d\x00" $status >&$IRIS_FD 2>/dev/null
    end

    function _iris_fish_precmd --on-event fish_prompt
        printf "IRIS_CWD:%%s\x00" "$PWD" >&$IRIS_FD 2>/dev/null

        # fish abbr is builtin and faster, so we check it first
        set -l abbrs (abbr --show 2>/dev/null | string collect)
        if test "$abbrs" != "$_iris_last_abbrs"
            set -g _iris_last_abbrs "$abbrs"
            printf "IRIS_ABBRS:%%s\x00" "$abbrs" >&$IRIS_FD 2>/dev/null
        end

        # alias in fish is wrapper of function, so we count function first to avoid call alias every prompt
        set -l func_count (functions -n 2>/dev/null | count)
        if test "$func_count" != "$_iris_last_func_count"
            set -g _iris_last_func_count "$func_count"
            set -l aliases (alias 2>/dev/null | string collect)
            set -g _iris_last_func_count (functions -n 2>/dev/null | count)
            set -g _iris_last_aliases "$aliases"
            printf "IRIS_ALIASES:%%s\x00" "$aliases" >&$IRIS_FD 2>/dev/null
        end
    end
    function _iris_fish_preexec --on-event fish_preexec
        printf "IRIS_CMD_START\x00" >&$IRIS_FD 2>/dev/null
    end
end
`, disableFishAutosuggest)
		}
	},
}

var SetupCmd = &cobra.Command{
	Use:   "setup [shell]",
	Short: "Automatically setup iris shell integration and install binary",
	Args:  cobra.MaximumNArgs(1),
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

		var shellName string
		if len(args) > 0 {
			shellName = filepath.Base(args[0])
		} else {
			shellPath := os.Getenv("SHELL")
			shellName = filepath.Base(shellPath)
		}
		var configFile string
		var evalCmd string

		switch shellName {
		case "zsh":
			configFile = filepath.Join(shell.GetZshConfigDir(), ".zshrc")
			evalCmd = `eval "$(iris init zsh)"`
		case "bash":
			configFile = filepath.Join(home, ".bashrc")
			evalCmd = `eval "$(iris init bash)"`
		case "fish":
			configFile = filepath.Join(shell.GetFishConfigDir(), "config.fish")
			evalCmd = `iris init fish | source`
		default:
			fmt.Printf("Unsupported shell: %s. Please add iris init manually.\n", shellName)
			return
		}

		content, readErr := os.ReadFile(configFile)
		if strings.Contains(string(content), "iris init") {
			fmt.Printf("Iris is already configured in %s\n", configFile)
		} else {
			var newContent string
			if readErr == nil && len(content) > 0 {
				newContent = "# Iris Autocomplete\n" + evalCmd + "\n\n" + string(content)
			} else {
				newContent = "# Iris Autocomplete\n" + evalCmd + "\n"
			}
			if mkdirErr := os.MkdirAll(filepath.Dir(configFile), 0755); mkdirErr != nil {
				fmt.Printf("Failed to create directory for %s: %v\n", configFile, mkdirErr)
				return
			}
			err = os.WriteFile(configFile, []byte(newContent), 0644)
			if err != nil {
				fmt.Printf("Failed to update %s: %v\n", configFile, err)
				return
			}
			fmt.Printf("✓ Added iris integration to top of %s\n", configFile)
		}

		if path, err := config.ConfigPath(); err == nil {
			if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
				_ = os.MkdirAll(filepath.Dir(path), 0755)
				defaultContent := `# ~/.config/iris/config.toml
# iris configuration file

[core]
version = 1
shell = ""
shell-login = false
mode = "last"
debug = false
expand-alias = true
auto-execute = false
navigate-closed = "history"

[ui]
style = "modern"
nerd-fonts = true
hidden-files = false
ghost-text = 1
max-suggestions = 100
max-height = 6
max-width = 0

[git]
filter-active-branch = true
deduplicate-branches = true

[updater]
check-on-startup = true
channel = "stable"
check-interval = "24h"
auto-update = 0

[zoxide]
extend-cd = false

[keybindings]
toggle-mode = "ctrl+r"
toggle-menu = "shift+tab"
select = "tab"
navigate-up = "up"
navigate-down = "down"
navigate-right = "right"
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
