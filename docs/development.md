# Development guide

This document outlines the architectural principles, development workflow, and mandatory coding conventions for contributors working on the Iris codebase

## Architecture overview

Iris functions as a transparent pseudoterminal (PTY) bridge sitting between the user terminal emulator and the underlying interactive shell (`zsh`, `bash`, or `fish`)

- PTY Bridge (`root/wrapper.go`): Intercepts keystrokes in raw mode, tracks typed queries, manages shell input/output buffers, and handles escape sequences
- Overlay Rendering (`integration/overlay.go`): Computes visual character widths (`lipgloss.Width`), handles cursor column positioning, and draws floating menu boxes directly on the terminal grid
- Suggestion Engine (`root/suggestions.go`): Collects, deduplicates, and sorts command completions from static specifications and fuzzy history search

## Development workflow

### Setting up environment

Clone the repository and verify that all dependencies are synced:

```bash
git clone https://github.com/versenilvis/iris.git
cd iris
go mod tidy
```

### Building and reloading

To rapidly compile and reload your local development binary, ensure `just` is installed and run:

```bash
just reload
```

### Running tests

Run the full automated test suite without cache to verify system integrity before submitting code:

```bash
go test -v -count=1 ./...
```

OR

```bash
just test
```

OR using my own test analyzer script
```bash
just ana
```
### Running linters

Ensure zero linting errors across the codebase:

```bash
golangci-lint run ./...
go vet ./...
```

OR

```bash
just lint
```

### Some available `just` commands

We use `just` as our primary command runner. Below are common shortcuts available during development:

- `just build`: Compile standard local binary
- `just optimized-build`: Compile stripped, optimized binary
- `just reload`: Compile binary and hot-reload any currently running Iris session
- `just run`: Execute local `./iris` binary
- `just debug`: Launch `./iris -d` in debug logging mode
- `just test`: Execute full Go test suite
- `just ana`: Run custom project test analyzer script
- `just lint`: Execute static analysis and linter checks
- `just build-release <version>`: Build versioned release binary (e.g. `just build-release v1.2.0`)

## Mandatory engineering rules

Contributors must strictly adhere to the following core UI/UX and architectural principles:

### Absolute positioning for overlay rendering

- Never rely on relative cursor movements (`\033[nA` or `\033[nD`) across multi-line inputs to draw menus, as wrapping causes severe visual jitter
- Always compute absolute target columns (`targetCol = PromptLen + typedLen`) and use carriage return (`\r`) followed by horizontal positioning (`\033[targetColC`) to anchor the overlay box

### UI rules

- Stationary navigation vs dynamic typing: When navigating up or down through the menu (`↑` or `↓`), the overlay box position must remain fixed in place. When typing characters, the menu must shift dynamically along with the cursor column
- Menu hiding: When clearing the input line completely or when typing up to a completed token where no further completion exists, the overlay menu must automatically hide
- Context preservation: When moving left or right with arrow keys (`←` or `→`) or deleting characters with backspace, the engine must preserve context and re-evaluate the exact substring position to generate accurate suggestions
- Empty buffer arrow suppression: When the command line buffer is entirely empty, pressing left or right arrow keys (`←` or `→`) must not open the menu or trigger suggestions
- Temporary vs permanent dismissal: Pressing `Esc` hides the overlay menu temporarily for the current command line, whereas pressing `Shift+Tab` disables overlay suggestions permanently until toggled back on
- Non-blocking navigation and context awareness: When moving the cursor back to the start of a typed command line, cursor movement must never be blocked or frozen, and the overlay menu must dynamically preserve context relative to the cursor position

### Responsive PTY handling and non-Blocking UI

- Do not block raw PTY keystroke interception loops with long operations or arbitrary debounce timers that freeze user typing
- When handling navigation keys (up/down arrow) or mode toggles (`Ctrl+R`), update internal state synchronously and render UI feedback immediately

### Mode switching integrity

- When intercepting mode toggle keys (`Ctrl+R`), always reset navigation flags (`userNavigated = false`) and restore the original typed query buffer
- This ensures that switching from specification suggestions to history navigation reloads clean completion lists starting from index 0

### Concurrency safety

- The PTY output read loop and keystroke interception loop run asynchronously in separate goroutines
- Always use explicit synchronization (`sync.Mutex` or `sync.RWMutex`) when reading or modifying shared state variables such as `naiveBuffer`, `activeMode`, or `overlay.Items`
