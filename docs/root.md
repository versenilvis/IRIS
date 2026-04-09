# IRIS: Central Integration & Event Loop (`root/`)

The `root` package is the entry point and the orchestration layer of IRIS. It handles the low-level terminal manipulation and the main interaction loop.

## How it works

1. **PTY Wrapper**: When you run `iris`, it starts a Pseudo-Terminal (PTY) and launches `bash` inside it.
2. **IO Interception**: It creates two "pumps":
   - **Output Pump**: Forwards everything from Bash to your screen using `TermWrite` (synchronized to prevent UI glitches).
   - **Input Pump**: Listens to your keyboard. If it detects you are typing a command, it records it in a `naiveBuffer` and triggers the suggestion engine.
3. **State Management**: It tracks whether you are in `spec` mode (command suggestions) or `history` mode (Ctrl+R).

## Key Components

### `root.go`
Contains the `runWrapper()` function which:
- Sets the terminal to **Raw Mode** so Iris can capture keys like `Tab`, `Esc`, or `Ctrl+C` before the shell does.
- Manages the `naiveBuffer`: a string that tracks exactly what you see on your prompt.
- Handles **Selection Logic**: When you press `Tab`, it modifies the Bash line using the `Ctrl+U` (clear line) + `selected command` sequence.

### `term_sync.go`
Provides `TermWrite`, a thread-safe wrapper around `os.Stdout`. It uses a `sync.Mutex` to ensure that if Bash and the Iris Overlay try to write at the exact same millisecond, the output doesn't get garbled.

## Example Flow
1. User types `g`.
2. `root` captures `g`, appends to `naiveBuffer`.
3. `root` calls `renderOverlay()`.
4. `renderOverlay` calls `Lookup("g")`.
5. `overlay` renders the result `git`.
6. User presses `Tab`.
7. `root` sends `Ctrl+U` to Bash, then sends `git ` (the completion).
