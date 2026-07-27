# Overlay rendering UI (`integration/overlay.go`)

The `overlay` package handles the visual representation of suggestions. It is designed to be non-destructive, drawing over the shell without corrupting the prompt or scrollback buffer.

## Design philosophy

- **ANSI ESC everywhere**: Uses ANSI escape sequences and CSI (Control Sequence Introducer) codes to manipulate terminal positioning.
- **Save and restore**: Uses DECSC (`\0337`) and DECRC (`\0338`) escape sequences to return to prompt anchor after drawing.
- **Fixed width**: Menu box uses a standard 72-character width (`boxWidth`) for consistent layout rendering.

## Technical details

The overlay engine manages positioning, scrolling protection, and styled cell drawing directly on the terminal grid.

### Terminal scrolling protection

When rendering near the bottom of the screen, printing new lines causes terminal scrolling which buries the active prompt. Iris prevents this by:

1. Moving to the prompt anchor.
2. Printing N empty newlines to pre-allocate scroll space.
3. Moving back up N lines.
4. Re-saving the cursor at this stabilized location.

### Styling with Lipgloss

While Iris handles cursor positioning via raw ANSI codes, visual component styling uses Lipgloss:

- **Colors**: Dracula-inspired color theme (`#BD93F9`, `#6272A4`).
- **Icons**: Fixed-width columns for category badges and completion types.
- **Selection**: Highlights the currently active suggestion item with a distinct background (`#44475A`).

## Example logic

To render line 1 of the menu dropdown:

1. `\0338` (Jump to prompt anchor).
2. `\033[2B` (Move 2 lines down).
3. `\033[K` (Clear line buffer).
4. Print formatted item: `│ [Icon] Command name... │`.
