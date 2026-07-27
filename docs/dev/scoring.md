# Scoring & ranking architecture (`internal/scoring/`)

The scoring engine ranks suggestions by combining frecency algorithms, workflow sequence learning, and item-type priority rules.

## Core concepts

### 1. Frecency calculation (`internal/scoring/frecency.go`)

Frecency combines execution **frequency** with **recency** decay:

$$\text{Score} = \text{Frequency} \times e^{-\lambda \Delta t}$$

Commands executed recently receive a higher score multiplier that decays over time.

### 2. Workflow sequence learning (`internal/scoring/context_rules.go`)

Iris tracks sequential command pairs to learn common developer workflows (e.g. `git add` $\rightarrow$ `git commit`, `go build` $\rightarrow$ `./iris`). When a parent command skeleton matches the previous command, related suggestions receive a priority boost.

### 3. Skeleton extraction (`internal/scoring/skeleton.go`)

`ExtractSkeleton(cmd)` normalizes full command strings into structural skeletons by removing specific arguments and flags (e.g. `git commit -m "feat: test"` $\rightarrow$ `git commit`).

### 4. Spec priority & flag gating (`spec/lookup.go`)

Within spec completion mode:
- Files and subcommands default to standard priority (`Priority = 30`).
- Flags and options default to low priority (`Priority = 10`) when typing arguments.
- When the user explicitly types `-` or `--`, flags are promoted (`Priority = 80`).