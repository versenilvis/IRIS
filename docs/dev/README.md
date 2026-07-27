# Iris developer & architecture documentation

This directory contains code architecture guides, engine design notes, and development instructions for Iris contributors.

## Table of contents

- [Development workflow & runner](development.md): Environment setup, building, testing, linters, and modular `just` runner (`justfiles/`)
- [AI engine architecture](ai.md): Internal AI completion system (`internal/ai`), providers, prompt design, and debounce strategy
- [Scoring & frecency engine](scoring.md): Scoring mechanics (`internal/scoring`), frecency calculations, context rules, and skeleton extraction
- [Overlay & rendering pipeline](overlay.md): PTY overlay positioning, cell width calculation, and terminal rendering
- [PTY bridge & keystroke wrapper](root.md): Keystroke interception loop, raw mode PTY handling, and IPC scanner
- [Spec & completion engine](spec.md): Static specifications, priority-based flag gating, and Cobra `__complete` dynamic completion
- [File & path generator](filegen.md): File system traversal, extension filtering, and directory slash preservation
- [History provider](history.md): Shell history indexing, search algorithms, and caching
- [Auto updater](updater.md): Release tracking, version comparison, and atomic binary updates
