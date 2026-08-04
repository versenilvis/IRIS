# To Clankers

## IRIS - Shell Auto-Completion Tool

**IRIS** is a lightweight, TTY-native shell auto-completion tool written in Go. No Electron, no GUI, no telemetry-runs on macOS and Linux with terminal emulators supporting ANSI colors.

## Product Language

- **Core value**: Real-time command suggestions rendered directly in the terminal (like IntelliSense for shells)
- **Alternative to Fig** (sunset Sept 2024): Lighter, open-source, multi-shell (bash/zsh/fish), runs on any terminal
- **Design principle**: Native Go binary, pure TTY, one process per session
- **Two suggestion modes**: "spec" (command specs/flags) or "history" (shell history + frecency)

## Key Architectural Files

### PTY & Input Handling
- **`root/wrapper.go:124-1230`** - Core event loop: PTY setup, terminal raw mode, keystroke interception, signal handling. Dense and critical.
  - Signal handling: SIGWINCH (resize), SIGUSR1 (reload)
  - Input parsing: Enter, Tab, Shift+Tab, arrow keys, Ctrl+shortcuts
  - Output bridging: Shell stdout to real stdout
  - IPC listening: shell-to-iris communication via FD13

### Suggestion Engine
- **`root/suggestions.go:18-111`** - MergeResults: dedupes, scores, and ranks suggestions from history/spec/AI
  - Two paths: history mode (confidence-based) vs spec mode (frecency scoring)
  - Frecency: uses `internal/scoring` package with context signals (cwd, prev command, etc.)
  - AI injection: optional LLM suggestions via `internal/ai`

### Configuration & Themes
- **`internal/config/config.go`** - Config struct, TOML parsing, hot-reload on file change
  - Default keybindings: Ctrl+R (toggle mode), Shift+Tab (toggle menu), Tab (select)
  - Validation rules: modes (last/spec/history), shells (bash/zsh/fish), channels (stable/nightly)
  - AutoDetectConfigChange: polls every 1s for config/theme changes

### Terminal Rendering
- **`integration/overlay.go`** - Overlay struct, menu rendering, ghost text, cursor tracking
  - ComputeCursorCol: parses ANSI escape sequences to track cursor position
  - Renders using lipgloss v2 (charm.sh)

### Data Structures
- **`spec/spec.go:1-59`** - Suggestion, Spec, Subcommand, Option types
  - Suggestion: Cmd, Desc, Icon, Source ("history"/"spec"/"ai"), Confidence (0-100), Priority
  - Registry: global map of spec definitions

### Shell Integration
- **`integration/shell/`** - Adapters for bash/zsh/fish
  - GetEnv: constructs shell environment with FD13 for IPC
  - RecordSessionCommand: captures executed commands

### Process Lifecycle
- **`root/root.go:81-192`** - Watchdog parent process
  - Monitors child stderr for panics/crashes
  - Restores terminal state on crash
  - Triggers rescue shell and logs to `~/.iris/crash.log`

### Commands
- **`root/config_cmd.go`** - `iris config init/show`
- **`root/theme_cmd.go`** - `iris theme init`
- **`root/uninstall.go`** - `iris uninstall` (cleanup)
- **`root/update.go`** - Background version check

## Common Flows

### User Types Character
1. Shell process sends query via FD13 to iris
2. `wrapper.go:480-582` receives on IPC scanner
3. `MergeResults` scores suggestions (history or spec mode)
4. `overlay.Render()` outputs menu to stdout
5. Terminal displays overlay above prompt

### User Presses Tab (Select Suggestion)
1. `wrapper.go:813-843` intercepts Tab
2. Gets current selection from overlay
3. If spec mode: adds trailing space
4. Writes to PTY (`0x15` = Ctrl+U clears, then new text)
5. Clears overlay, triggers render

### User Presses Ctrl+R (Toggle Mode)
1. `wrapper.go:773-795` intercepts Ctrl+R
2. Swaps `activeMode` between "spec" and "history"
3. Saves mode to state file for next session
4. Re-renders menu with new scoring

### Terminal Resizes (SIGWINCH)
1. `wrapper.go:213-215` catches signal
2. `pty.InheritSize` syncs PTY to new size
3. Overlay re-renders with new width

### Config Reloads
1. User edits `~/.config/iris/config.toml`
2. `config.AutoDetectConfigChange` polls, detects mtime change
3. Calls `config.Load()`, updates global config
4. Callback triggers overlay re-render if visible

## Design Patterns

- **Mutex-protected state**: naiveBuffer, cursorOffset, activeMode (avoid races in goroutines)
- **Atomic flags**: isCommandActive, userNavigated, disableGhostText (fast checks without locks)
- **Debounced rendering**: 20ms timer on keystroke to batch render calls
- **AI suggestions**: separate debounce (400ms default), cancellable context
- **Ghost text**: faded hint of next completion, cleared on user navigation
- **Frecency scoring**: combines frequency + recency + context signals (cwd, prev command, exit code)

## Comments
- NO FUNTION COMMENTS UNLESS IT'S TOO COMPLEX TO UNDERSTAND
- NO TOP FILE COMMENTS
- NO COMMENT LONGER THAN 2 LINES UNLESS ASKED EXPLICITLY
- ONLY COMMENT WHERE IT'S IMPORTANT TO UNDERSTAND
- THE COMMENT SHOULD EXPLAIN "WHY" OR "WHY WE NEED THIS" INSTEAD OF TELLING "WHAT IT DOES"

## Configuration Paths

- Config: `~/.config/iris/config.toml`
- Theme: `~/.config/iris/theme.toml`
- State: `~/.local/share/iris/state.toml` (or XDG_DATA_HOME)
- Logs: `~/.local/share/iris/iris.log` (or XDG_CACHE_HOME)
- Crash: `~/.iris/crash.log`

## Testing

- Test files match source: `*_test.go` (unit tests, no integration tests yet)
- Mocks for git provider in `root/mock_git_provider_test.go`
- `spec/cobra_complete_test.go` tests command completion

## Verify
- Run `just test` or `go test ./... -v` to make sure all tests passed
- Run `just lint` or `golangci-lint run ./...` to make sure zero linter problems
- Run `go vet ./...` to make sure there is no suspicious constructs, structural mistakes, and high-probability bugs
- Run `go fix ./...` to modernize all Go code

## CommitS
- If the agent is the one who commits, make sure follow the conventional commits style
- The commit should be as short as possible, but enough to describe changes
---

**Minimal reference for rapid navigation. Read files directly-do not fabricate code from this index.**
