package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/versenilvis/iris/spec"
)

type AIHandler func(ctx context.Context, buf string, env EnvSnapshot, dynamicCtx string) (*spec.Suggestion, error)

type AIEngine struct {
	handler AIHandler
}

func NewAIEngine(h AIHandler) *AIEngine {
	if h == nil {
		h = defaultAIHandler
	}
	return &AIEngine{handler: h}
}

func (e *AIEngine) Suggest(ctx context.Context, buf string, env EnvSnapshot, dynamicCtx string) (*spec.Suggestion, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	sugg, err := e.handler(ctx, buf, env, dynamicCtx)
	if err != nil {
		return nil, err
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return sugg, nil
}

func ShouldOverwrite(originalBuf string, currentBuf string, newSugg *spec.Suggestion, currentConfidence int) bool {
	if newSugg == nil {
		return false
	}
	if !strings.HasPrefix(currentBuf, originalBuf) {
		return false
	}
	if !strings.HasPrefix(strings.ToLower(newSugg.Cmd), strings.ToLower(currentBuf)) {
		return false
	}
	return newSugg.Confidence > currentConfidence
}

func defaultAIHandler(ctx context.Context, buf string, env EnvSnapshot, dynamicCtx string) (*spec.Suggestion, error) {
	select {
	case <-time.After(100 * time.Millisecond):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if strings.HasPrefix(buf, "git commit -m \"") || strings.HasPrefix(buf, "git commit -m '") {
		msg := "feat: update codebase"
		if dynamicCtx != "" {
			msg = "feat: " + strings.TrimSpace(dynamicCtx)
		} else if env.GitStatus != "" {
			msg = "chore: update " + strings.TrimSpace(env.GitStatus)
		}
		quote := string(buf[len("git commit -m ")])
		return &spec.Suggestion{
			Cmd:        "git commit -m " + string(quote) + msg + string(quote),
			Desc:       "ai suggestion",
			Icon:       "ai",
			Source:     string(SourceAI),
			Confidence: 85,
		}, nil
	}

	return nil, fmt.Errorf("no ai suggestion available")
}
