# Iris developer & architecture documentation

This directory contains code architecture guides, engine design notes, and development instructions for Iris contributors.

## Table of contents

- [Development Workflow & Runner](development.md): Environment setup, building, testing, linters, and modular `just` runner (`justfiles/`)
- [AI Engine Architecture](ai.md): Internal AI completion system (`internal/ai`), providers, prompt design, and debounce strategy
- [Scoring & Frecency Engine](scoring.md): Scoring mechanics (`internal/scoring`), frecency calculations, context rules, and skeleton extraction
- [Overlay & Rendering Pipeline](overlay.md): PTY overlay positioning, cell width calculation, and terminal rendering
- [PTY Bridge & Keystroke Wrapper](root.md): Keystroke interception loop, raw mode PTY handling, and IPC scanner
- [Spec & Completion Engine](spec.md): Static specifications, priority-based flag gating, and Cobra `__complete` dynamic completion
- [File & Path Generator](filegen.md): File system traversal, extension filtering, and directory slash preservation
- [History Provider](history.md): Shell history indexing, search algorithms, and caching
- [Auto Updater](updater.md): Release tracking, version comparison, and atomic binary updates
