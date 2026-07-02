<div align="center">
  <!-- <img width="50%" alt="banner" src="https://github.com/user-attachments/assets/c5ec623b-8259-473f-b7c3-3d01a64deb5d" /> -->
  <!-- <img width="25%" alt="logo" src="https://github.com/user-attachments/assets/79d3913c-56b7-42cb-8b07-53e98f39322b" /> -->
  <img width="15%" alt="logo" src="https://github.com/user-attachments/assets/10b7ca98-872b-44a2-bdcd-265f18aa0564" />
  
  <!-- <h1>IRIS</h1> -->
  <p>IRIS (Intelligent Real-time Input Suggestion) - A shell auto-completion tool that works like code editor's IntelliSense</p>
  
  [![GitHub Actions](https://img.shields.io/github/actions/workflow/status/versenilvis/IRIS/release.yml?branch=main&style=for-the-badge&logo=github&logoColor=white&label=Actions)](https://github.com/versenilvis/IRIS/actions/workflows/release.yml)
  [![License: 0BSD](https://img.shields.io/badge/License-0BSD-blue?style=for-the-badge&logo=github&logoColor=white)](./LICENSE)
  
  <a href="#why-iris-instead-of-fig-or-zsh-auto-plugins">Comparison</a> · <a href="#installation">Installation</a> · <a href="#shortcuts">Shortcuts</a> · <a href="#reporting-bugs">Reporting bugs</a> 

</div>

**IRIS is built on top of TTY, so it runs everywhere. It just needs a terminal!**

Run iris wherever you already work; your local machine, a remote server, or anywhere you can ssh. Each suggestion menu renders directly inline inside your real terminal session, not an app's imitation of one, so it never breaks full-screen TUIs or terminal formatting. Automatically index your aliases and shell history to suggest commands that match your actual workflow in real time. Change configurations and propagate them instantly without restarting your shell. One single local native Go binary, not an app: no gui, no electron, no mac-only wrapper, no account, no telemetry. (if you've used fig: it's that, rebuilt to run purely on TTY)

<!--
  [![Status](https://img.shields.io/badge/status-beta-yellow?style=for-the-badge&logo=github&logoColor=white)]()
  [![Documentation](https://img.shields.io/badge/docs-available-brightgreen?style=for-the-badge&logo=github&logoColor=white)](./docs)
  [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen?style=for-the-badge&logo=github&logoColor=white)](./.github/CONTRIBUTING.md)
-->


## Why Iris instead of Fig or Zsh auto-plugins

> [!IMPORTANT]
> **[Fig](https://app.fig.io/) was officially sunset in September 2024 and migrated to Amazon Q Developer (which requires cloud authentication and proprietary bloat)**  
> **IRIS is the lightweight, open-source, zero-telemetry alternative built purely on native Go and TTY with no accounts, no GUI app, and no background daemons required**

How it compares

| Feature                           | Iris                 | Fig                        | Zsh plugins           |
| :-------------------------------- | :------------------- | :------------------------- | :-------------------- |
| **Engine**                        | Native TTY binary    | Electron app               | Shell scripts         |
| **Startup**                       | No startup cost      | Minimal                    | Slower shell startup  |
| **UI**                            | Inline overlay       | GUI popover                | Inline suggestions    |
| **SSH**                           | Works out of the box | Desktop companion required | Remote setup required |
| **Tmux & Linux virtual terminal** | Full support         | Not supported              | Full support          |
| **Memory**                        | < 15 MB              | Heavy                      | Moderate              |

Supported shell & environment

| Environment                                                                      | Iris  |  Fig  |       Shell plugins       |
| :------------------------------------------------------------------------------- | :---: | :---: | :-----------------------: |
| **Zsh**                                                                          |   ✓   |   ✓   |      ✓ (Zsh plugins)      |
| **Bash**                                                                         |   ✓   |   ✓   |             -             |
| **Fish**                                                                         |   ✓   |   ✓   |      ✓ (Fish native)      |
| **Nushell**                                                                      |   ✓   |   -   |             -             |
| **POSIX sh / Dash**                                                              |   ✓   |   -   |             -             |
| **Busybox ash**                                                                  |   ✓   |   -   |             -             |
| **Tmux**                                                                         |   ✓   |   -   |     ✓ (only Zsh/Fish)     |
| **Linux virtual terminal (<kbd>Ctrl</kbd> + <kbd>Alt</kbd> + <kbd>F1-F9</kbd>)** |   ✓   |   -   |             -             |
| **SSH sessions**                                                                 |   ✓   |   -   | ✓ (requires remote setup) |

## Installation

```bash
curl -sSL https://raw.githubusercontent.com/versenilvis/iris/main/scripts/install.sh | sh
```

> [!NOTE]
> Currently, Windows is not supported

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

## Reporting bugs
> [!NOTE]
> Describing the bug you are facing, along with the relevant log  
> Enabling debug mode and then performing actions that led to the error

Run IRIS with debug mode:
```bash
iris -d
```
or `config.toml`:
```toml
debug=true
```

> [!IMPORTANT]
> **Since IRIS logs everything you type, you should only enable debug mode when you need to report bugs**

## License

This project is licensed under the [0BSD License](LICENSE) - no strings attached. Meaning you can do whatever you want with it.

For those who fork it and want to publish a new version or something else; if you can, a credit or co-author mention is always welcome :) (though never required).

Thank you!

## Feedback

I'd love to hear your feedback

Feel free to reach out via:
* [Email](mailto:versedev.store@proton.me)
* [Twitter](https://twitter.com/versenilvis)
* [GitHub issues](https://github.com/versenilvis/iris/issues/new)
