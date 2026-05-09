# Shell History Integration (`integration/history.go`)

The `history` module provides the Ctrl+R fuzzy search functionality by reading the user's persistent shell command history.

## Features

- **Zsh Extended Support**: Specifically parses the `: <timestamp>;<command>` format common in Zsh.
- **Lazy Loading**: Doesn't read the disk until the user actually requests history. This keeps startup time instantaneous.
- **Fuzzy Search**: Integrated with the `fuzzy` search engine.
- **Deduplication**: Automatically hides duplicate entries, showing only the most unique command variants.

## Data Flow

1. User presses `Ctrl+R`.
2. `root` sets `mode = "history"`.
3. `SearchHistory("")` is called.
4. If `cache` is empty:
   - Reads `~/.zsh_history`.
   - Strips metadata using `;` delimiter.
   - Populates a slice of commands.
5. Search matches are sorted by the `fuzzy` engine.
6. Suggestions are returned to the `overlay` for rendering.

## Configuration
It currently looks for `.zsh_history` in the user's home directory.
