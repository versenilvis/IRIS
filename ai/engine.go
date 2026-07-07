package ai

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/versenilvis/iris/config"
	"github.com/versenilvis/iris/spec"
)

type suggestionCacheItem struct {
	query string
	sugg  *spec.Suggestion
	time  time.Time
}

type AIEngine struct {
	handler        AIHandler
	cache          *ProviderCache
	providers      []ContextProvider
	lastCallTime   time.Time
	lastSuggestion *suggestionCacheItem
	mu             sync.Mutex
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
	trimmed := strings.TrimSpace(buf)
	if len(trimmed) < 3 {
		return nil, nil
	}

	e.mu.Lock()
	if e.lastSuggestion != nil && time.Since(e.lastSuggestion.time) < 30*time.Second {
		lastQ := e.lastSuggestion.query
		lastCmd := e.lastSuggestion.sugg.Cmd
		if buf == lastQ || (strings.HasPrefix(buf, lastQ) && strings.HasPrefix(strings.ToLower(lastCmd), strings.ToLower(buf))) || (strings.HasPrefix(lastQ, buf) && strings.HasPrefix(strings.ToLower(lastCmd), strings.ToLower(buf))) {
			cached := *e.lastSuggestion.sugg
			e.mu.Unlock()
			return &cached, nil
		}
	}

	minIntervalMS := config.Get().AI.MinIntervalMS
	if minIntervalMS <= 0 {
		minIntervalMS = 3000
	}
	if !e.lastCallTime.IsZero() && time.Since(e.lastCallTime) < time.Duration(minIntervalMS)*time.Millisecond {
		e.mu.Unlock()
		return nil, nil
	}
	e.lastCallTime = time.Now()
	e.mu.Unlock()

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
		e.mu.Lock()
		e.lastSuggestion = &suggestionCacheItem{
			query: buf,
			sugg:  sugg,
			time:  time.Now(),
		}
		e.mu.Unlock()
	}
	return sugg, nil
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
