# IRIS documentation

Iris is a fast terminal autocomplete assistant written in Go. It wraps around your shell (Zsh, Bash, or Fish) to give you real-time command suggestions, a floating dropdown menu, and smart history search right where you type

## Table of contents

- [Getting started](#getting-started)
- [Usage guide](#usage-guide)
- [Commands reference](#commands-reference)
- [Configuration guide](#configuration-guide)
- [Troubleshooting guide](#troubleshooting-guide)
- [Developer guide](development.md)

## Getting started

### Dependencies

Before installing Iris, ensure your system meets the following requirements:
- OS: Linux or macOS
- Terminal emulator with ANSI escape sequence support
- Go 1.24 or newer if building from source

### Installation

You can install Iris using any of the three methods below:

#### Method 1: Install script (recommended)

Quickly download and install the precompiled binary using our install script:

```bash
curl -sSL https://raw.githubusercontent.com/versenilvis/iris/main/scripts/install.sh | sh
```

#### Method 2: Go install

If you already have Go installed on your system, install directly to your `GOBIN`:

```bash
go install github.com/versenilvis/iris@latest
```

Ensure that your `$GOPATH/bin` or `$HOME/go/bin` directory is added to your system `PATH`

#### Method 3: Build from source

To develop or compile Iris locally from source, install `just` command runner and execute `just reload`:

```bash
git clone https://github.com/versenilvis/iris.git
cd iris
just reload
```

### Shell integration setup

To enable intelligent auto-completion and overlay rendering, connect Iris to your shell environment

For Zsh (`~/.zshrc`):

```zsh
if command -v iris >/dev/null 2>&1; then
    alias i="iris"
fi
```

For Bash (`~/.bashrc`):

```bash
if command -v iris >/dev/null 2>&1; then
    alias i="iris"
fi
```

For Fish (`~/.config/fish/config.fish`):

```fish
if command -v iris >/dev/null 2>&1
    alias i="iris"
end
```

Verify your installation by running:

```bash
iris version
```

## Usage guide

Once inside an interactive Iris session, your shell receives powerful real-time auto-completion overlays and navigation enhancements

### Core navigation

When you type a command or query, Iris displays a floating overlay box positioned directly below your cursor with matching suggestions

- Up arrow (`↑`): Move the selection cursor up through the suggestion list
- Down arrow (`↓`): Move the selection cursor down through the suggestion list
- Tab: Insert the currently highlighted suggestion directly into your command line buffer without executing it
- Enter: Execute the currently highlighted command immediately or submit the text typed in your prompt
- Esc: Temporarily dismiss and hide the overlay suggestion menu for the current line
- Shift+Tab: Permanently disable overlay suggestions until toggled back on

### Mode switching

Iris operates in two primary autocomplete modes: specification suggestions and history navigation

You can toggle between these modes instantly at any time by pressing `Ctrl+R`

- Specification mode (`spec`): Suggests available subcommands, flags, and arguments based on built-in tool specifications (such as git, npm, docker, or go)
- History mode (`history`): Suggests previous commands from your shell history that match your typed prefix using fuzzy matching

When navigating with arrow keys and pressing `Ctrl+R`, Iris resets the menu selection and immediately loads suggestions for your typed query under the newly active mode

### Instant alias expansion

Iris supports POSIX and shell aliases defined in your configuration

When you type an alias keyword and press the `Space` key, Iris immediately expands the alias into its full command inside your command line buffer

### Ghost text autosuggestions

When ghost text is enabled, Iris displays inline completion suggestions ahead of your cursor in a muted style

- Right Arrow (`→`): Accept the inline ghost text suggestion and append it to your current command line buffer

## Commands reference

Iris provides a comprehensive set of CLI commands to manage shell integration, updates, configuration, and diagnostics:

```bash
# start interactive autocomplete session (wrapped terminal wrapper)
iris [flags]
  -s, --shell <shell>   specify target shell environment (zsh, bash, fish)
  -d, --debug           enable runtime debug logging to ~/.cache/iris/iris.log

# shell integration setup and initialization
iris setup [shell]      automatically configure shell integration in RC file and initialize default config
iris init <shell>       output raw shell wrapper code for manual evaluation in profile scripts

# configuration management
iris config init        initialize default configuration file at ~/.config/iris/config.toml
iris config show        output current active/resolved configuration in TOML format

# maintenance and diagnostics
iris update             check GitHub release tracks and update binary to latest release
iris version            print current semantic version string
iris uninstall          remove shell integration hooks from RC files and uninstall Iris binary
iris crash-log          display file path to the latest captured stack trace report
iris crash-log --clear  remove all stored crash logs from ~/.cache/iris/crashes
```

## Configuration guide

Iris uses a clean TOML configuration file located at `~/.config/iris/config.toml` to customize UI presentation, suggestion behavior, and core engine settings

### Default configuration structure

Below is a complete sample configuration template with all available parameters and comments:

```toml
# ~/.config/iris/config.toml
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
```

### Configuration sections

- `[core]`: Defines core engine settings including schema version, forced target shell (`bash`, `zsh`, `fish`), initial startup mode (`last`, `spec`, `history`), and diagnostic debug logging
- `[ui]`: Controls visual layout such as inline ghost text autocompletion, total suggestions rendered, and maximum overlay box height
- `[git]`: Specialized Git autocompletion behavior such as filtering out the currently active branch when switching and deduplicating local and remote branch names
- `[updater]`: Configures automated background updates, checking frequency (`check-interval`), and release track (`channel`)

### CLI configuration commands

Iris provides built-in commands to inspect and modify settings directly from your shell:

```bash
iris config init
iris config show
```

## Troubleshooting guide

If you encounter unexpected behavior while running Iris, use the diagnostic procedures outlined below

### Debug mode and logging

To inspect overlay positioning or command interception bugs in real time, launch Iris in debug mode:

```bash
iris -d
```

Runtime activities are written directly to the log file located at `~/.cache/iris/iris.log`. You can follow real-time events by running:

```bash
tail -f ~/.cache/iris/iris.log
```

> [!Caution]
> ### Common issues and solutions
> - Overlay Position Misaligned: Multi-line shell prompts (such as Starship) emitting ANSI codes can cause offset issues. Iris automatically parses carriage returns (`\r`) and horizontal escape codes (`\033[nC`) to anchor visual columns accurately
> - Menu Disappears When Toggling Mode: When pressing `Ctrl+R` after arrow navigation, Iris restores the original search query automatically so suggestions reload cleanly
> - Inspecting Crash Logs: If the session terminates unexpectedly, inspect captured stack traces by running `iris crash-log`
