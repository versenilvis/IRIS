# Command specification & lookup engine (`spec/`)

The `spec` module is the language parser of Iris. It understands the relationship between commands, subcommands, and flags.

## Data structures

- **`Spec`**: Top-level command definition (e.g. `git`).
- **`Subcommand`**: Recursively defined children (e.g. `commit` under `git`).
- **`Generator`**: A function providing dynamic content (e.g. file paths, docker IDs).

## The `Lookup` algorithm

The core path-traversal function:

1. **Tokenization**: Splits `"git commit -m"` into `["git", "commit", "-m"]`. Empty tokens (from trailing space) indicate the user is ready for the next suggestion level.
2. **Tree walking**: Starts at the root node (`git`) and matches each token against available subcommands.
3. **Context identification**: When traversal stops (partial word or option prefix), it defines:
   - `prefix`: the path already traversed.
   - `partial`: the word currently being typed.
4. **Result collection**: Gathers all subcommands and options matching the `partial` prefix.

## Example

Input: `git com`
1. Tokens: `["git", "com"]`.
2. Walk: root is `git`.
3. Next token `com` doesn't match `commit` exactly.
4. Stop. `partial = "com"`.
5. Suggestions: subcommands of `git` starting with `com`.
6. Return: `git commit`.

## Shell aliases & priority

- **Dynamic alias parsing**: Scans shell config files (`.bashrc`, `.zshrc`) for defined aliases.
- **Highest priority**: Aliases appear above spec and system command suggestions.
- **Token injection**: When a root command is an alias (e.g. `gr` for `go run`), Iris injects expanded tokens into the lookup engine for accurate subcommand suggestions.
- **Display**: Shows the expanded command (e.g. `tmux a -t`) with the alias name in the description (e.g. `alias: ta`).
