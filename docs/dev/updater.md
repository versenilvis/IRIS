# Versioning and updates

Iris has a built-in update notification system designed to be non-intrusive and zero-latency

The update system is designed to keep Iris up to date while staying out of the way. It performs an asynchronous network check when the shell starts and notifies you exactly once per version after you run a command. This prevents redundant notifications and ensures that you only see an update message when a new release is actually available. Optionally, it can also install updates for you in the background. The system also includes dedicated debugging tools to verify the notification and auto-update logic without requiring real releases

## How it works

1. **Background check**: every time you open a new terminal, Iris launches a background goroutine to check for updates
   - checks the network at most once every `check-interval` (24h by default) to avoid GitHub API rate limiting
   - has a 5 second timeout so it never hangs on a slow or missing network connection
   - if there is no network, it fails silently with zero impact on startup

2. **State persistence**: state is stored in `$XDG_DATA_HOME/iris/state.toml` (or `~/.local/share/iris/state.toml`), including:
   - `updater.last-check-time` / `updater.seen-version` - when the last check ran, and the latest version already notified about
   - `updater.auto-update-target` / `updater.auto-update-attempt` - the auto-update escalation ladder's state (see [Auto-update](#auto-update) below)
   - `updater.declined-version` - a version the user explicitly declined at an auto-update confirm prompt, never re-prompted

3. **Smart notification**:
   - the notice only appears after you run your first command (triggered by the `IRIS_CMD_STOP` IPC signal from the shell hook)
   - appears only once per session, never again even if you keep the terminal open
   - if you have already been notified about a specific version, it will not show again until a newer GitHub release tag is detected
   - when `iris update` is run successfully, the `seen-version` flag is cleared
   - includes up to two changelog bullets pulled from the release's grouped notes (see [Implementation details](#implementation-details) for where those come from)

## Auto-update

`updater.auto-update` (see [Configuration](#configuration)) controls whether the background check can install updates itself, not just notify:

- **`0` (default)**: notify only, exactly the behavior above.
- **`1`**: auto-install, with a per-version escalation ladder so a flaky network or a failing installer can't turn into a silent retry loop:
  1. first detection of a newer version: installs it in the background, no prompt. Output is captured (not streamed) since the wrapper's terminal is in raw mode - a clean one-line notice is printed once the install finishes, if you're still in the session.
  2. if that version is still current on the *next* check (the install didn't stick), it stops auto-installing and shows a confirm prompt instead: `[y/N]`, pressed like any other single keystroke.
  3. if it's still not current after that, Iris gives up on that version - prints one notice pointing at `iris update`, and won't try again until a *newer* version is released.
- **`2`**: always confirm first, every time - never installs silently, never gives up, just keeps asking until you say yes or the version changes.

A confirm prompt declined with `n`/Esc/Ctrl+C records that version so it's never re-prompted (until a newer one comes along).

## Build-time versioning

Iris uses Go `ldflags` to inject the version string at build time:

```bash
go build -ldflags="-X github.com/versenilvis/iris/root.Version=v1.2.0" -o iris main.go
```

If not provided, the version defaults to dev. The dev version will never trigger an update notification

## Configuration

The `[updater]` section of `~/.config/iris/config.toml`:

```toml
[updater]
check-on-startup = true
channel = "stable" # "stable" or "nightly"
check-interval = "24h"
auto-update = 0 # 0 = off, 1 = auto-install, 2 = always confirm first
```

## Commands

- `iris version` - print the current version of the running binary
- `iris update` - manually check for and apply the latest release
- `iris changelog` - show grouped release notes fetched from GitHub (`-n <count>`, `--refresh`)

## Debugging and testing

Four separate things to test: the update command, the in-session notification banner, the changelog command, and the background auto-updater

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
  - <first changelog bullet>
```

### Test 3: changelog command

Tests fetching, grouping, the response cache, and rate-limit fallback. No full Iris session needed

```bash
just build-release v0.0.1
just debug-changelog v1.99.0
```

### Test 4: background auto-updater

Tests the escalation ladder (silent install, then escalating to a confirm-worthy state, then giving up rather than looping) fully offline, by inspecting `state.toml` after each simulated check cycle. The confirm-prompt UI itself is interactive and isn't exercised here - set `auto-update` in a real config and watch a live session with shell hooks active to see it rendered

```bash
just debug-autoupdate
```

### Environment variables

| Variable                   | Purpose                                                                                              |
| --------------------------- | -------------------------------------------------------------------------------------------------------- |
| `IRIS_UPDATE_URL`          | override the release-check endpoint with a custom URL (used by `debug-update`, `debug-autoupdate`)  |
| `IRIS_MOCK_LATEST_VERSION` | skip network entirely, resolve to this version immediately (used by `debug-notify`)                 |
| `IRIS_CHANGELOG_URL`       | override the changelog-fetch endpoint with a custom URL (used by `debug-changelog`)                 |
| `IRIS_INSTALL_URL`         | override the install-script fetch URL (used by `debug-autoupdate`)                                  |

### State reset

To force a fresh check on next launch, delete the state file:

```bash
rm ~/.local/share/iris/state.toml
```

## Implementation details

| File                    | Purpose                                                                                     |
| ------------------------ | ----------------------------------------------------------------------------------------------- |
| `root/version.go`       | holds the `Version` variable, defaults to `dev`                                             |
| `root/update.go`        | version fetch/compare, state persistence, the `iris update`/`iris version` commands         |
| `root/changelog_cmd.go` | the `iris changelog` command: fetch, cache, grouping-aware rendering                         |
| `root/autoupdate.go`    | `performUpdate`, the auto-update escalation ladder, and the confirm-prompt keystroke handler |
| `root/wrapper.go`       | wires the background check into the IPC loop, prints notices on `IRIS_CMD_STOP`             |
| `.goreleaser.yaml`      | `changelog:` block groups commits by type for the GitHub release body                        |
| `cliff.toml`            | generates the committed `CHANGELOG.md` from the same commit history on each stable release   |
