# Changes from upstream (versenilvis/iris)

This fork adds multi-provider AI support and security hardening.

## Features Added

### Anthropic Native Client
- New `internal/ai/client_anthropic.go` — implements Anthropic's Messages API natively
  - System prompt as top-level field (not a message)
  - `x-api-key` header auth (not Bearer)
  - Anthropic-specific response parsing (`content[0].text`)
  - `anthropic-version: 2023-06-01` header
  - Same timeout handling, error capping, and rate limiting as the OpenAI client
- Wired into client factory via `inherited_from = "anthropic"` in provider config

### Multi-Provider Config
- `iris config init` now generates 7 commented-out provider examples:
  - OpenRouter, Groq, Ollama, Anthropic, DeepSeek, LM Studio, OpenAI
- Each shows correct `inherited_from` protocol and endpoint
- Docs updated with full multi-provider sample config
- Backward compatible: existing Groq/Ollama configs still work

### Dangerous Command Detection
- 16 regex patterns covering destructive shell commands
- Patterns: `rm -rf /`, `dd if=`, `mkfs`, `chmod 777`, fork bombs,
  `curl | sh`, `git push --force`, `git reset --hard`, `eval`,
  docker injection, `chown root`, overwrite block devices and /etc
- Added `Dangerous bool` field to `spec.Suggestion`
- Both OpenAI and Anthropic clients set this flag on dangerous suggestions
- Overlay can surface a warning icon (⚠️) for flagged suggestions

## Security Hardening

| # | Issue | Fix |
|---|-------|-----|
| 1 | Default HTTP client had no global timeout | Added `Timeout: 10s` to `sharedHTTPClient` |
| 2 | Updater used `curl \| sh` (pipe-to-shell anti-pattern) | Now downloads→validates shebang→temp file→executes. Refuses non-script responses. 1MB cap. |
| 3 | Debug mode logs all keystrokes including passwords/secrets | Added regex-based sanitizer: redacts API keys, tokens, passwords, auth headers, `export SECRET=...` from DEBUG-level log messages |
| 4 | AI suggestions could contain destructive commands without warning | Added 16-pattern dangerous command detector (see above) |

## Files Changed

```
 M docs/README.md                       # Multi-provider docs
 M internal/ai/client.go                # HTTP timeout, dangerous flag
 M internal/ai/utils.go                 # Dangerous command detector
 M internal/logger/logger.go            # Debug log sanitizer
 M root/config_cmd.go                  # Config init with 7 providers
 M root/update.go                      # Safe updater (download→validate→run)
 M spec/spec.go                        # Dangerous bool field
 A internal/ai/client_anthropic.go      # NEW: Anthropic native client
```

## Test Results

- All 22 existing tests pass (no regressions)
- Anthropic client: connects and authenticates successfully to api.anthropic.com
- Dangerous detector: 15/15 destructive commands flagged, 5/5 safe commands pass clean

## How to Submit Upstream

1. Go to https://github.com/Opfour/iris
2. Click **Contribute** → **Open pull request**
3. Base repository: `versenilvis/iris` / base: `main`
4. Title suggestion: `feat: multi-provider AI support + security hardening`
5. Link to this file in the PR description