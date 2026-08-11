<div align="center">
  <!-- <img width="50%" alt="banner" src="https://github.com/user-attachments/assets/c5ec623b-8259-473f-b7c3-3d01a64deb5d" /> -->
  <!-- <img width="25%" alt="logo" src="https://github.com/user-attachments/assets/79d3913c-56b7-42cb-8b07-53e98f39322b" /> -->
  <img width="15%" alt="logo" src="https://github.com/user-attachments/assets/10b7ca98-872b-44a2-bdcd-265f18aa0564" />

  <!-- <h1>IRIS</h1> -->
  
   
  
  [![macOS](https://img.shields.io/badge/macOS-FFFFFF?style=for-the-badge&logo=apple&logoColor=black)](https://www.apple.com/macos/)
  [![Linux](https://img.shields.io/badge/Linux-131415?style=for-the-badge&logo=linux&logoColor=white)](https://www.kernel.org/)
<br>
 <a href="https://trendshift.io/repositories/95378?utm_source=repository-badge&amp;utm_medium=badge&amp;utm_campaign=badge-repository-95378" target="_blank" rel="noopener noreferrer"><img src="https://trendshift.io/api/badge/repositories/95378" alt="versenilvis%2FIRIS | Trendshift" width="250" height="55"/></a>
  <!--[![GitHub Actions](https://img.shields.io/github/actions/workflow/status/versenilvis/IRIS/release.yml?branch=main&style=for-the-badge&logo=github&logoColor=white&label=Actions)](https://github.com/versenilvis/IRIS/actions/workflows/release.yml)-->
  [![Status](https://img.shields.io/badge/status-beta-yellow?style=for-the-badge&logo=github&logoColor=white)]()
  [![License: 0BSD](https://img.shields.io/badge/License-0BSD-blue?style=for-the-badge&logo=github&logoColor=white)](./LICENSE)
  [![Documentation](https://img.shields.io/badge/docs-available-brightgreen?style=for-the-badge&logo=github&logoColor=white)](#install)
  [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen?style=for-the-badge&logo=github&logoColor=white)](./CONTRIBUTING.md)
  
  <a href="#why-iris-instead-of-fig">Comparison</a> · <a href="#install">Install</a> · <a href="#default-shortcuts">Shortcuts</a> · <a href="#configuration-guide">Configuration</a> · <a href="#reporting-bugs">Reporting bugs</a>
  <p>IRIS (Intelligent Real-time Input Suggestion) - A shell auto-completion tool that works like code editor's IntelliSense</p>

</div>
<div align="center">
  <!-- <img width="800" height="450" alt="output" src="https://github.com/user-attachments/assets/fd586c5b-89b6-4f6a-af9c-29248db5edc3" />-->
  <!-- <img width="1280" height="720" alt="output" src="https://github.com/user-attachments/assets/a1003bda-7722-4185-9514-f1bf83fbb504" /> -->
  <img width="1920" height="1080" alt="Ghostty terminal showcase" src="https://github.com/user-attachments/assets/dc2434c2-4f59-4432-b77f-78e65c4d412e" />
  👻 <i>Ghostty terminal</i>
</div>

<br>

**IRIS is built on top of TTY, so it runs everywhere. It just needs a terminal!**

Run iris wherever you already work; your local machine, a remote server, or anywhere you can ssh. Each suggestion menu renders directly inline inside your real terminal session, not an app's imitation of one, so it never breaks full-screen TUIs or terminal formatting. Automatically index your aliases and shell history to suggest commands that match your actual workflow in real time. Change configurations and propagate them instantly without restarting your shell. One single local native Go binary, not an app: no gui, no electron, no mac-only wrapper, no account, no telemetry. (if you've used fig: it's that, rebuilt to run purely on TTY)

<!--
  [![Status](https://img.shields.io/badge/status-beta-yellow?style=for-the-badge&logo=github&logoColor=white)]()
  [![Documentation](https://img.shields.io/badge/docs-available-brightgreen?style=for-the-badge&logo=github&logoColor=white)](./docs)
  [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen?style=for-the-badge&logo=github&logoColor=white)](./.github/CONTRIBUTING.md)
-->

## AI suggestions

<div align="center">
  <img width="1920" height="1080" alt="output" src="https://github.com/user-attachments/assets/ab80fe75-b5dd-4acd-84bb-45bca17ee3b7" />
  <i>IRIS has AI suggestions like your code editor (API key/Local)</i>
</div>

## Why IRIS instead of Fig

> [!IMPORTANT]
> **[Fig](https://app.fig.io/) was officially sunset in September 2024 and migrated to Amazon Q Developer (which requires cloud authentication and proprietary bloat)**  
> **IRIS is the lightweight, open-source, zero-telemetry alternative built purely on native Go and TTY with no accounts, no GUI app, and no background daemons required**

### How it compares

| Feature                     | IRIS                 | Fig              |
| :-------------------------- | :------------------- | :--------------- |
| **Platforms**               | Linux, macOS         | macOS only       |
| **Engine**                  | Native Go (TTY)      | Electron         |
| **Startup**                 | Near-zero overhead   | Low overhead     |
| **UI**                      | Inline overlay       | GUI popover      |
| **Remote SSH**              | TTY-native, portable | macOS GUI-bound  |
| **Tmux**                    | ✓                    | Limited          |
| **Linux virtual terminals** | ✓                    | -                |
| **Memory**                  | Lightweight          | Electron runtime |

## Why not shell autocomplete plugins?

Shell plugins are great, but they also come with trade-offs. And also, not everyone use Zsh or Fish especially on SSH.

| Feature                     | IRIS                    | Shell plugins                |
| :-------------------------- | :---------------------- | :--------------------------- |
| **Installation**            | Single binary           | Plugin manager required      |
| **Shell support**           | Most shells supported   | Usually shell-specific       |
| **Startup**                 | No shell initialization | Increases shell startup time |
| **SSH**                     | One config              | Per-shell config             |
| **Tmux**                    | ✓                       | Depends on the shell         |
| **Linux virtual terminals** | ✓                       | Depends on the shell         |

## Install

#### Dependencies

- OS: Linux or macOS
- Terminal emulator with ANSI color support
- Go 1.24 or newer (if building from source)

> [!WARNING]
> Currently, Windows is not supported

#### Method 1: Package managers

<details>
<summary><b>Arch Linux (AUR)</b></summary>
<br>

IRIS is available on the Arch User Repository. You can install it using your favorite AUR helper:

```bash
yay -S iris-autocomplete
```

OR

```bash
paru -S iris-autocomplete
```
</details>

<details>
<summary><b>macOS / Linux (Homebrew)</b></summary>
<br>

Install IRIS via Homebrew tap:

```bash
brew install versenilvis/iris/iris
```
</details>

<details>
<summary><b>Nix Flakes</b></summary>
<br>

If you are using Nix Flakes, you can consume this module directly without building it manually.

1. In your `flake.nix` inputs, add:

   ```nix
   iris.url = "github:versenilvis/iris/main";
   ```

2. Then, use one of the following options to add IRIS to your system:

   **Option A: Try without installing (ephemeral)**
   ```bash
   nix run github:versenilvis/iris
   ```

   **Option B: Install to your user profile**
   ```bash
   nix profile add github:versenilvis/iris
   ```

   **Option C: Using Home Manager**
   *(Ensure you pass `inputs` to modules via `extraSpecialArgs`)*
   ```nix
   home.packages = [ inputs.iris.packages.${pkgs.system}.default ];
   ```

   **Option D: Without Home Manager (NixOS System)**
   *(Ensure you pass `inputs` to modules via `specialArgs`)*
   ```nix
   environment.systemPackages = [ inputs.iris.packages.${pkgs.system}.default ];
   ```
</details>

<details>
<summary><b>Debian / Ubuntu (.deb)</b></summary>
<br>

```bash
curl -sLO https://github.com/versenilvis/iris/releases/latest/download/iris_linux_amd64.deb
sudo dpkg -i iris_linux_amd64.deb
rm iris_linux_amd64.deb
```
*(For ARM64 architecture, replace `amd64` with `arm64`)*
</details>

<details>
<summary><b>Fedora / RHEL (.rpm)</b></summary>
<br>

```bash
curl -sLO https://github.com/versenilvis/iris/releases/latest/download/iris_linux_amd64.rpm
sudo rpm -i iris_linux_amd64.rpm
rm iris_linux_amd64.rpm
```
*(For ARM64 architecture, replace `amd64` with `arm64`)*
</details>

<details>
<summary><b>Aqua</b></summary>
<br>

If you use [aqua](https://aquaproj.github.io/), you can install IRIS by adding it to your `aqua.yaml`:

```yaml
packages:
  - name: versenilvis/iris
```

Then run `aqua i`.
</details>

<details>
<summary><b>asdf</b></summary>
<br>

If you use [asdf](https://asdf-vm.com/), you can install IRIS via its plugin:

```bash
asdf plugin add iris https://github.com/versenilvis/asdf-iris.git
asdf install iris latest
asdf set -u iris latest
asdf current iris
```
</details>

#### Method 2: Install script (recommended)

```bash
curl -sSL https://raw.githubusercontent.com/versenilvis/iris/main/scripts/install.sh | sh
```

#### Method 3: Go install

```bash
go install github.com/versenilvis/iris/cmd/iris@latest
```

#### Method 4: Build from source (for developers)

```bash
git clone https://github.com/versenilvis/iris.git
cd iris
just reload
```

## Update

```bash
iris update
```

Checks for and installs the latest release. To see what's changed first:

```bash
iris changelog       # latest release
iris changelog -n 3  # last 3 releases
```

IRIS can also check for and install updates on its own - see `updater.auto-update` in the [configuration guide](#configuration-guide) below.

## Uninstall

To completely uninstall IRIS, remove all configurations, and clean up your shell integration files, simply run:

```bash
iris uninstall
```

## Shell setup

> [!WARNING]
> **IRIS may cause visual conflicts and keybinding overlaps with other shell autosuggestion plugins or third-party completion tools. To prevent this, please disable them safely (e.g., zsh-autosuggestions, zsh-autocomplete, atuin, flyline, ...)**

> [!NOTE]
> **This is already added from installation script, but if your shell config is missing it, please add manually**

To automatically start IRIS every time you open a new shell session, add the init command to your shell config:

**Zsh (`~/.zshrc`):**
```zsh
eval "$(iris init zsh)"
```

**Bash (`~/.bashrc`):**
```bash
eval "$(iris init bash)"
```

**Fish (`~/.config/fish/config.fish`):**
```fish
iris init fish | source
```

## Configuration guide

IRIS uses a clean TOML configuration file located at `~/.config/iris/config.toml`

### Creating & viewing config

```bash
iris config init
iris config show
```

### Sample `config.toml`

```toml
[core]
version = 1                    # config schema version
shell = ""                     # "zsh", "bash", "fish", or empty for auto-detect
shell-login = false            # run shell as a login shell (also: iris --shell-login)
mode = "last"                  # "last", "spec", or "history"
debug = false                  # verbose logging to iris.log (also: iris -d)
expand-alias = true            # expand aliases before matching
auto-execute = false           # run suggestion immediately instead of inserting it
atuin-history = 0              # 0 = shell history, 1 = atuin, 2 = both
atuin-db-path = ""             # path to atuin's history.db, empty = use default
cobra-probe-enabled = true     # fall back to probing cobra binaries for completions

[ui]
style = "modern"       # "modern" or "classic"
ghost-text = true      # faded inline preview of top suggestion
hidden-files = false   # include dotfiles in suggestions
max-suggestions = 100  # max suggestions ranked before display
max-height = 15        # max visible rows in the menu
max-width = 0          # max menu width, 0 = auto
nerd-fonts = true      # use nerd-font icons

[keybindings]
toggle-mode = "ctrl+r"     # switch spec/history mode
toggle-menu = "shift+tab"  # show or hide the suggestion menu
select = "tab"             # accept the selected suggestion
navigate-up = "up"         # select up, or open history if empty
navigate-down = "down"     # select down, or open history if empty
navigate-right = "right"   # accept ghost text

[git]
filter-active-branch = true  # exclude current branch from suggestions
deduplicate-branches = true  # merge same-name local/remote branches

[updater]
check-on-startup = true  # check for updates on startup
channel = "stable"       # "stable" or "nightly"
check-interval = "24h"   # min time between update checks
auto-update = 0          # 0 = off, 1 = auto-install, 2 = confirm first

[ai]
enabled = false
provider = "groq"  # "groq" or "ollama"
debounce_ms = 400  # ms to wait before querying the AI

# please use free subscription, that is enough for your daily usage
[ai.providers.groq]
endpoint = "https://api.groq.com/openai/v1/chat/completions"
api_key_env = "GROQ_API_KEY"       # or set api_key directly
model = "llama-3.3-70b-versatile"
timeout_ms = 3000                  # ms before giving up

[ai.providers.ollama]
endpoint = "http://localhost:11434/v1/chat/completions"
model = "qwen2.5-coder"
timeout_ms = 5000  # ms before giving up
```

> [!NOTE]
> Using `api_key_env` is recommended over hardcoding `api_key` in plain text to keep credentials out of configuration files.


## Default shortcuts

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

## Theme

<div align="center">
  <img width="1920" height="1080" alt="Kitty terminal showcase" src="https://github.com/user-attachments/assets/d45cf36e-6d7d-437a-9af9-28cd994bf55f" />
  😺 <i>Kitty terminal</i>
</div>

> [!TIP]
> **We keep all available theme templates in the [`themes/`](./themes) directory**  
> **Feel free to create your own theme or contribute a new color scheme by adding it to this directory**

IRIS has theme TOML configuration file located at `~/.config/iris/theme.toml`

### Creating theme config

```bash
iris theme init
```

### Default theme
IRIS automatically falls back to the default theme if `theme.toml` is missing, empty, or contains missing configuration options

```toml
border = "#a277ff"
accent = "#61ffca"
muted = "#6d6a7f"
text = "#edecee"
text_sel = "#ffffff"
key = "#a277ff"
match = "#61ffca"
desc = "#9692a8"
desc_sel = "#edecee"
sel_bg = "#3d375e"
sel_text = "#110f18"
scroll_info = "#a277ff"
ghost_text = "#4B4A4C"
sys = "#1e1d28"
sys_sel = "#a277ff"
hist = "#1a2d36"
hist_sel = "#61ffca"
alias = "#2a2342"
alias_sel = "#a277ff"
```

> [!NOTE]
> IRIS also has 2 basic styles

<table>
  <tr>
    <td align="center"><b>Modern style</b></td>
    <td align="center"><b>Classic style</b></td>
  </tr>
  <tr>
    <td><img width="100%" alt="spec" src="https://github.com/user-attachments/assets/cb92cc64-a08d-43ad-a412-22ab479e53aa" /></td>
    <td><img width="100%" alt="spec" src="https://github.com/user-attachments/assets/6e68b330-39d9-4750-ac95-e94514ba4e7b" /></td>
  </tr>
  <tr>
    <td><img width="100%" alt="history" src="https://github.com/user-attachments/assets/fd1a272a-19fe-472a-83d6-eec790813403" /></td>
    <td><img width="100%" alt="history" src="https://github.com/user-attachments/assets/a1983d05-b771-4277-91cf-96009875979b" /></td>
  </tr>
</table>

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

If IRIS crashes, it will automatically save a crash log and show the path on your terminal (`~/.iris/crash.log`). Or you can find the path to the latest crash log by running:
```bash
iris crash-log
```
Please include this file when reporting a crash.
 
## Developer documentation

For system architecture overview, engine design, and contribution guide, please refer to the [Developer documentation](./docs/dev/README.md).

## License

This project is licensed under the [0BSD License](LICENSE) - no strings attached. Meaning you can do whatever you want with it.

For those who fork it and want to publish a new version or something else; if you can, a credit or co-author mention is always welcome :) (though never required).

Thank you!

## Feedback

I'd love to hear your feedback

Feel free to reach out via:

- [Email](mailto:versedev.store@proton.me)
- [Twitter](https://twitter.com/versenilvis)
- [GitHub issues](https://github.com/versenilvis/iris/issues/new)
