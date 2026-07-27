<div align="center">
  <!-- <img width="50%" alt="banner" src="https://github.com/user-attachments/assets/c5ec623b-8259-473f-b7c3-3d01a64deb5d" /> -->
  <!-- <img width="25%" alt="logo" src="https://github.com/user-attachments/assets/79d3913c-56b7-42cb-8b07-53e98f39322b" /> -->
  <img width="15%" alt="logo" src="https://github.com/user-attachments/assets/10b7ca98-872b-44a2-bdcd-265f18aa0564" />
  
  <!-- <h1>IRIS</h1> -->
  <p>IRIS (Intelligent Real-time Input Suggestion) - A shell auto-completion tool that works like code editor's IntelliSense</p>
  
  [![GitHub Actions](https://img.shields.io/github/actions/workflow/status/versenilvis/IRIS/release.yml?branch=main&style=for-the-badge&logo=github&logoColor=white&label=Actions)](https://github.com/versenilvis/IRIS/actions/workflows/release.yml)
  [![License: 0BSD](https://img.shields.io/badge/License-0BSD-blue?style=for-the-badge&logo=github&logoColor=white)](./LICENSE)
  
  <a href="#why-iris-instead-of-fig">Comparison</a> · <a href="#installation">Installation</a> · <a href="./docs/README.md">Docs</a> · <a href="./docs/README.md#shortcuts">Shortcuts</a> · <a href="./docs/README.md#reporting-bugs">Reporting bugs</a>

</div>
<div align="center">
  <!-- <img width="800" height="450" alt="output" src="https://github.com/user-attachments/assets/fd586c5b-89b6-4f6a-af9c-29248db5edc3" />-->
  <!-- <img width="1280" height="720" alt="output" src="https://github.com/user-attachments/assets/a1003bda-7722-4185-9514-f1bf83fbb504" /> -->
  <img width="1920" height="1080" alt="Ghostty terminal showcase" src="https://github.com/user-attachments/assets/dc2434c2-4f59-4432-b77f-78e65c4d412e" />
  👻 <i>Ghostty terminal</i>
</div>

<br>

**IRIS is built on top of TTY, so it runs everywhere. It just needs a terminal!**

Run iris wherever you already work; your local machine, a remote server, or anywhere you can ssh. Each suggestion menu renders directly inline inside your real terminal session, not an app's imitation of one, so it never breaks full-screen TUIs or terminal formatting. Automatically index your aliases and shell history to suggest commands that match your actual workflow in real time. Change configurations and propagate them instantly without restarting your shell. One single local native Go binary, not an app: no gui, no electron, no mac-only wrapper, no account, no telemetry. (if you've used fig: it's that, rebuilt to run purely on TTY)

<!--
  [![Status](https://img.shields.io/badge/status-beta-yellow?style=for-the-badge&logo=github&logoColor=white)]()
  [![Documentation](https://img.shields.io/badge/docs-available-brightgreen?style=for-the-badge&logo=github&logoColor=white)](./docs)
  [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen?style=for-the-badge&logo=github&logoColor=white)](./.github/CONTRIBUTING.md)
-->


## Why Iris instead of Fig

> [!IMPORTANT]
> **[Fig](https://app.fig.io/) was officially sunset in September 2024 and migrated to Amazon Q Developer (which requires cloud authentication and proprietary bloat)**  
> **IRIS is the lightweight, open-source, zero-telemetry alternative built purely on native Go and TTY with no accounts, no GUI app, and no background daemons required**

### How it compares

| Feature                     | Iris                 | Fig              |
| :-------------------------- | :------------------- | :--------------- |
| **Platforms**               | Linux, macOS         | macOS only       |
| **Engine**                  | Native Go (TTY)      | Electron         |
| **Startup**                 | Near-zero overhead   | Low overhead     |
| **UI**                      | Inline overlay       | GUI popover      |
| **Remote SSH**              | TTY-native, portable | macOS GUI-bound  |
| **Tmux**                    | ✓                    | Limited          |
| **Linux virtual terminals** | ✓                    | -                |
| **Memory**                  | < 15 MB              | Electron runtime |

## Why not shell autocomplete plugins?

Shell plugins are great, but they also come with trade-offs. And also, not everyone use Zsh or Fish especially on SSH.

| Feature                     | Iris                    | Shell plugins                |
| :-------------------------- | :---------------------- | :--------------------------- |
| **Installation**            | Single binary           | Plugin manager required      |
| **Shell support**           | Most shells supported   | Usually shell-specific       |
| **Startup**                 | No shell initialization | Increases shell startup time |
| **SSH**                     | One config              | Per-shell config             |
| **Tmux**                    | ✓                       | Depends on the shell         |
| **Linux virtual terminals** | ✓                       | Depends on the shell         |

## Installation

```bash
curl -sSL https://raw.githubusercontent.com/versenilvis/iris/main/scripts/install.sh | sh
```

> [!WARNING]
> Currently, Windows is not supported

## Theme
<div align="center">
  <img width="1920" height="1080" alt="Kitty terminal showcase" src="https://github.com/user-attachments/assets/d45cf36e-6d7d-437a-9af9-28cd994bf55f" />
  😺 <i>Kitty terminal</i>
</div>

> [!NOTE]
> Currently, IRIS doesn't have custom theme but it does have 2 basic styles

<table>
  <tr>
    <td align="center"><b>Modern style</b></td>
    <td align="center"><b>Classic style</b></td>
  </tr>
  <tr>
    <td><img width="100%" alt="spec" src="https://github.com/user-attachments/assets/cb92cc64-a08d-43ad-a412-22ab479e53aa" /></td>
    <td><img width="100%" alt="spec" src="https://github.com/user-attachments/assets/6e68b330-39d9-4750-ac95-e94514ba4e7b" /></td>
  </tr>
  <tr>
    <td><img width="100%" alt="history" src="https://github.com/user-attachments/assets/fd1a272a-19fe-472a-83d6-eec790813403" /></td>
    <td><img width="100%" alt="history" src="https://github.com/user-attachments/assets/a1983d05-b771-4277-91cf-96009875979b" /></td>
  </tr>
</table>

## Docs

- [Getting started](./docs/README.md#getting-started): dependencies, installation methods, and shell integration setup
- [Shortcuts](./docs/README.md#shortcuts): core navigation, shortcuts table, mode switching, and ghost text
- [Configuration guide](./docs/README.md#configuration-guide): TOML configuration settings including AI provider options
- [Reporting bugs](./docs/README.md#reporting-bugs): debug mode, log inspection, and crash reporting
- [Developer documentation](./docs/dev/README.md): system architecture overview, engine design, and contribution guide
 
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
