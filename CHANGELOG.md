# Changelog

All notable changes to this project are documented in this file.

## [Unreleased]

### Bug fixes

- Follow the shell working directory in the wrapper process (#129) (4a4ab3c)

## [v0.5.2](https://github.com/versenilvis/iris/releases/tag/v0.5.2) - 2026-08-06

### Bug fixes

- Start iris in multiplexer panes instead of inheriting a dead one (#120) (fc1634b)

## [v0.5.1](https://github.com/versenilvis/iris/releases/tag/v0.5.1) - 2026-08-06

### Bug fixes

- Mode doesn't work correctly (#128) (7a97ea9)

### Documentation

- Remove ⌘ (ce40433)

## [v0.5.0](https://github.com/versenilvis/iris/releases/tag/v0.5.0) - 2026-08-06

### Bug fixes

- Update stale vendorHash after lipgloss v2 upgrade (#99) (476ca39)
- Support kitty keyboard protocol for ctrl keybindings (#101) (d669e97)
- Isolate __complete probe from the controlling terminal (#111) (0b77c93)
- Make Right Arrow configurable and correct overlay positioning (#110) (dfaa051)
- Handle non-ascii characters command (#118) (75aa1fb)
- Update nightly releases (#125) (0595ddd)

### Documentation

- Add navigate-right config (b79d25a)
- Theme config and list of theme templates (#116) (b5734b4)
- Add themes directory reference (3a97efa)
- AI guide (#117) (30619a8)
- Add auto-start setup instructions (53f1045)
- Small note in shell config and fix line space (994ff83)

### Features

- Adopt x/ansi instead of hand rolling escape sequences (#96) (df42c95)
- Support tool aliases (#107) (6faa2e3)
- Add theme customization (#81) (b14b26e)

### Refactors

- Use XDG config for all platforms (#119) (d299d1c)

## [v0.4.21](https://github.com/versenilvis/iris/releases/tag/v0.4.21) - 2026-07-31

### Documentation

- Reorder badge and description (23bd718)

### Features

- Support for XDG base directories (#95) (bf75dc2)
- Upgrade to lipgloss v2 (#94) (efc49ba)

## [v0.4.20](https://github.com/versenilvis/iris/releases/tag/v0.4.20) - 2026-07-31

### Bug fixes

- Dynamic select key ui component (#90) (8a1cec3)

## [v0.4.19](https://github.com/versenilvis/iris/releases/tag/v0.4.19) - 2026-07-31

### Bug fixes

- Resolve enter key submission and keybinding issues (#89) (235a7f9)

## [v0.4.15](https://github.com/versenilvis/iris/releases/tag/v0.4.15) - 2026-07-31

### Bug fixes

- Respect zdotdir config (#83) (e9e50c7)

## [v0.4.14](https://github.com/versenilvis/iris/releases/tag/v0.4.14) - 2026-07-31

### Bug fixes

- Sync config_cmd.go template with init.go template (9663788)

## [v0.4.13](https://github.com/versenilvis/iris/releases/tag/v0.4.13) - 2026-07-31

### Documentation

- Warning about potential conflicts (4421d6e)

### Features

- Ui width, auto exec, selection and navigation configs (#80) (a885c74)

## [v0.4.12](https://github.com/versenilvis/iris/releases/tag/v0.4.12) - 2026-07-31

### Bug fixes

- Arrow up history entries (#77) (cee756e)
- Sync shell working directory for suggestions (#78) (9f43fcb)

## [v0.4.11](https://github.com/versenilvis/iris/releases/tag/v0.4.11) - 2026-07-30

### Features

- Add shell login option (#75) (2573ad7)

## [v0.4.10](https://github.com/versenilvis/iris/releases/tag/v0.4.10) - 2026-07-30

### Bug fixes

- Substring fuzzy search (#74) (37182f7)

### Documentation

- Move user guide to top level readme (#70) (0049256)

## [v0.4.9](https://github.com/versenilvis/iris/releases/tag/v0.4.9) - 2026-07-30

### Bug fixes

- Show symbolic links in suggestion (#69) (979bd6f)

## [v0.4.8](https://github.com/versenilvis/iris/releases/tag/v0.4.8) - 2026-07-29

### Bug fixes

- Iris stops working after exec a command in bash (#66) (812d73b)

## [v0.4.7](https://github.com/versenilvis/iris/releases/tag/v0.4.7) - 2026-07-29

### Bug fixes

- Stdin handling and exit command (#58) (d73921e)

## [v0.4.6](https://github.com/versenilvis/iris/releases/tag/v0.4.6) - 2026-07-29

### Features

- Add some more config options (#64) (d6cfb9c)

## [v0.4.5](https://github.com/versenilvis/iris/releases/tag/v0.4.5) - 2026-07-28

### Features

- Add Nix flake for NixOS support (#52) (69fb779)

## [v0.4.4](https://github.com/versenilvis/iris/releases/tag/v0.4.4) - 2026-07-28

### Bug fixes

- ZDOTDIR and XDG_CONFIG_HOME config detecting (#51) (ebb0bea)

## [v0.4.3](https://github.com/versenilvis/iris/releases/tag/v0.4.3) - 2026-07-28

### Bug fixes

- Set hex color for ghost text instead of hardcoded ANSI bright-black (#49) (6f488c0)

## [v0.4.2](https://github.com/versenilvis/iris/releases/tag/v0.4.2) - 2026-07-28

### Bug fixes

- Ranking (#46) (93061c8)

## [v0.4.0](https://github.com/versenilvis/iris/releases/tag/v0.4.0) - 2026-07-27

### Bug fixes

- Fuzzy search not working (#44) (84f2c8f)

### Documentation

- Update docs (#43) (078e5e8)

### Features

- Transition workflow learning and fix history UX (#39) (6ffb9f8)
- Better suggestion with cobra's cli __complete subcmd (#40) (c4e9bcd)
- Better file path suggestions (#41) (27e0792)

### Refactors

- Split justfile into small files (#42) (1540b96)

## [v0.3.3](https://github.com/versenilvis/iris/releases/tag/v0.3.3) - 2026-07-14

### Bug fixes

- Skip alias expansion on paste (#38) (85ac76e)

### Features

- Scoring and frecency (#37) (860b475)

## [v0.3.1](https://github.com/versenilvis/iris/releases/tag/v0.3.1) - 2026-07-13

### Bug fixes

- Enhance prompt (#36) (f54e97a)

### Refactors

- Better project structure (#35) (1e04949)

## [v0.3.0](https://github.com/versenilvis/iris/releases/tag/v0.3.0) - 2026-07-11

### Documentation

- Add showcase for kitty terminal (#33) (4728ca0)

### Features

- AI suggestion (#34) (94c8af7)

### Refactors

- Separate commands and its core into two folders (#32) (ebde9ba)

## [v0.2.8](https://github.com/versenilvis/iris/releases/tag/v0.2.8) - 2026-07-04

### Bug fixes

- Cannot update in dev mode (#31) (6b98455)

## [v0.2.7](https://github.com/versenilvis/iris/releases/tag/v0.2.7) - 2026-07-04

### Documentation

- Update to latest config and commands (#28) (34f49ed)

### Features

- More commands (#30) (2ef3c90)

## [v0.2.6](https://github.com/versenilvis/iris/releases/tag/v0.2.6) - 2026-07-03

### Bug fixes

- Remove old binary before cp to avoid Text file busy error (#27) (ffb2b99)

## [v0.2.5](https://github.com/versenilvis/iris/releases/tag/v0.2.5) - 2026-07-03

### Bug fixes

- Actually execute install script when updating (#26) (37b0651)

## [v0.2.4](https://github.com/versenilvis/iris/releases/tag/v0.2.4) - 2026-07-03

### Features

- Allow passing specific shell to setup command (#25) (ee1ef2d)

## [v0.2.3](https://github.com/versenilvis/iris/releases/tag/v0.2.3) - 2026-07-03

### Bug fixes

- Allow left/right movement after up arrow history recall (#24) (4ca7449)

## [v0.2.2](https://github.com/versenilvis/iris/releases/tag/v0.2.2) - 2026-07-03

### Bug fixes

- Installer (#22) (df968e8)

### Features

- Iris uninstall command (#21) (2223f17)

## [v0.2.0](https://github.com/versenilvis/iris/releases/tag/v0.2.0) - 2026-07-03

### Bug fixes

- Better suggestion and overlay movement on key presses (#12) (aca3d16)
- Menu deboucing (#14) (c7960aa)
- Data race and overlay bugs (#20) (ca9d30a)

### Documentation

- Update guidelines for new user and developer (#16) (b9cec2b)
- Introducing to iris (#11) (24fddac)

### Features

- Better debugger (#13) (f37f068)

### Performance

- Prevent memory leak, data race and speed up ComputeCursorCol processing (#15) (8f89c15)

## [v0.1.0](https://github.com/versenilvis/iris/releases/tag/v0.1.0) - 2026-05-30

### Bug fixes

- Hide history when delete whole line (afa3c28)
- Continue to suggest command after tab (61f05a4)
- Handle quote when tokenize (e465d90)
- Handle value attached flags (1109034)
- Add token injection to resolve shell aliases (503d38e)
- Remove sync.Once so we can always up-to-date with shell config (8b25bda)
- Prioritize aliases (2b760fb)
- Iris now reloads via signal handoff (bed9c98)
- Return last shell insteadd of default shell on reload (ccbf592)
- Print reloading (8dff686)
- Eliminate child shell leakage and terminal corruption on reload (11225df)
- Ctrl R open reverse search instead of switching to history (f306566)
- Iris keep suggesting when press tab continually (fc21a6b)
- Up/down arrow hide menu (c798673)
- Extra left border (38118a9)
- Should show in priority (8032448)
- Prioritize exact matching then order (a5c27bd)
- Now with same fuzzy tier, which command has higher score will be in higher order (b045aaf)
- Key reading loop conflict (27c7b97)
- Wrong suggestion on chmod command (d2909d9)
- Text file busy bug on setup (44ff34f)
- Cd and zoxide only shows dir (9720b5f)
- Support folder contains space (3c5a058)
- Fix all problems related to test (e587f02)
- All issues by golangci-lint (c765412)
- Fuzzy search shows commands that no one need (88d6ee4)
- Typo (07b7d19)
- Err check rules and fix config (56e04be)
- Add user email and name to support ci/cd (ae067fa)
- Git doesnt show branch suggesstion properly (0ec946f)
- Git (9b206fe)
- Issue #4 and #5 (#6) (d06526c)
- Correct nightly tag, fix tmux startup, full-word suggestions only, and keep current dir on reload (#8) (8556fa6)

### Documentation

- LICENSE (2eac0cd)
- Add hot reload and shell aliases (97d009f)
- Updater (72dd694)

### Features

- Entry point (c384c57)
- Build and run (f69e390)
- Ipc server (55bd6e1)
- Threat-safe terminal writer (551a096)
- Iris command (1fb279e)
- History search (515d6f2)
- Drop down menu tui (25ff2bc)
- Spec (d384694)
- Some common commands (aab3c6b)
- Add partial (e559c96)
- File generator for file suggestion (76458ce)
- Refactor and add more commands (f4816d2)
- Tracking CWD and simply return file desc (7b37b8a)
- Change priority (b018446)
- Optimize build command (4cc7853)
- Dynamic git command completion (f64c802)
- System PATH scannee (525eac1)
- Support aliases (1f483c5)
- Scan shell aliases like .zshrc, .bashrc, ... (4c70c24)
- Rebuild and update pkg (02fdd99)
- In-place reloading on current TTY (78e1919)
- Shift tab to toggle menu (7a9c28e)
- Stop showing suggestion when delete entire line (887b5af)
- Now pressing space  show alias as full command (7c11a3a)
- Continue last mode (c7106f9)
- Add bash and fish (00826a5)
- Ghost text (edab802)
- Disable ghost text when using arrow keys (ff8874f)
- Run with terminal (5c8ac9e)
- Add more options and limit file extension (667cf9f)
- Scan 1 level deeper in folder (66b2b22)
- : install script (ef26602)
- Setup (b6e5096)
- Detect another command is running to stop iris (cea2a7b)
- Show preview command when using arrow keys (23e5b13)
- Up/down arrow key now can choose last commands (3c8eca3)
- Zoxide (5d81871)
- Copy to local bin (be3c081)
- All tests (17f33ff)
- More commands, tests, add Registry conflict between tests (d74c959)
- Version (e6f4bdd)
- Updater (1e923a1)
- Update test (f9eb190)
- Uninstaller (eaa7344)
- Updates and versioning (#1) (9fd54c8)
- Issue template and release ci/cd (b5c0e5e)
- Issue template and ci/cd (#2) (8f4d8b7)
- Nightly built (9c1c774)
- Crash log (#7) (bf67e50)
- Toml config support (#9) (3d76a3e)

### Performance

- Using string.Builder instead of concatenation (81b80d0)
- Finding the first space and using a case-insensitive comparison (caf4aa3)

### Refactors

- Improve reading aliases in shell config (71a75ce)
- Shell adaptor (736723f)
- Seperate into small files (db8b6ec)
- Seperate into small files (3adc045)
