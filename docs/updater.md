# Versioning and updates

Iris has a built-in update notification system designed to be non-intrusive and zero-latency

The update system is designed to keep Iris up to date while staying out of the way. It performs an asynchronous network check when the shell starts and notifies you exactly once per version after you run a command. This prevents redundant notifications and ensures that you only see an update message when a new release is actually available. The system also includes dedicated debugging tools to verify the notification logic without requiring real releases

## How it works

1. **Background check**: every time you open a new terminal, Iris launches a background goroutine to check for updates
   - checks the network at most once every 6 hours to avoid GitHub API rate limiting
   - has a 5 second timeout so it never hangs on a slow or missing network connection
   - if there is no network, it fails silently with zero impact on startup

2. **State persistence**: state is stored in `~/.iris/update_state.json` with two fields:
   - `last_check` - unix timestamp of the last network check
   - `seen_version` - the latest version the user was already notified about

3. **Smart notification**:
   - the notice only appears after you run your first command (triggered by the `IRIS_CMD_STOP` IPC signal from the shell hook)
   - appears only once per session, never again even if you keep the terminal open
   - if you have already been notified about a specific version, it will not show again until a newer GitHub release tag is detected
   - when `iris update` is run successfully, the `seen_version` flag is cleared

## Build-time versioning

Iris uses Go `ldflags` to inject the version string at build time:

```bash
go build -ldflags="-X github.com/versenilvis/iris/root.Version=v1.2.0" -o iris main.go
```

If not provided, the version defaults to dev. The dev version will never trigger an update notification

## Commands

- `iris version` - print the current version of the running binary
- `iris update` - manually check for and apply the latest release

## Debugging and testing

There are two separate things to test: the update command itself, and the in-session notification banner

### Test 1: update command

Tests version fetching, comparison and the output message. No full Iris session needed

```bash
just build-release v0.0.1
just debug-update v1.99.0
```

Expected output
```
--- testing iris update command ---
checking for updates (current: v0.0.1)...
[IRIS] updating v0.0.1 -> v1.99.0
running: curl -sS https://raw.githubusercontent.com/versenilvis/iris/main/scripts/install.sh | sh

[IRIS] restart your terminal to use the new version
```

### Test 2: in-session notification banner

Tests the yellow notice that appears after you run your first command inside a live Iris session. Requires `iris.zsh` to be active in the inner shell so `IRIS_CMD_STOP` fires through the IPC pipe

```bash
just build-release v0.0.1
just debug-notify v1.99.0
```

Inside the new session, run any command. Expected output after the command
```
[IRIS] new version v0.0.1 -> v1.99.0 available, run iris update to upgrade
```

### Environment variables

| Variable                   | Purpose                                                                             |
| -------------------------- | ----------------------------------------------------------------------------------- |
| `IRIS_UPDATE_URL`          | override the GitHub API endpoint with a custom URL (used by `debug-update`)         |
| `IRIS_MOCK_LATEST_VERSION` | skip network entirely, resolve to this version immediately (used by `debug-notify`) |

### State reset

To force a fresh network check on next launch, delete the state file:

```bash
rm ~/.iris/update_state.json
```

## Implementation details

| File              | Purpose                                                                                  |
| ----------------- | ---------------------------------------------------------------------------------------- |
| `root/version.go` | holds the `Version` constant, defaults to `dev`                                          |
| `root/update.go`  | all update logic: state persistence, network fetch, version compare, commands            |
| `root/wrapper.go` | wires the background check into the IPC loop, prints the notice on first `IRIS_CMD_STOP` |
