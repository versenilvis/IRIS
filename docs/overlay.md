# Overlay Rendering UI (`integration/overlay.go`)

The `overlay` package handles the visual representation of suggestions. It is designed to be "non-destructive," meaning it draws over the shell without corrupting the prompt or the scrollback buffer.

## Design Philosophy

- **ANSI ESC Everywhere**: It uses CSI (Command Sequence Introducer) codes to manipulate the terminal.
- **Save/Restore**: It uses `\0337` (DECSC) and `\0338` (DECRC) to jump back to the prompt after drawing the menu.
- **Fixed Width**: The box is exactly 72 characters wide (`boxWidth`) to ensure a consistent, premium feel.

## Technical Details

### Terminal Scrolling Protection
When rendering near the bottom of the screen, simply printing lines would cause the terminal to scroll and bury the prompt. Iris solves this by:
1. Moving to the prompt.
2. Printing N empty newlines.
3. Moving back up N lines.
4. Saving the cursor *again* at this new, stabilized location.

### Styling with Lipgloss
While Iris handles the positioning via raw ANSI codes, it uses **Lipgloss** for the "interior design":
- **Colors**: Dracula-inspired palette (`#BD93F9`, `#6272A4`).
- **Icons**: Fixed-width columns for badges.
- **Selection**: Highlights the active item with a background color (`#44475A`).

## Example Logic
To draw line 1 of the menu:
1. `\0338` (Jump to prompt anchor).
2. `\033[2B` (Move 2 lines down).
3. `\033[K` (Clear the entire line).
4. Print `│ [Icon] Command name... │`.
