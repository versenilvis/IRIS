# Iris central integration & event loop (`root/`)

The `root` package is the entry point and orchestration layer of Iris, handling low-level terminal PTY manipulation and the main interaction loop.

## How it works

1. **PTY wrapper**: Starts a pseudoterminal (PTY) and launches target shell (`zsh`, `bash`, `fish`) inside it.
2. **IO interception**: Runs two pumps:
   - **Output pump**: Forwards shell output to terminal screen via synchronized `TermWrite`.
   - **Input pump**: Listens to keystrokes in raw mode, tracks typed characters in `naiveBuffer`, and triggers suggestion rendering.
3. **State management**: Tracks active completion mode (`spec` vs `history`).

## Key components

### `root/wrapper.go`
Contains the core PTY loop:
- Sets terminal to raw mode to intercept keys like `Tab`, `Esc`, or `Ctrl+C`.
- Manages `naiveBuffer` string tracking prompt input state.
- Handles suggestion insertion when pressing `Tab` or `Enter`.

### `root/term_sync.go`
Provides `TermWrite`, a thread-safe stdout wrapper using `sync.Mutex` to prevent screen garbling when shell output and overlay rendering overlap.

## Example flow

1. User types `g`.
2. `root` captures `g`, appends to `naiveBuffer`.
3. `root` calls `renderOverlay()`.
4. `renderOverlay` calls `Lookup("g")`.
5. `overlay` renders suggestions box.
6. User presses `Tab`.
7. `root` inserts completion into prompt buffer.

## Hot-reload

Iris includes an atomic hot-reload mechanism for rapid local development:

- **Signal listener**: The root process listens for `SIGUSR1`.
- **Process replace**: Uses `syscall.Exec` to replace the running binary in-place without killing the underlying PTY shell session.
