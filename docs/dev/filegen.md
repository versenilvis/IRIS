# Dynamic file suggestion generator (`spec/filegen.go`)

The `FileGenerator` is a dynamic suggestion engine bridging static command specs with your live filesystem.

## How it works

File lookup is the default for command arguments, so commands do not register an unfiltered file generator. Existing directory/extension filters and custom generators (such as git, chmod, and zoxide) keep their special handling.

A matching path also triggers lookup for unknown commands, invalid subcommands, assignments such as `export ROOT=./src/`, and redirections. Command validation and argument limits cannot block this fallback. The lookup preserves quotes and shell syntax when extending a path.

File results take priority over history and AI suggestions when an existing path is being entered.

### Path resolution

Handles both local and nested paths:
- `cat m` -> Scans `./` for entries starting with `m`.
- `cat src/` -> Scans `src/` for all entries.
- `cat src/main` -> Scans `src/` for entries starting with `main`.

### Intelligence features

1. **Extension filtering**: Restricts the generator to specific file extensions (e.g. `FileGenerator(".go")` for `go build`).
2. **Directory suffixing**: When a directory is suggested, Iris appends `/` without an extra trailing space, allowing continuous folder traversal.
3. **Descriptions**: Converts file extensions into human-readable labels (e.g. `.mp4` -> `video`, `.zip` -> `archive`).
4. **Directory-only mode**: Commands like `cd` filter out files to display directories only.

## Code example

```go
// registration for 'go run'
Generator: spec.FileGenerator(".go")
```

When typing `go run `, Iris scans the target directory, filters for `.go` extensions, and ignores unrelated binaries or assets.
