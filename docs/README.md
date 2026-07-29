# User guide

Iris is a fast terminal autocomplete assistant written in Go. It wraps around your shell (Zsh, Bash, or Fish) to give you real-time command suggestions, a floating dropdown menu, and smart history search right where you type.

## Table of contents

- [Getting started](#getting-started)
- [Shortcuts](#shortcuts)
- [Configuration guide](#configuration-guide)
- [Reporting bugs](#reporting-bugs)
- [Developer documentation](dev/README.md)

## Getting started

### Dependencies

- OS: Linux or macOS
- Terminal emulator with ANSI color support
- Go 1.24 or newer (if building from source)

### Installation

#### Method 1: Install script (recommended)

```bash
curl -sSL https://raw.githubusercontent.com/versenilvis/iris/main/scripts/install.sh | sh
```

#### Method 2: Go install

```bash
go install github.com/versenilvis/iris@latest
```

#### Method 3: Build from source

```bash
git clone https://github.com/versenilvis/iris.git
cd iris
just reload
```

### Shell setup

Add an alias to your shell configuration file to launch Iris easily:

**Zsh (`~/.zshrc`):**
```zsh
if command -v iris >/dev/null 2>&1; then
    alias i="iris"
fi
```

**Bash (`~/.bashrc`):**
```bash
if command -v iris >/dev/null 2>&1; then
    alias i="iris"
fi
```

**Fish (`~/.config/fish/config.fish`):**
```fish
if command -v iris >/dev/null 2>&1
    alias i="iris"
end
```

## Shortcuts

| Shortcut                           | Action                  | Description                                                               |
| :--------------------------------- | :---------------------- | :------------------------------------------------------------------------ |
| <kbd>Shift</kbd> + <kbd>Tab</kbd>  | Toggle menu             | Show or hide the suggestion menu.                                         |
| <kbd>Esc</kbd>                     | Hide menu               | Temporarily hide the menu until the next key press.                       |
| <kbd>Tab</kbd>                     | Accept suggestion       | Insert the currently selected suggestion into the prompt.                 |
| <kbd>Enter</kbd>                   | Execute command         | Close the menu and send the current command to the shell.                 |
| <kbd>↑</kbd>                       | Navigate up / history   | Move the selection up, or open command history when the prompt is empty.  |
| <kbd>↓</kbd>                       | Navigate down / history | Move the selection down, or open command history when the prompt is empty.|
| <kbd>→</kbd>                       | Accept ghost text       | Accept the faded ghost text suggestion when the menu is open.             |
| <kbd>←</kbd> / <kbd>→</kbd>        | Move cursor             | Move the cursor inside the input buffer. Disabled when the prompt is empty|
| <kbd>Ctrl</kbd> + <kbd>R</kbd>     | Switch mode             | Toggle between `spec` and `history` mode.                                 |
| <kbd>Ctrl</kbd> + <kbd>A</kbd>     | Beginning of line       | Move the cursor to the start of the command line.                         |
| <kbd>Ctrl</kbd> + <kbd>E</kbd>     | End of line             | Move the cursor to the end of the command line.                           |
| <kbd>Ctrl</kbd> + <kbd>L</kbd>     | Clear screen            | Clear the terminal while preserving the input buffer and redrawing the menu. |
| <kbd>Ctrl</kbd> + <kbd>U</kbd>     | Clear command           | Remove the entire current command and close the menu.                     |
| <kbd>Ctrl</kbd> + <kbd>C</kbd>     | Cancel command          | Send `SIGINT`, clear the input buffer, and close the menu.                |
| <kbd>Ctrl</kbd> + <kbd>W</kbd>     | Delete word             | Delete the word immediately before the cursor.                            |

> [!NOTE]
> With <kbd>Ctrl</kbd> + <kbd>A</kbd>, <kbd>Ctrl</kbd> + <kbd>E</kbd>, <kbd>Ctrl</kbd> + <kbd>W</kbd>, <kbd>Ctrl</kbd> + <kbd>U</kbd>, <kbd>Ctrl</kbd> + <kbd>L</kbd>, and <kbd>Ctrl</kbd> + <kbd>C</kbd>: they belong to your shell by default. IRIS handles them directly in raw mode so your cursor and menu stay in sync

## Configuration guide

Iris uses a clean TOML configuration file located at `~/.config/iris/config.toml`.

### Creating & viewing config

```bash
iris config init
iris config show
```

### Sample `config.toml`

```toml
[core]
version = 1
shell = ""        # "zsh", "bash", "fish", or empty for auto-detection
mode = "last"     # "last", "spec", or "history"
debug = false
expand-alias = true

[ui]
style = "modern"  # "modern" or "classic"
ghost-text = true
hidden-files = false
max-suggestions = 100
max-height = 15
nerd-fonts = true

[keybindings]
toggle-mode = "ctrl+r"
toggle-menu = "ctrl+space"

[git]
filter-active-branch = true
deduplicate-branches = true

[updater]
check-on-startup = true
channel = "stable" # "stable" or "nightly"
check-interval = "24h"

[ai]
enabled = false
provider = "groq" # "groq" or "ollama"
debounce_ms = 400

[ai.providers.groq]
endpoint = "https://api.groq.com/openai/v1/chat/completions"
api_key_env = "GROQ_API_KEY" # or set api_key directly
model = "llama-3.3-70b-versatile"
timeout_ms = 3000

[ai.providers.ollama]
endpoint = "http://localhost:11434/v1/chat/completions"
model = "qwen2.5-coder"
timeout_ms = 5000
```

> [!NOTE]
> Using `api_key_env` is recommended over hardcoding `api_key` in plain text to keep credentials out of configuration files.

## Reporting bugs
> [!NOTE]
> When submitting a bug report, please include:
> - A detailed description of the bug and steps to reproduce it
> - Relevant log files captured while running in debug mode

Run IRIS with debug mode:
```bash
iris -d
```
or `config.toml`:
```toml
debug=true
```

> [!IMPORTANT]
> **Since IRIS logs everything you type, you should only enable debug mode when you need to report bugs**
