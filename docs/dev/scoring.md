# Scoring & ranking architecture (`internal/scoring/`)

The scoring engine ranks suggestions and completion candidates by combining frecency algorithms, workflow sequence learning, workspace tiering, and item-type priority rules.

## Core concepts

### 1. Spec mode composite scoring (`internal/scoring/scorer.go`)

In spec mode, suggestions from static command specs are evaluated against runtime signals and ranked through a composite score:

$$\text{FinalScore} = \sum w_i \cdot S_i$$

Default weights:
- **BasePriority** ($w_1 = 0.20$): Subcommands default to 30, flags default to 10 (boosted to 80 when typing `-` or `--`).
- **ContextBonus** ($w_2 = 0.20$): Active directory context and argument hints.
- **Frecency** ($w_3 = 0.20$): Normalized frecency score from execution history.
- **Transition** ($w_4 = 0.20$): Sequential pattern match based on previous command skeleton.
- **MatchQuality** ($w_5 = 0.20$): Exact match vs prefix match vs substring proximity.

### 2. Frecency decay calculation (`internal/scoring/frecency.go`)

Frecency combines execution count with step-based recency weights (`RawScore`):

$$\text{Score} = \text{Count} \times \text{Weight}(\Delta t)$$

- **$\le 1$ hour**: weight = 100
- **$\le 24$ hours**: weight = 50
- **$\le 7$ days**: weight = 20
- **$\le 30$ days**: weight = 5
- **$> 30$ days**: weight = 1
- Commands with non-zero exit codes are never recorded.

### 3. Workflow sequence learning (`command_sequences` & `command_transitions`)

Iris tracks sequential command pairs to suggest developer workflows:
- **Skeleton transitions (`command_transitions`)**: Structural transitions between base commands (`git add` $\rightarrow$ `git commit`, `go build` $\rightarrow$ `./iris`).
- **Full command sequences (`command_sequences`)**: Exact command pairs $(C_{prev}, C_{next})$ preserving arguments, working directory, and `project_id`. Bootstrapped in the background from shell history.

### 4. Workspace candidate tiering for predictions

When retrieving candidates for ghost text prediction (`QuerySequenceCandidates` and `QueryHistoryCandidates`), results are categorized into workspace tiers calculated in Go:

- **Tier 4 (Exact CWD)**: Recorded working directory matches `cwd` exactly.
- **Tier 3 (Descendant)**: Recorded working directory is beneath the current directory within the same project.
- **Tier 2 (Ancestor)**: Recorded working directory is above the current directory within the same project.
- **Tier 1 (Project Siblings)**: Different directory branches sharing the same `project_id`.
- **Tier 0 (Foreign Scope)**: Outside the current project scope or different non-git directories.

### 5. Two-phase candidate retrieval & lazy scope gating

1. **Local phase**: Queries `history_entries` and `command_sequences` where `cwd = ? OR project_id = ?`, bounded by `cmd >= prefix AND cmd < prefixUpperBound` and `instr(cmd, prefix) = 1`.
2. **Global phase**: If local candidate pool is below threshold, queries foreign scopes (`project_id != ? OR project_id IS NULL`), deduplicating against local candidates.
3. **Lazy scope gate**: Instead of running expensive `COUNT(DISTINCT)` aggregates across all candidates in the global SQL query, `store.ScopeCount` queries `COUNT(DISTINCT COALESCE(NULLIF(project_id, ''), cwd))` only on-demand for Tier 0 candidates. Tier 0 candidates require `ScopeCount >= 3` to pass the admission gate.

### 6. Storage & migration

- SQLite database (`~/.local/share/iris/history.db`) runs with WAL mode and `PRAGMA busy_timeout = 5000`.
- Safe schema migration: adds `project_id` column dynamically with an automatic backup file (`history.db.bak`, permissions `0600`) created only when legacy schema migration runs.