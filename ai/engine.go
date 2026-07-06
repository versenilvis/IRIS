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
	handler   AIHandler
	cache     *ProviderCache
	providers []ContextProvider
}

func NewAIEngine(h AIHandler) *AIEngine {
	if h == nil {
		h = defaultAIHandler
	}
	return &AIEngine{
		handler: h,
		cache:   NewProviderCache(4 * time.Second),
		providers: []ContextProvider{
			DockerExecProvider{},
			GitCommitProvider{},
		},
	}
}

func (e *AIEngine) RegisterProvider(p ContextProvider) {
	e.providers = append([]ContextProvider{p}, e.providers...)
}

func (e *AIEngine) GatherDynamicContext(ctx context.Context, buf string) string {
	for _, p := range e.providers {
		if p.Matches(buf) {
			return e.cache.GetOrGather(ctx, p)
		}
	}
	return ""
}

func (e *AIEngine) Cache() *ProviderCache {
	return e.cache
}

func (e *AIEngine) Suggest(ctx context.Context, buf string, env EnvSnapshot, dynamicCtx string) (*spec.Suggestion, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if dynamicCtx == "" {
		dynamicCtx = e.GatherDynamicContext(ctx, buf)
	}
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

	if strings.HasPrefix(buf, "docker exec") {
		container := "my-container"
		if dynamicCtx != "" {
			lines := strings.Split(strings.TrimSpace(dynamicCtx), "\n")
			if len(lines) > 0 {
				fields := strings.Fields(lines[0])
				if len(fields) > 0 {
					container = fields[0]
				}
			}
		}
		return &spec.Suggestion{
			Cmd:        "docker exec -it " + container + " bash",
			Desc:       "ai suggestion",
			Icon:       "ai",
			Source:     string(SourceAI),
			Confidence: 85,
		}, nil
	}

	return nil, fmt.Errorf("no ai suggestion available")
}
