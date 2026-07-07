package ai

import (
	"context"
	"strings"
	"time"

	"github.com/versenilvis/iris/config"
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
		handler:   h,
		cache:     NewProviderCache(4 * time.Second),
		providers: []ContextProvider{},
	}
}

func (e *AIEngine) RegisterProvider(p ContextProvider) {
	e.providers = append(e.providers, p)
}

func (e *AIEngine) GatherDynamicContext(ctx context.Context, buf string, cwd string) string {
	for _, p := range e.providers {
		if p.Matches(buf) {
			return e.cache.GetOrGather(ctx, p)
		}
	}
	return e.cache.GetOrGather(ctx, &universalProvider{cwd: cwd, buf: buf})
}

func (e *AIEngine) Cache() *ProviderCache {
	return e.cache
}

func (e *AIEngine) Suggest(ctx context.Context, buf string, env EnvSnapshot, dynamicCtx string) (*spec.Suggestion, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if dynamicCtx == "" {
		dynamicCtx = e.GatherDynamicContext(ctx, buf, env.Cwd)
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
	if sugg != nil {
		sugg.Cmd = NormalizeSuggestion(buf, sugg.Cmd)
	}
	return sugg, nil
}

func NormalizeSuggestion(buf string, suggCmd string) string {
	suggCmd = CleanSuggestion(suggCmd)

	if strings.Contains(buf, "-m \"") || strings.Contains(buf, "-am \"") || strings.Contains(buf, "--message \"") {
		if !strings.HasPrefix(strings.ToLower(suggCmd), strings.ToLower(buf)) {
			for _, flag := range []string{"-m ", "-am ", "--message "} {
				idx := strings.Index(suggCmd, flag)
				if idx != -1 {
					afterFlag := suggCmd[idx+len(flag):]
					if !strings.HasPrefix(afterFlag, "\"") && !strings.HasPrefix(afterFlag, "'") {
						suggCmd = suggCmd[:idx+len(flag)] + "\"" + afterFlag + "\""
						break
					}
				}
			}
		}
	}

	return suggCmd
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
	cfg := config.Get()
	if !cfg.AI.Enabled {
		return nil, nil
	}
	pCfg, ok := cfg.AI.GetActiveProvider()
	if !ok {
		return nil, nil
	}
	client, err := NewClient(pCfg)
	if err != nil {
		return nil, err
	}
	return client.Suggest(ctx, buf, env, dynamicCtx)
}
