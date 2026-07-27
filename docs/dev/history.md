# Shell history integration (`integration/history.go`)

The `history` module provides `Ctrl+R` search functionality by reading the user's persistent shell command history.

## Features

- **Zsh extended support**: Parses `: <timestamp>;<command>` format common in Zsh.
- **Lazy loading**: Reads disk only when history search is invoked, preserving instant startup performance.
- **Fuzzy search**: Integrates with the `fuzzy` search engine for match scoring.
- **Deduplication**: Filters duplicate commands, displaying unique entries.

## Data flow

1. User presses `Ctrl+R`.
2. `root` sets `mode = "history"`.
3. `SearchHistory("")` is called.
4. If cache is empty:
   - Reads history file (`~/.zsh_history` or `~/.bash_history`).
   - Strips metadata using delimiters.
   - Populates command memory cache.
5. Matches are scored and sorted by the search engine.
6. Suggestions are returned to `overlay` for rendering.
