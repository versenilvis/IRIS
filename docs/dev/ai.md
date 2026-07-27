# AI engine architecture (`internal/ai/`)

The AI subsystem provides real-time, context-aware command completions powered by cloud or local LLM providers (Groq, Ollama).

## Core components

- `internal/ai/client.go`: Defines the `Client` interface and handles HTTP requests to OpenAI-compatible chat completion endpoints.
- `internal/ai/env.go`: Captures runtime context snapshots (`EnvSnapshot`) including CWD, previous command, exit code, and recent history.
- `internal/ai/prompts.go`: Constructs system and user prompts optimized for shell command completion.
- `root/wrapper.go`: Manages async debounce timers, request cancellation (`context.WithCancel`), and ghost text injection.

## Provider interface

All providers use a unified HTTP pattern matching OpenAI's `/v1/chat/completions` format:

```go
type Client interface {
    Suggest(ctx context.Context, prompt string, env *EnvSnapshot, currentCmd string) (*AISuggestion, error)
}
```

## Request lifecycle

1. User types in the prompt buffer.
2. `root/wrapper.go` triggers a debounce timer (`debounce_ms`, default 500ms).
3. If typing continues, previous in-flight contexts are cancelled (`aiCancel()`).
4. An `EnvSnapshot` is created capturing CWD, last executed command, and up to 3 recent history entries.
5. Provider sends an HTTP POST request with structured JSON payload.
6. Response is parsed and rendered as inline ghost text via `overlay.InjectAISuggestion()`.
