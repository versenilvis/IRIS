# Development guide

This document outlines the development workflow, build system, and coding conventions for contributors working on the Iris codebase.

## Environment setup

Clone the repository and sync dependencies:

```bash
git clone https://github.com/versenilvis/iris.git
cd iris
go mod tidy
```

## Modular `just` command runner

We use `just` as our primary command runner. Recipes are organized into modular files under `justfiles/`:

- `justfiles/build.just` (`[build]` group): `build`, `optimized-build`, `build-release <version>`
- `justfiles/dev.just` (`[dev]` group): `run`, `config-init`, `reload`, `copy`
- `justfiles/test.just` (`[test]` group): `test`, `lint`, `analyze` (alias `ana`)
- `justfiles/gen.just` (`[gen]` group): `gen-docs`
- `justfiles/pkg.just` (`[pkg]` group): `pkg`
- `justfiles/debug.just` (`[debug]` group): `debug`, `debug-update`, `debug-changelog`, `debug-notify`, `debug-install`, `debug-autoupdate`

Run `just -l` to quickly view all available recipes organized by group.

## Building and testing

### Rapid hot-reload

```bash
just reload
```

### Running unit tests

```bash
just test
# OR
go test -v ./...
```

### Running linter

```bash
just lint
# OR
golangci-lint run ./...
```
