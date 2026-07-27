# Dynamic File Suggestion Generator (`commands/core/filegen.go`)

The `FileGenerator` is a dynamic suggestion engine that bridges the gap between static command specs and your live filesystem.

## How it works

The generator is assigned to commands or options that expect file paths (e.g., `git add`, `cat`, `go run`).

### Path Resolution
It handles both local and nested paths:
- `cat m` -> Scans `./` for entries starting with `m`.
- `cat src/` -> Scans `src/` for all entries.
- `cat src/main` -> Scans `src/` for entries starting with `main`.

### Intelligence Features

1. **Extension Filtering**: You can restrict the generator to specific files.
   - `FileGenerator(".go")` -> Only shows `.go` files (used in `go build`).
2. **Directory Suffixing**: When a directory is suggested, Iris appends a `/` to the command string. This allows you to immediately continue typing to dive into the next level of the folder tree.
3. **Descriptions**: It automatically converts common extensions into human-readable labels (e.g., `.mp4` -> `video`, `.zip` -> `archive`).
4. **Directory-Only Mode**: Commands like `cd` use a special filter to hide all files and only show folders.

## Code Example
```go
// Registration for 'go run'
Generator: core.FileGenerator(".go")
```
When typing `go run `, Iris will scan the current directory, filter for `.go` extensions, and ignore images, binaries, or other unrelated files.
