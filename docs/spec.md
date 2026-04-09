# Command Specification & Navigation (`commands/core/spec.go`)

The `spec` module is the "language parser" of IRIS. it understands the relationship between commands, subcommands, and flags.

## Data Structures

- **`Spec`**: The top-level definition (e.g., for `git`).
- **`Subcommand`**: Recursively defined children (e.g., `commit` under `git`).
- **`Generator`**: A function that provides dynamic content (like files or docker IDs).

## The `Lookup` Algorithm

This is the most critical function in IRIS. It follows these steps:

1. **Tokenization**: Splits `"git commit -m"` into `["git", "commit", "-m"]`. 
   - *Empty tokens* (from a trailing space) indicate the user is ready for the next level of suggestions.
2. **Tree Walking**: It starts at the root (`git`) and tries to match each following token against the current node's subcommands.
3. **Context Identification**: Once it can't walk any further (e.g., a partial word or an option), it defines:
   - `prefix`: The path already traveled.
   - `partial`: The word currently being typed.
4. **Result Collection**: It gathers all valid subcommands and options that match the `partial` prefix.

## Example
Input: `git com`
1. Tokens: `["git", "com"]`
2. Walk: Root is `git`.
3. Next token `com` doesn't exactly match `commit`.
4. Stop walking. `partial` = `com`.
5. Suggestions: Look for subcommands of `git` starting with `com`.
6. Return: `git commit`.
