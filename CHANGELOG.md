# Changelog

All notable changes to this project are documented in this file.

## [Unreleased]

### Bug fixes

- Follow the shell working directory in the wrapper process ([#129](https://github.com/versenilvis/iris/issues/129)) ([4a4ab3c102067f7a24f737987e80fbb7cedbe86c](https://github.com/versenilvis/iris/commit/4a4ab3c102067f7a24f737987e80fbb7cedbe86c))
- Fall back to raw body for nightly releases ([005ed5ae469953eec4d20d91fe647cb42151b98b](https://github.com/versenilvis/iris/commit/005ed5ae469953eec4d20d91fe647cb42151b98b))
- Invalidate cache when the updater channel changes ([5845fcb0511b3a6b6bfae2cadba3735e244e3331](https://github.com/versenilvis/iris/commit/5845fcb0511b3a6b6bfae2cadba3735e244e3331))
- Only strip GoReleaser's exact heading, not real content ([a8779d56e67beb859caaee3af7371bbdb829ef77](https://github.com/versenilvis/iris/commit/a8779d56e67beb859caaee3af7371bbdb829ef77))
- Use full commit SHAs in CHANGELOG.md for GitHub autolinking ([4d5e86ac9d142c38be7c9231f6576cc73371a73d](https://github.com/versenilvis/iris/commit/4d5e86ac9d142c38be7c9231f6576cc73371a73d))

### Documentation

- Document changelog and auto-update ([80ee178932a7d2b1b75dd516e9b95aee035c05b5](https://github.com/versenilvis/iris/commit/80ee178932a7d2b1b75dd516e9b95aee035c05b5))

### Features

- Generate CHANGELOG.md with git-cliff ([6f9a2b9500535a70f08b2bd47949d9e560122013](https://github.com/versenilvis/iris/commit/6f9a2b9500535a70f08b2bd47949d9e560122013))
- Group release notes by commit type ([0421ed853b79759ba6e215a95288c9f5a72ef9c3](https://github.com/versenilvis/iris/commit/0421ed853b79759ba6e215a95288c9f5a72ef9c3))
- Add changelog command ([5745dc7ca46872076d3f9369d85fb3743db856e5](https://github.com/versenilvis/iris/commit/5745dc7ca46872076d3f9369d85fb3743db856e5))
- Show changelog in update notification ([27e7453ab27998c0ffd643cb8821c0a56c998989](https://github.com/versenilvis/iris/commit/27e7453ab27998c0ffd643cb8821c0a56c998989))
- Add auto-update option ([db2055a4c3e50a11099a5a32c881c5241014bfc9](https://github.com/versenilvis/iris/commit/db2055a4c3e50a11099a5a32c881c5241014bfc9))
- Auto install updates in background ([359480d77c065f783761ee86e988eaee58c1ec31](https://github.com/versenilvis/iris/commit/359480d77c065f783761ee86e988eaee58c1ec31))
- Shape nightly changelog like stable release groups ([ae39b4731f381d2e5f192b805e49a57f463f38e7](https://github.com/versenilvis/iris/commit/ae39b4731f381d2e5f192b805e49a57f463f38e7))
- Render release notes with glamour ([8b2dba99fa0022016d7d731ef617fb2aaa1868f8](https://github.com/versenilvis/iris/commit/8b2dba99fa0022016d7d731ef617fb2aaa1868f8))
- Compact rendering and per-version lookup ([445553b15180dff0af5113e8a89be55442f54d64](https://github.com/versenilvis/iris/commit/445553b15180dff0af5113e8a89be55442f54d64))

## [v0.5.2](https://github.com/versenilvis/iris/releases/tag/v0.5.2) - 2026-08-06

### Bug fixes

- Start iris in multiplexer panes instead of inheriting a dead one ([#120](https://github.com/versenilvis/iris/issues/120)) ([fc1634b00799b2c88ef8f02710da1e5e39787e7d](https://github.com/versenilvis/iris/commit/fc1634b00799b2c88ef8f02710da1e5e39787e7d))

## [v0.5.1](https://github.com/versenilvis/iris/releases/tag/v0.5.1) - 2026-08-06

### Bug fixes

- Mode doesn't work correctly ([#128](https://github.com/versenilvis/iris/issues/128)) ([7a97ea9016a40693f89115be089b4b8e6ad58c2d](https://github.com/versenilvis/iris/commit/7a97ea9016a40693f89115be089b4b8e6ad58c2d))

### Documentation

- Remove ⌘ ([ce404335ed8114c64a2d7068dcd1ea2437a751d2](https://github.com/versenilvis/iris/commit/ce404335ed8114c64a2d7068dcd1ea2437a751d2))

## [v0.5.0](https://github.com/versenilvis/iris/releases/tag/v0.5.0) - 2026-08-06

### Bug fixes

- Update stale vendorHash after lipgloss v2 upgrade ([#99](https://github.com/versenilvis/iris/issues/99)) ([476ca3946b5e129471663a9269d1481b22fd4963](https://github.com/versenilvis/iris/commit/476ca3946b5e129471663a9269d1481b22fd4963))
- Support kitty keyboard protocol for ctrl keybindings ([#101](https://github.com/versenilvis/iris/issues/101)) ([d669e97423a7ca9326d17b1289d06ba90942bd77](https://github.com/versenilvis/iris/commit/d669e97423a7ca9326d17b1289d06ba90942bd77))
- Isolate __complete probe from the controlling terminal ([#111](https://github.com/versenilvis/iris/issues/111)) ([0b77c9334b218eb1d7249983ceae424904ebe071](https://github.com/versenilvis/iris/commit/0b77c9334b218eb1d7249983ceae424904ebe071))
- Make Right Arrow configurable and correct overlay positioning ([#110](https://github.com/versenilvis/iris/issues/110)) ([dfaa051f31671eaf45b0bd0bc5a2694df00e3db0](https://github.com/versenilvis/iris/commit/dfaa051f31671eaf45b0bd0bc5a2694df00e3db0))
- Handle non-ascii characters command ([#118](https://github.com/versenilvis/iris/issues/118)) ([75aa1fbe383c05e33ca10c13798225af2c23cdf3](https://github.com/versenilvis/iris/commit/75aa1fbe383c05e33ca10c13798225af2c23cdf3))
- Update nightly releases ([#125](https://github.com/versenilvis/iris/issues/125)) ([0595ddd4b903639464b6248be60de3c76ed77738](https://github.com/versenilvis/iris/commit/0595ddd4b903639464b6248be60de3c76ed77738))

### Documentation

- Add navigate-right config ([b79d25a0827805a5eaed12ca56b022e5e9e97969](https://github.com/versenilvis/iris/commit/b79d25a0827805a5eaed12ca56b022e5e9e97969))
- Theme config and list of theme templates ([#116](https://github.com/versenilvis/iris/issues/116)) ([b5734b4e85dfb3cf6654300b42663adbb6ace231](https://github.com/versenilvis/iris/commit/b5734b4e85dfb3cf6654300b42663adbb6ace231))
- Add themes directory reference ([3a97efad183a59d46dbd3f7b7aa2f64a4e2eb118](https://github.com/versenilvis/iris/commit/3a97efad183a59d46dbd3f7b7aa2f64a4e2eb118))
- AI guide ([#117](https://github.com/versenilvis/iris/issues/117)) ([30619a862182371d3654989cc4d273eba381f0f5](https://github.com/versenilvis/iris/commit/30619a862182371d3654989cc4d273eba381f0f5))
- Add auto-start setup instructions ([53f104577f5b7f5a0c398b061b03e77f74c0f71a](https://github.com/versenilvis/iris/commit/53f104577f5b7f5a0c398b061b03e77f74c0f71a))
- Small note in shell config and fix line space ([994ff836752bf4c5246a3aec236dd597421dccf0](https://github.com/versenilvis/iris/commit/994ff836752bf4c5246a3aec236dd597421dccf0))

### Features

- Adopt x/ansi instead of hand rolling escape sequences ([#96](https://github.com/versenilvis/iris/issues/96)) ([df42c95fa935d647fe27fad4414630e606951187](https://github.com/versenilvis/iris/commit/df42c95fa935d647fe27fad4414630e606951187))
- Support tool aliases ([#107](https://github.com/versenilvis/iris/issues/107)) ([6faa2e3274b71c9056f3b8d8445f2d3fc7e609ff](https://github.com/versenilvis/iris/commit/6faa2e3274b71c9056f3b8d8445f2d3fc7e609ff))
- Add theme customization ([#81](https://github.com/versenilvis/iris/issues/81)) ([b14b26e663fe0c5b6144d46370c2ea6bb2d956e4](https://github.com/versenilvis/iris/commit/b14b26e663fe0c5b6144d46370c2ea6bb2d956e4))

### Refactors

- Use XDG config for all platforms ([#119](https://github.com/versenilvis/iris/issues/119)) ([d299d1c7c7c66394c6a260451410f87b6559ac9c](https://github.com/versenilvis/iris/commit/d299d1c7c7c66394c6a260451410f87b6559ac9c))

## [v0.4.21](https://github.com/versenilvis/iris/releases/tag/v0.4.21) - 2026-07-31

### Documentation

- Reorder badge and description ([23bd7188990253500bfd142bc8213c886d7be78a](https://github.com/versenilvis/iris/commit/23bd7188990253500bfd142bc8213c886d7be78a))

### Features

- Support for XDG base directories ([#95](https://github.com/versenilvis/iris/issues/95)) ([bf75dc239bd2432c10c3e36bc313b5407c23807f](https://github.com/versenilvis/iris/commit/bf75dc239bd2432c10c3e36bc313b5407c23807f))
- Upgrade to lipgloss v2 ([#94](https://github.com/versenilvis/iris/issues/94)) ([efc49bacfe7249880e3d2d2410abc43df51f5c18](https://github.com/versenilvis/iris/commit/efc49bacfe7249880e3d2d2410abc43df51f5c18))

## [v0.4.20](https://github.com/versenilvis/iris/releases/tag/v0.4.20) - 2026-07-31

### Bug fixes

- Dynamic select key ui component ([#90](https://github.com/versenilvis/iris/issues/90)) ([8a1cec331aff82cb7c9463be77240c1775eeb8df](https://github.com/versenilvis/iris/commit/8a1cec331aff82cb7c9463be77240c1775eeb8df))

## [v0.4.19](https://github.com/versenilvis/iris/releases/tag/v0.4.19) - 2026-07-31

### Bug fixes

- Resolve enter key submission and keybinding issues ([#89](https://github.com/versenilvis/iris/issues/89)) ([235a7f98ea2cac3f9c546955dc0b70097c9e442f](https://github.com/versenilvis/iris/commit/235a7f98ea2cac3f9c546955dc0b70097c9e442f))

## [v0.4.15](https://github.com/versenilvis/iris/releases/tag/v0.4.15) - 2026-07-31

### Bug fixes

- Respect zdotdir config ([#83](https://github.com/versenilvis/iris/issues/83)) ([e9e50c7dea0186efe45e86a194fa8327e9c8d441](https://github.com/versenilvis/iris/commit/e9e50c7dea0186efe45e86a194fa8327e9c8d441))

## [v0.4.14](https://github.com/versenilvis/iris/releases/tag/v0.4.14) - 2026-07-31

### Bug fixes

- Sync config_cmd.go template with init.go template ([9663788de65b06767ab21a9f5ac1c41386364c55](https://github.com/versenilvis/iris/commit/9663788de65b06767ab21a9f5ac1c41386364c55))

## [v0.4.13](https://github.com/versenilvis/iris/releases/tag/v0.4.13) - 2026-07-31

### Documentation

- Warning about potential conflicts ([4421d6e6a53cc3ed43a2688abac16f9d46535772](https://github.com/versenilvis/iris/commit/4421d6e6a53cc3ed43a2688abac16f9d46535772))

### Features

- Ui width, auto exec, selection and navigation configs ([#80](https://github.com/versenilvis/iris/issues/80)) ([a885c7462abaacb44607df11bbf6fa3b37c4ccd3](https://github.com/versenilvis/iris/commit/a885c7462abaacb44607df11bbf6fa3b37c4ccd3))

## [v0.4.12](https://github.com/versenilvis/iris/releases/tag/v0.4.12) - 2026-07-31

### Bug fixes

- Arrow up history entries ([#77](https://github.com/versenilvis/iris/issues/77)) ([cee756e2c110f98681f3e9f4283e44bf1963999b](https://github.com/versenilvis/iris/commit/cee756e2c110f98681f3e9f4283e44bf1963999b))
- Sync shell working directory for suggestions ([#78](https://github.com/versenilvis/iris/issues/78)) ([9f43fcb3527d5a72c108d53ccc381cb5fd7943f9](https://github.com/versenilvis/iris/commit/9f43fcb3527d5a72c108d53ccc381cb5fd7943f9))

## [v0.4.11](https://github.com/versenilvis/iris/releases/tag/v0.4.11) - 2026-07-30

### Features

- Add shell login option ([#75](https://github.com/versenilvis/iris/issues/75)) ([2573ad73b4dd55932248f54558ff0e31d6fb349e](https://github.com/versenilvis/iris/commit/2573ad73b4dd55932248f54558ff0e31d6fb349e))

## [v0.4.10](https://github.com/versenilvis/iris/releases/tag/v0.4.10) - 2026-07-30

### Bug fixes

- Substring fuzzy search ([#74](https://github.com/versenilvis/iris/issues/74)) ([37182f72127b6ebf134edaae22c7d87b0a491b4c](https://github.com/versenilvis/iris/commit/37182f72127b6ebf134edaae22c7d87b0a491b4c))

### Documentation

- Move user guide to top level readme ([#70](https://github.com/versenilvis/iris/issues/70)) ([004925684443b7f250a8174a0286a99b638590b7](https://github.com/versenilvis/iris/commit/004925684443b7f250a8174a0286a99b638590b7))

## [v0.4.9](https://github.com/versenilvis/iris/releases/tag/v0.4.9) - 2026-07-30

### Bug fixes

- Show symbolic links in suggestion ([#69](https://github.com/versenilvis/iris/issues/69)) ([979bd6f5d608b9b5f34efc682a940190fd9a99ce](https://github.com/versenilvis/iris/commit/979bd6f5d608b9b5f34efc682a940190fd9a99ce))

## [v0.4.8](https://github.com/versenilvis/iris/releases/tag/v0.4.8) - 2026-07-29

### Bug fixes

- Iris stops working after exec a command in bash ([#66](https://github.com/versenilvis/iris/issues/66)) ([812d73b42a27f920ed08a1bbc1c2fd83c13c3a74](https://github.com/versenilvis/iris/commit/812d73b42a27f920ed08a1bbc1c2fd83c13c3a74))

## [v0.4.7](https://github.com/versenilvis/iris/releases/tag/v0.4.7) - 2026-07-29

### Bug fixes

- Stdin handling and exit command ([#58](https://github.com/versenilvis/iris/issues/58)) ([d73921e4199330e8996056f639d68fd9c3c62781](https://github.com/versenilvis/iris/commit/d73921e4199330e8996056f639d68fd9c3c62781))

## [v0.4.6](https://github.com/versenilvis/iris/releases/tag/v0.4.6) - 2026-07-29

### Features

- Add some more config options ([#64](https://github.com/versenilvis/iris/issues/64)) ([d6cfb9c90981db34a188744d95f4352261f33332](https://github.com/versenilvis/iris/commit/d6cfb9c90981db34a188744d95f4352261f33332))

## [v0.4.5](https://github.com/versenilvis/iris/releases/tag/v0.4.5) - 2026-07-28

### Features

- Add Nix flake for NixOS support ([#52](https://github.com/versenilvis/iris/issues/52)) ([69fb77910deab604ae16d1f4752ccb0972c70a8f](https://github.com/versenilvis/iris/commit/69fb77910deab604ae16d1f4752ccb0972c70a8f))

## [v0.4.4](https://github.com/versenilvis/iris/releases/tag/v0.4.4) - 2026-07-28

### Bug fixes

- ZDOTDIR and XDG_CONFIG_HOME config detecting ([#51](https://github.com/versenilvis/iris/issues/51)) ([ebb0bead7dce702c3e111ed7c1ed4e5f86ee6eb9](https://github.com/versenilvis/iris/commit/ebb0bead7dce702c3e111ed7c1ed4e5f86ee6eb9))

## [v0.4.3](https://github.com/versenilvis/iris/releases/tag/v0.4.3) - 2026-07-28

### Bug fixes

- Set hex color for ghost text instead of hardcoded ANSI bright-black ([#49](https://github.com/versenilvis/iris/issues/49)) ([6f488c02cd79219eb33c16c72fa8f90a0dd2315e](https://github.com/versenilvis/iris/commit/6f488c02cd79219eb33c16c72fa8f90a0dd2315e))

## [v0.4.2](https://github.com/versenilvis/iris/releases/tag/v0.4.2) - 2026-07-28

### Bug fixes

- Ranking ([#46](https://github.com/versenilvis/iris/issues/46)) ([93061c89cab1c0191d429af9f822c72afadef017](https://github.com/versenilvis/iris/commit/93061c89cab1c0191d429af9f822c72afadef017))

## [v0.4.0](https://github.com/versenilvis/iris/releases/tag/v0.4.0) - 2026-07-27

### Bug fixes

- Fuzzy search not working ([#44](https://github.com/versenilvis/iris/issues/44)) ([84f2c8f49a82a197a36eae3fdaa3ee1af4ec05c5](https://github.com/versenilvis/iris/commit/84f2c8f49a82a197a36eae3fdaa3ee1af4ec05c5))

### Documentation

- Update docs ([#43](https://github.com/versenilvis/iris/issues/43)) ([078e5e817c276b65fafda158007210b6d77b0ab2](https://github.com/versenilvis/iris/commit/078e5e817c276b65fafda158007210b6d77b0ab2))

### Features

- Transition workflow learning and fix history UX ([#39](https://github.com/versenilvis/iris/issues/39)) ([6ffb9f83f21ee71d40ad77c9feb1cd494584e316](https://github.com/versenilvis/iris/commit/6ffb9f83f21ee71d40ad77c9feb1cd494584e316))
- Better suggestion with cobra's cli __complete subcmd ([#40](https://github.com/versenilvis/iris/issues/40)) ([c4e9bcdd88cdb3db818549e2974e58a46ab543f2](https://github.com/versenilvis/iris/commit/c4e9bcdd88cdb3db818549e2974e58a46ab543f2))
- Better file path suggestions ([#41](https://github.com/versenilvis/iris/issues/41)) ([27e0792560cb2217f6c88485c78fb5c54399bbcb](https://github.com/versenilvis/iris/commit/27e0792560cb2217f6c88485c78fb5c54399bbcb))

### Refactors

- Split justfile into small files ([#42](https://github.com/versenilvis/iris/issues/42)) ([1540b9655af472345a7c8ef2e5b109c23758fe37](https://github.com/versenilvis/iris/commit/1540b9655af472345a7c8ef2e5b109c23758fe37))

## [v0.3.3](https://github.com/versenilvis/iris/releases/tag/v0.3.3) - 2026-07-14

### Bug fixes

- Skip alias expansion on paste ([#38](https://github.com/versenilvis/iris/issues/38)) ([85ac76e85a276890b8924ab73d733d22431aab04](https://github.com/versenilvis/iris/commit/85ac76e85a276890b8924ab73d733d22431aab04))

### Features

- Scoring and frecency ([#37](https://github.com/versenilvis/iris/issues/37)) ([860b475b53e285cb2f80b16a4bb27b4cb4cfcc92](https://github.com/versenilvis/iris/commit/860b475b53e285cb2f80b16a4bb27b4cb4cfcc92))

## [v0.3.1](https://github.com/versenilvis/iris/releases/tag/v0.3.1) - 2026-07-13

### Bug fixes

- Enhance prompt ([#36](https://github.com/versenilvis/iris/issues/36)) ([f54e97a6c0f3e8003f00f43ae3b07d647f90a942](https://github.com/versenilvis/iris/commit/f54e97a6c0f3e8003f00f43ae3b07d647f90a942))

### Refactors

- Better project structure ([#35](https://github.com/versenilvis/iris/issues/35)) ([1e049492421aa2b10cfcc584ac585d800d1883eb](https://github.com/versenilvis/iris/commit/1e049492421aa2b10cfcc584ac585d800d1883eb))

## [v0.3.0](https://github.com/versenilvis/iris/releases/tag/v0.3.0) - 2026-07-11

### Documentation

- Add showcase for kitty terminal ([#33](https://github.com/versenilvis/iris/issues/33)) ([4728ca0c7f8202094e5956655550e64355063b03](https://github.com/versenilvis/iris/commit/4728ca0c7f8202094e5956655550e64355063b03))

### Features

- AI suggestion ([#34](https://github.com/versenilvis/iris/issues/34)) ([94c8af7b1d4408caddf1c6ed323defdc1867d037](https://github.com/versenilvis/iris/commit/94c8af7b1d4408caddf1c6ed323defdc1867d037))

### Refactors

- Separate commands and its core into two folders ([#32](https://github.com/versenilvis/iris/issues/32)) ([ebde9ba68f72194ba0358aef2bca3ce8690617d0](https://github.com/versenilvis/iris/commit/ebde9ba68f72194ba0358aef2bca3ce8690617d0))

## [v0.2.8](https://github.com/versenilvis/iris/releases/tag/v0.2.8) - 2026-07-04

### Bug fixes

- Cannot update in dev mode ([#31](https://github.com/versenilvis/iris/issues/31)) ([6b98455957303045071073738ac7a000feefa97c](https://github.com/versenilvis/iris/commit/6b98455957303045071073738ac7a000feefa97c))

## [v0.2.7](https://github.com/versenilvis/iris/releases/tag/v0.2.7) - 2026-07-04

### Documentation

- Update to latest config and commands ([#28](https://github.com/versenilvis/iris/issues/28)) ([34f49ed690aae3e232ac95f5a5c8f5dbf8614fce](https://github.com/versenilvis/iris/commit/34f49ed690aae3e232ac95f5a5c8f5dbf8614fce))

### Features

- More commands ([#30](https://github.com/versenilvis/iris/issues/30)) ([2ef3c9081904885111d1880fb41df50c6d784fd7](https://github.com/versenilvis/iris/commit/2ef3c9081904885111d1880fb41df50c6d784fd7))

## [v0.2.6](https://github.com/versenilvis/iris/releases/tag/v0.2.6) - 2026-07-03

### Bug fixes

- Remove old binary before cp to avoid Text file busy error ([#27](https://github.com/versenilvis/iris/issues/27)) ([ffb2b991a08d1a9a39b522d3db6f07577ea7a183](https://github.com/versenilvis/iris/commit/ffb2b991a08d1a9a39b522d3db6f07577ea7a183))

## [v0.2.5](https://github.com/versenilvis/iris/releases/tag/v0.2.5) - 2026-07-03

### Bug fixes

- Actually execute install script when updating ([#26](https://github.com/versenilvis/iris/issues/26)) ([37b0651651459c62d0677fbb22026e446084d2d2](https://github.com/versenilvis/iris/commit/37b0651651459c62d0677fbb22026e446084d2d2))

## [v0.2.4](https://github.com/versenilvis/iris/releases/tag/v0.2.4) - 2026-07-03

### Features

- Allow passing specific shell to setup command ([#25](https://github.com/versenilvis/iris/issues/25)) ([ee1ef2dd733607f74125c28e8385e562f90ca38c](https://github.com/versenilvis/iris/commit/ee1ef2dd733607f74125c28e8385e562f90ca38c))

## [v0.2.3](https://github.com/versenilvis/iris/releases/tag/v0.2.3) - 2026-07-03

### Bug fixes

- Allow left/right movement after up arrow history recall ([#24](https://github.com/versenilvis/iris/issues/24)) ([4ca744983c3c1e2234debf68add1fa755ec543e4](https://github.com/versenilvis/iris/commit/4ca744983c3c1e2234debf68add1fa755ec543e4))

## [v0.2.2](https://github.com/versenilvis/iris/releases/tag/v0.2.2) - 2026-07-03

### Bug fixes

- Installer ([#22](https://github.com/versenilvis/iris/issues/22)) ([df968e8d7f1f93cf9091425096df443b4eb1fcd8](https://github.com/versenilvis/iris/commit/df968e8d7f1f93cf9091425096df443b4eb1fcd8))

### Features

- Iris uninstall command ([#21](https://github.com/versenilvis/iris/issues/21)) ([2223f1736be05c16df31ae15ce7805261b6bec49](https://github.com/versenilvis/iris/commit/2223f1736be05c16df31ae15ce7805261b6bec49))

## [v0.2.0](https://github.com/versenilvis/iris/releases/tag/v0.2.0) - 2026-07-03

### Bug fixes

- Better suggestion and overlay movement on key presses ([#12](https://github.com/versenilvis/iris/issues/12)) ([aca3d1670cc7792765d3175c51da57d9b3634281](https://github.com/versenilvis/iris/commit/aca3d1670cc7792765d3175c51da57d9b3634281))
- Menu deboucing ([#14](https://github.com/versenilvis/iris/issues/14)) ([c7960aac46f6b7b4b985f4c01cf3282dcb37fb4c](https://github.com/versenilvis/iris/commit/c7960aac46f6b7b4b985f4c01cf3282dcb37fb4c))
- Data race and overlay bugs ([#20](https://github.com/versenilvis/iris/issues/20)) ([ca9d30af9c0f9361cb5515d46afa36fdd12140d4](https://github.com/versenilvis/iris/commit/ca9d30af9c0f9361cb5515d46afa36fdd12140d4))

### Documentation

- Update guidelines for new user and developer ([#16](https://github.com/versenilvis/iris/issues/16)) ([b9cec2b3ac0e2baa62add488948b987138bc0763](https://github.com/versenilvis/iris/commit/b9cec2b3ac0e2baa62add488948b987138bc0763))
- Introducing to iris ([#11](https://github.com/versenilvis/iris/issues/11)) ([24fddaccfef8302a7850434f3ddebe5a0d2d8015](https://github.com/versenilvis/iris/commit/24fddaccfef8302a7850434f3ddebe5a0d2d8015))

### Features

- Better debugger ([#13](https://github.com/versenilvis/iris/issues/13)) ([f37f0681ce2e22054fc636f37acc79d8347c6a1f](https://github.com/versenilvis/iris/commit/f37f0681ce2e22054fc636f37acc79d8347c6a1f))

### Performance

- Prevent memory leak, data race and speed up ComputeCursorCol processing ([#15](https://github.com/versenilvis/iris/issues/15)) ([8f89c15e725e8487ded62912ef264ca7a220a76c](https://github.com/versenilvis/iris/commit/8f89c15e725e8487ded62912ef264ca7a220a76c))

## [v0.1.0](https://github.com/versenilvis/iris/releases/tag/v0.1.0) - 2026-05-30

### Bug fixes

- Hide history when delete whole line ([afa3c28a8862f72cb489fe4781a0b75a8164a3d7](https://github.com/versenilvis/iris/commit/afa3c28a8862f72cb489fe4781a0b75a8164a3d7))
- Continue to suggest command after tab ([61f05a4c82b1c3a40af0d1c870a7ed13998dd98a](https://github.com/versenilvis/iris/commit/61f05a4c82b1c3a40af0d1c870a7ed13998dd98a))
- Handle quote when tokenize ([e465d900a81ad6c847cfc4b9e408f13c2117ad56](https://github.com/versenilvis/iris/commit/e465d900a81ad6c847cfc4b9e408f13c2117ad56))
- Handle value attached flags ([1109034c8729234424d47bd94e7cc01dcec15ce3](https://github.com/versenilvis/iris/commit/1109034c8729234424d47bd94e7cc01dcec15ce3))
- Add token injection to resolve shell aliases ([503d38e5c497ad8e15a06560c072d371555a849c](https://github.com/versenilvis/iris/commit/503d38e5c497ad8e15a06560c072d371555a849c))
- Remove sync.Once so we can always up-to-date with shell config ([8b25bdac0fc321f4b4d9bb53c4f36fdc9ab2bb80](https://github.com/versenilvis/iris/commit/8b25bdac0fc321f4b4d9bb53c4f36fdc9ab2bb80))
- Prioritize aliases ([2b760fb82ae6c0e64a0c79963c30499a538ba295](https://github.com/versenilvis/iris/commit/2b760fb82ae6c0e64a0c79963c30499a538ba295))
- Iris now reloads via signal handoff ([bed9c98747f5ed1395c7b255faf8cebfa6575c67](https://github.com/versenilvis/iris/commit/bed9c98747f5ed1395c7b255faf8cebfa6575c67))
- Return last shell insteadd of default shell on reload ([ccbf59201aa3953aa5c771ee43918a46bbe10328](https://github.com/versenilvis/iris/commit/ccbf59201aa3953aa5c771ee43918a46bbe10328))
- Print reloading ([8dff68601e5e2fa630d407074ef8c092547b0c15](https://github.com/versenilvis/iris/commit/8dff68601e5e2fa630d407074ef8c092547b0c15))
- Eliminate child shell leakage and terminal corruption on reload ([11225df6c839b5c287bf31dd5f7cc38085ea2251](https://github.com/versenilvis/iris/commit/11225df6c839b5c287bf31dd5f7cc38085ea2251))
- Ctrl R open reverse search instead of switching to history ([f306566602d2eae39428864b8e8b14934d32e708](https://github.com/versenilvis/iris/commit/f306566602d2eae39428864b8e8b14934d32e708))
- Iris keep suggesting when press tab continually ([fc21a6b0cc4ea13887a7809b36b72a190e715481](https://github.com/versenilvis/iris/commit/fc21a6b0cc4ea13887a7809b36b72a190e715481))
- Up/down arrow hide menu ([c798673069e1b37f8ffdb12d0fe9742c57bd2285](https://github.com/versenilvis/iris/commit/c798673069e1b37f8ffdb12d0fe9742c57bd2285))
- Extra left border ([38118a9a4789329be922ac57d457e2c0b12a1b94](https://github.com/versenilvis/iris/commit/38118a9a4789329be922ac57d457e2c0b12a1b94))
- Should show in priority ([8032448cb6f30fd6d4a412a308abbbc08e32a1b4](https://github.com/versenilvis/iris/commit/8032448cb6f30fd6d4a412a308abbbc08e32a1b4))
- Prioritize exact matching then order ([a5c27bdb942e3c3dd2fc2e934c20e379c5ee155e](https://github.com/versenilvis/iris/commit/a5c27bdb942e3c3dd2fc2e934c20e379c5ee155e))
- Now with same fuzzy tier, which command has higher score will be in higher order ([b045aaf6af71ed9cc7332849f8147972762bd6b7](https://github.com/versenilvis/iris/commit/b045aaf6af71ed9cc7332849f8147972762bd6b7))
- Key reading loop conflict ([27c7b9732658745e7afdc201d930f5799e7f641a](https://github.com/versenilvis/iris/commit/27c7b9732658745e7afdc201d930f5799e7f641a))
- Wrong suggestion on chmod command ([d2909d93aaeee6667d2d0c042126b05fd7eb7dee](https://github.com/versenilvis/iris/commit/d2909d93aaeee6667d2d0c042126b05fd7eb7dee))
- Text file busy bug on setup ([44ff34f3789d0c4f5f61671365a701960872d0c8](https://github.com/versenilvis/iris/commit/44ff34f3789d0c4f5f61671365a701960872d0c8))
- Cd and zoxide only shows dir ([9720b5f4bcf2dd888973b91c3984187c5015ec20](https://github.com/versenilvis/iris/commit/9720b5f4bcf2dd888973b91c3984187c5015ec20))
- Support folder contains space ([3c5a0580d1e6709dbc054fd4373a2f7a654e7ea0](https://github.com/versenilvis/iris/commit/3c5a0580d1e6709dbc054fd4373a2f7a654e7ea0))
- Fix all problems related to test ([e587f027982701d390fc1bc75702c93e0d66ee7c](https://github.com/versenilvis/iris/commit/e587f027982701d390fc1bc75702c93e0d66ee7c))
- All issues by golangci-lint ([c76541258acc79abce4a3685590e3e35dbdb0ed7](https://github.com/versenilvis/iris/commit/c76541258acc79abce4a3685590e3e35dbdb0ed7))
- Fuzzy search shows commands that no one need ([88d6ee4e57d324076ae6b41aacc49241ad980ca5](https://github.com/versenilvis/iris/commit/88d6ee4e57d324076ae6b41aacc49241ad980ca5))
- Typo ([07b7d19347bfe3fe7188e331bfc6b10af937ad58](https://github.com/versenilvis/iris/commit/07b7d19347bfe3fe7188e331bfc6b10af937ad58))
- Err check rules and fix config ([56e04be4098b22662a7ad96b3a10cbe4cf8752e0](https://github.com/versenilvis/iris/commit/56e04be4098b22662a7ad96b3a10cbe4cf8752e0))
- Add user email and name to support ci/cd ([ae067fa2004790f677f8350b744b3ff34f365ff1](https://github.com/versenilvis/iris/commit/ae067fa2004790f677f8350b744b3ff34f365ff1))
- Git doesnt show branch suggesstion properly ([0ec946f408f4416aca1c520f8c326326a565a10d](https://github.com/versenilvis/iris/commit/0ec946f408f4416aca1c520f8c326326a565a10d))
- Git ([9b206fe701ab024a14903c9ff1315bb055991676](https://github.com/versenilvis/iris/commit/9b206fe701ab024a14903c9ff1315bb055991676))
- Issue #4 and #5 ([#6](https://github.com/versenilvis/iris/issues/6)) ([d06526cf21d59336d2d2f7fee5d492b1ead2f73c](https://github.com/versenilvis/iris/commit/d06526cf21d59336d2d2f7fee5d492b1ead2f73c))
- Correct nightly tag, fix tmux startup, full-word suggestions only, and keep current dir on reload ([#8](https://github.com/versenilvis/iris/issues/8)) ([8556fa626aecd555e1a87f0bd39621fe517cebd5](https://github.com/versenilvis/iris/commit/8556fa626aecd555e1a87f0bd39621fe517cebd5))

### Documentation

- LICENSE ([2eac0cd38e21cde76cd45a5d9d450805e759e140](https://github.com/versenilvis/iris/commit/2eac0cd38e21cde76cd45a5d9d450805e759e140))
- Add hot reload and shell aliases ([97d009f7ef6c0b9e020f69be155440fa7fa0a5ef](https://github.com/versenilvis/iris/commit/97d009f7ef6c0b9e020f69be155440fa7fa0a5ef))
- Updater ([72dd6940e422932caa4977f61778f38917511ed6](https://github.com/versenilvis/iris/commit/72dd6940e422932caa4977f61778f38917511ed6))

### Features

- Entry point ([c384c57cb9ee749cf594505c5ec7c499cc062b47](https://github.com/versenilvis/iris/commit/c384c57cb9ee749cf594505c5ec7c499cc062b47))
- Build and run ([f69e390e3b55c87dea82a716045081a4e62faf63](https://github.com/versenilvis/iris/commit/f69e390e3b55c87dea82a716045081a4e62faf63))
- Ipc server ([55bd6e1b481df0a4d773b0165543e54020f72a96](https://github.com/versenilvis/iris/commit/55bd6e1b481df0a4d773b0165543e54020f72a96))
- Threat-safe terminal writer ([551a096d2d0b23ee3e019695778702b3d8ee04d8](https://github.com/versenilvis/iris/commit/551a096d2d0b23ee3e019695778702b3d8ee04d8))
- Iris command ([1fb279e588f790e3402a8b8734eda72d6dfa71f2](https://github.com/versenilvis/iris/commit/1fb279e588f790e3402a8b8734eda72d6dfa71f2))
- History search ([515d6f2a0a772708227936490aa77878fcb2f601](https://github.com/versenilvis/iris/commit/515d6f2a0a772708227936490aa77878fcb2f601))
- Drop down menu tui ([25ff2bcb2c4cf60da3a884cc48c30ed2c82ffd0f](https://github.com/versenilvis/iris/commit/25ff2bcb2c4cf60da3a884cc48c30ed2c82ffd0f))
- Spec ([d384694ab341fcc384fd323bcf89ee5788088627](https://github.com/versenilvis/iris/commit/d384694ab341fcc384fd323bcf89ee5788088627))
- Some common commands ([aab3c6be2ed85e897149bf1a273761f5e4248b5d](https://github.com/versenilvis/iris/commit/aab3c6be2ed85e897149bf1a273761f5e4248b5d))
- Add partial ([e559c960d9b5589b2609045f23eac7be4cc4c531](https://github.com/versenilvis/iris/commit/e559c960d9b5589b2609045f23eac7be4cc4c531))
- File generator for file suggestion ([76458ce52ef27f090d204d1885a274fc7119cc50](https://github.com/versenilvis/iris/commit/76458ce52ef27f090d204d1885a274fc7119cc50))
- Refactor and add more commands ([f4816d2b5f0fc371a413f61ebb13f321e610e8d7](https://github.com/versenilvis/iris/commit/f4816d2b5f0fc371a413f61ebb13f321e610e8d7))
- Tracking CWD and simply return file desc ([7b37b8ae2823a4bad1660cba8486fce7a70e9534](https://github.com/versenilvis/iris/commit/7b37b8ae2823a4bad1660cba8486fce7a70e9534))
- Change priority ([b01844639bfc02dc7a727df3aae93a08648a2b6b](https://github.com/versenilvis/iris/commit/b01844639bfc02dc7a727df3aae93a08648a2b6b))
- Optimize build command ([4cc7853256480dbe1a2c0cbe1c9d2abf2b1acbb5](https://github.com/versenilvis/iris/commit/4cc7853256480dbe1a2c0cbe1c9d2abf2b1acbb5))
- Dynamic git command completion ([f64c8022b264e78ac3248e52b8b3d6386868c100](https://github.com/versenilvis/iris/commit/f64c8022b264e78ac3248e52b8b3d6386868c100))
- System PATH scannee ([525eac1bfa6e1c20246c6cc6d200be70e75790f1](https://github.com/versenilvis/iris/commit/525eac1bfa6e1c20246c6cc6d200be70e75790f1))
- Support aliases ([1f483c5ab291bea8cdc35eec2a13ec9da21b6096](https://github.com/versenilvis/iris/commit/1f483c5ab291bea8cdc35eec2a13ec9da21b6096))
- Scan shell aliases like .zshrc, .bashrc, ... ([4c70c242990acc4d1311712dfd7d189994a1b8ab](https://github.com/versenilvis/iris/commit/4c70c242990acc4d1311712dfd7d189994a1b8ab))
- Rebuild and update pkg ([02fdd99e400b09b8c39dd1a86f7b2140cf056a48](https://github.com/versenilvis/iris/commit/02fdd99e400b09b8c39dd1a86f7b2140cf056a48))
- In-place reloading on current TTY ([78e19199b0478efa1de09d549ef044c8f5d0b491](https://github.com/versenilvis/iris/commit/78e19199b0478efa1de09d549ef044c8f5d0b491))
- Shift tab to toggle menu ([7a9c28e6f68f3914d69649de6fd87e519ec5fa67](https://github.com/versenilvis/iris/commit/7a9c28e6f68f3914d69649de6fd87e519ec5fa67))
- Stop showing suggestion when delete entire line ([887b5af9580b29fa10bdf86fb478e9622d0aec79](https://github.com/versenilvis/iris/commit/887b5af9580b29fa10bdf86fb478e9622d0aec79))
- Now pressing space  show alias as full command ([7c11a3a620e2e6b0fd704b2a094399e4f98a0fe2](https://github.com/versenilvis/iris/commit/7c11a3a620e2e6b0fd704b2a094399e4f98a0fe2))
- Continue last mode ([c7106f9fc2a849f8262d36bcc73e92dd3791868a](https://github.com/versenilvis/iris/commit/c7106f9fc2a849f8262d36bcc73e92dd3791868a))
- Add bash and fish ([00826a5c075020a000948546cb9800489e7b182e](https://github.com/versenilvis/iris/commit/00826a5c075020a000948546cb9800489e7b182e))
- Ghost text ([edab80202690c088faaa7af0ba9e535a73b66260](https://github.com/versenilvis/iris/commit/edab80202690c088faaa7af0ba9e535a73b66260))
- Disable ghost text when using arrow keys ([ff8874fefcbdcc5ec036e7e5381efbe8bff96753](https://github.com/versenilvis/iris/commit/ff8874fefcbdcc5ec036e7e5381efbe8bff96753))
- Run with terminal ([5c8ac9e4a65434f11004d430fcc694001dbb7d8c](https://github.com/versenilvis/iris/commit/5c8ac9e4a65434f11004d430fcc694001dbb7d8c))
- Add more options and limit file extension ([667cf9feae0b22690f5c911b5d36bd6f6f8076f6](https://github.com/versenilvis/iris/commit/667cf9feae0b22690f5c911b5d36bd6f6f8076f6))
- Scan 1 level deeper in folder ([66b2b22d94a4f60618bbd31b905bf1242a900341](https://github.com/versenilvis/iris/commit/66b2b22d94a4f60618bbd31b905bf1242a900341))
- : install script ([ef26602ded90ba3fdd64f6f74addcd9dc7773814](https://github.com/versenilvis/iris/commit/ef26602ded90ba3fdd64f6f74addcd9dc7773814))
- Setup ([b6e5096092a8de4ec92ec35070999cf4e3998432](https://github.com/versenilvis/iris/commit/b6e5096092a8de4ec92ec35070999cf4e3998432))
- Detect another command is running to stop iris ([cea2a7bee72ec9dc5c689eb76d1dbdf89a59dca7](https://github.com/versenilvis/iris/commit/cea2a7bee72ec9dc5c689eb76d1dbdf89a59dca7))
- Show preview command when using arrow keys ([23e5b13f6d11b436dcac52126f023b4362eb1f64](https://github.com/versenilvis/iris/commit/23e5b13f6d11b436dcac52126f023b4362eb1f64))
- Up/down arrow key now can choose last commands ([3c8eca3a62315b9594fde752125dfc88a1a5f504](https://github.com/versenilvis/iris/commit/3c8eca3a62315b9594fde752125dfc88a1a5f504))
- Zoxide ([5d81871aef68ad6b1e88b76832066579bb094e8f](https://github.com/versenilvis/iris/commit/5d81871aef68ad6b1e88b76832066579bb094e8f))
- Copy to local bin ([be3c08171bdbf7b312551efa5413a9c43bc33c69](https://github.com/versenilvis/iris/commit/be3c08171bdbf7b312551efa5413a9c43bc33c69))
- All tests ([17f33ff480ff8f443107d682322d72f4b51cceb0](https://github.com/versenilvis/iris/commit/17f33ff480ff8f443107d682322d72f4b51cceb0))
- More commands, tests, add Registry conflict between tests ([d74c959be4ebc8ac535686ab18eadb3f0d648af5](https://github.com/versenilvis/iris/commit/d74c959be4ebc8ac535686ab18eadb3f0d648af5))
- Version ([e6f4bdd36d6024b2286985def87df30888244222](https://github.com/versenilvis/iris/commit/e6f4bdd36d6024b2286985def87df30888244222))
- Updater ([1e923a13d9816ce0414dbbd05d0fc219e40027e2](https://github.com/versenilvis/iris/commit/1e923a13d9816ce0414dbbd05d0fc219e40027e2))
- Update test ([f9eb1900464cb0190a0120acacde2846b1726e4a](https://github.com/versenilvis/iris/commit/f9eb1900464cb0190a0120acacde2846b1726e4a))
- Uninstaller ([eaa7344a24a0a2b5feaad1daba379a3611152557](https://github.com/versenilvis/iris/commit/eaa7344a24a0a2b5feaad1daba379a3611152557))
- Updates and versioning ([#1](https://github.com/versenilvis/iris/issues/1)) ([9fd54c8c3dfa6e87750b1365615ccabef3c53a44](https://github.com/versenilvis/iris/commit/9fd54c8c3dfa6e87750b1365615ccabef3c53a44))
- Issue template and release ci/cd ([b5c0e5ec3d99ec4d754c6722cbe15d49f7a11fed](https://github.com/versenilvis/iris/commit/b5c0e5ec3d99ec4d754c6722cbe15d49f7a11fed))
- Issue template and ci/cd ([#2](https://github.com/versenilvis/iris/issues/2)) ([8f4d8b7db7c32ac919f3f1121af251f3bbcff887](https://github.com/versenilvis/iris/commit/8f4d8b7db7c32ac919f3f1121af251f3bbcff887))
- Nightly built ([9c1c7749b08983db182cee0f0d4ffaf485a6ed6d](https://github.com/versenilvis/iris/commit/9c1c7749b08983db182cee0f0d4ffaf485a6ed6d))
- Crash log ([#7](https://github.com/versenilvis/iris/issues/7)) ([bf67e50ea612b1a1e146e4dfc67de843877374e2](https://github.com/versenilvis/iris/commit/bf67e50ea612b1a1e146e4dfc67de843877374e2))
- Toml config support ([#9](https://github.com/versenilvis/iris/issues/9)) ([3d76a3ee8816ac8fa2c5026551ff93d77bbd4762](https://github.com/versenilvis/iris/commit/3d76a3ee8816ac8fa2c5026551ff93d77bbd4762))

### Performance

- Using string.Builder instead of concatenation ([81b80d079ce20956ee99e0343fc90c0d4d8b9295](https://github.com/versenilvis/iris/commit/81b80d079ce20956ee99e0343fc90c0d4d8b9295))
- Finding the first space and using a case-insensitive comparison ([caf4aa344e3f9ae1c0691cde3c91e89fa23d4629](https://github.com/versenilvis/iris/commit/caf4aa344e3f9ae1c0691cde3c91e89fa23d4629))

### Refactors

- Improve reading aliases in shell config ([71a75ced2857f9ded98185ca2f110555f3c2146b](https://github.com/versenilvis/iris/commit/71a75ced2857f9ded98185ca2f110555f3c2146b))
- Shell adaptor ([736723fbe8a6d520ec70ecd939e3a91bf013f624](https://github.com/versenilvis/iris/commit/736723fbe8a6d520ec70ecd939e3a91bf013f624))
- Seperate into small files ([db8b6ec91bf741519ce3a056ab2d93d8dafcd1d6](https://github.com/versenilvis/iris/commit/db8b6ec91bf741519ce3a056ab2d93d8dafcd1d6))
- Seperate into small files ([3adc045c09575996ce966e8b83b4932a761bd64e](https://github.com/versenilvis/iris/commit/3adc045c09575996ce966e8b83b4932a761bd64e))


