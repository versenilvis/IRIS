# Directory-Scoped Ghost Text Prediction (`internal/ctxcheck` & `root/wrapper.go`)

This document describes the design and implementation of context-aware, directory-scoped ghost text prediction in IRIS.

## Overview

Iris predicts commands based on shell history, transitions, and frecency. To prevent cross-project command leakage (e.g. suggesting `just reload` in projects without a `justfile`), predictions are evaluated through a tiered scoping model and context validation engine.

## Core Concepts

### 1. Project Scoping (`project_id`)

A workspace scope is identified by `project_id`, calculated by `workspace.DetectProjectIDCached(cwd)`:
- Traverses upward from `cwd` searching for `.git` (or project markers).
- Canonicalizes paths resolving symlinks (`/var/folders` to `/private/var` on macOS).
- Commands in non-git directories fall back to directory path (`COALESCE(project_id, cwd)`).

### 2. Candidate Tiering

Tiers are calculated in Go during candidate retrieval (`internal/scoring/frecency.go`):
- **Tier 4 (Local CWD)**: Exact directory match (`cwd == current_cwd`).
- **Tier 3 (Subdir to Ancestor)**: Current directory is descendant of command's recorded working directory within the same project.
- **Tier 2 (Ancestor to Subdir)**: Current directory is ancestor of command's recorded working directory within the same project.
- **Tier 1 (Project Siblings)**: Same project (`project_id`), different directory branch.
- **Tier 0 (Foreign / Global)**: Different project or non-matching directories.

### 3. Lazy Scope Count (`ScopeCount`)

Rather than running expensive `COUNT(DISTINCT)` aggregates across all candidates in the global SQL query, `ScopeCount` is evaluated lazily only when candidate validation reaches Tier 0:
- Commands with `Tier > 0` bypass scope count evaluation.
- Commands at `Tier 0` query the number of distinct scopes (`project_id` or `cwd`) where the command was executed.

### 4. Validation Engine & Verdicts (`internal/ctxcheck`)

The validator parses candidate commands against the local filesystem with a 15ms deadline:
- **`Valid`**: Contextually valid in current directory (e.g. recipe exists in `justfile`, script in `package.json`, target in `Makefile`, executable/directory exists).
- **`Invalid`**: Target explicitly missing (e.g. `just non_existent_recipe`, `cd nonexistent_dir`). Blocked across all tiers.
- **`Free`**: Generic shell commands or commands without known manifests (`cargo run`, `git status`, `./app`).
- **`Unknown`**: Indeterminate or timed-out parsing (e.g. complex compound commands with `cd`, network-mounted slow disks).

### 5. Admission Matrix (`ctxcheck.Allow`)

| Tier | Verdict `Valid` | Verdict `Invalid` | Verdict `Free` | Verdict `Unknown` |
| :--- | :--- | :--- | :--- | :--- |
| **Tier 4** (Exact CWD) | Allowed | Blocked | Allowed | Allowed |
| **Tier 1–3** (Same Project) | Allowed | Blocked | Allowed | Blocked |
| **Tier 0** (Foreign Project) | Allowed | Blocked | ScopeCount >= 3 | Blocked |

### 6. Known Limitations

- **Non-existent destination targets**: Commands like `cp source ./dest` where `dest` does not yet exist are classified as `Invalid`.
- **Non-git directories**: If `cwd` is outside a git repository, cross-directory sharing relies on exact directory paths or global scope threshold (Case 2 requires `project_id`).
- **Shell aliases and functions**: Custom shell aliases or functions without matching binaries are treated as `Free`.
- **Nested monorepos**: Repositories containing nested `.git` directories or submodules treat each git boundary as a distinct project scope.

## Debugging & Observability

### `IRIS_DEBUG_PREDICT`

Set `IRIS_DEBUG_PREDICT=1` to log prediction candidate evaluations:
- Logs current directory, prefix, chosen prediction, and candidate breakdown (tier, scopes, verdict, allow).
- **Warning**: Log entries include verbatim command strings, which may contain sensitive arguments, tokens, or passwords.
- Only enable during dogfooding/troubleshooting sessions and remove log files after analysis.
