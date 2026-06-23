<div align="center">
  <img width="50%" alt="banner" src="https://github.com/user-attachments/assets/c5ec623b-8259-473f-b7c3-3d01a64deb5d" />
  <!-- <h1>IRIS</h1> -->
  <p>IRIS (Intelligent Real-time Input Suggestion) - A shell auto-completion tool that works like code editor's IntelliSense.</p>
  
  [![GitHub Actions](https://img.shields.io/github/actions/workflow/status/versenilvis/IRIS/release.yml?branch=main&style=for-the-badge&logo=github&logoColor=white&label=Actions)](https://github.com/versenilvis/IRIS/actions/workflows/release.yml)
  [![License: 0BSD](https://img.shields.io/badge/License-0BSD-blue?style=for-the-badge&logo=github&logoColor=white)](./LICENSE)

</div>


<!--
  [![Status](https://img.shields.io/badge/status-beta-yellow?style=for-the-badge&logo=github&logoColor=white)]()
  [![Documentation](https://img.shields.io/badge/docs-available-brightgreen?style=for-the-badge&logo=github&logoColor=white)](./docs)
  [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen?style=for-the-badge&logo=github&logoColor=white)](./.github/CONTRIBUTING.md)
-->



## Installation

```bash
curl -sSL https://raw.githubusercontent.com/versenilvis/iris/main/scripts/install.sh | sh
```

## Shortcuts

| Shortcut                           | Action                  | Description                                                               |
| :--------------------------------- | :---------------------- | :------------------------------------------------------------------------ |
| <kbd>Shift</kbd> + <kbd>Tab</kbd>  | Toggle menu             | Show or hide the suggestion menu.                                         |
| <kbd>Esc</kbd>                     | Hide menu               | Temporarily hide the menu until the next key press.                       |
| <kbd>Tab</kbd>                     | Accept suggestion       | Insert the currently selected suggestion into the prompt.                 |
| <kbd>Enter</kbd>                   | Execute command         | Close the menu and send the current command to the shell.                 |
| <kbd>↑</kbd>                       | Navigate up / history   | Move the selection up, or open command history when the prompt is empty.  |
| <kbd>↓</kbd>                       | Navigate down / history | Move the selection down, or open command history when the prompt is empty. |
| <kbd>→</kbd>                       | Accept ghost text       | Accept the faded ghost text suggestion when the menu is open.             |
| <kbd>←</kbd> / <kbd>→</kbd>        | Move cursor             | Move the cursor inside the input buffer. Disabled when the prompt is empty. |
| <kbd>Ctrl</kbd> + <kbd>R</kbd>     | Switch mode             | Toggle between `spec` and `history` mode.                                 |
| <kbd>Ctrl</kbd> + <kbd>A</kbd>     | Beginning of line       | Move the cursor to the start of the command line.                         |
| <kbd>Ctrl</kbd> + <kbd>E</kbd>     | End of line             | Move the cursor to the end of the command line.                           |
| <kbd>Ctrl</kbd> + <kbd>L</kbd>     | Clear screen            | Clear the terminal while preserving the input buffer and redrawing the menu. |
| <kbd>Ctrl</kbd> + <kbd>U</kbd>     | Clear command           | Remove the entire current command and close the menu.                     |
| <kbd>Ctrl</kbd> + <kbd>C</kbd>     | Cancel command          | Send `SIGINT`, clear the input buffer, and close the menu.                |
| <kbd>Ctrl</kbd> + <kbd>W</kbd>     | Delete word             | Delete the word immediately before the cursor.                            |

## License

This project is licensed under the [0BSD License](LICENSE) - no strings attached. Meaning you can do whatever you want with it.

For those who fork it and want to publish a new version or something else; if you can, a credit or co-author mention is always welcome :) (though never required).

Thank you!
