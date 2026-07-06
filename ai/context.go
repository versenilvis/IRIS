package ai

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/versenilvis/iris/spec"
)

type ContextSuggester interface {
	SuggestOnEmpty(ctx context.Context, env EnvSnapshot) (*spec.Suggestion, error)
}

type RuleBasedSuggester struct{}

func (s RuleBasedSuggester) SuggestOnEmpty(ctx context.Context, env EnvSnapshot) (*spec.Suggestion, error) {
	if env.LastExitCode != 0 && env.LastCmd != "" {
		return &spec.Suggestion{
			Cmd:        env.LastCmd,
			Desc:       "retry failed command",
			Icon:       "retry",
			Source:     string(SourceSpec),
			Confidence: 80,
		}, nil
	}

	if strings.TrimSpace(env.LastCmd) == "git status" {
		return &spec.Suggestion{
			Cmd:        "git diff",
			Desc:       "view modifications",
			Icon:       "git",
			Source:     string(SourceSpec),
			Confidence: 75,
		}, nil
	}

	if env.GitStatus != "" {
		return &spec.Suggestion{
			Cmd:        "git status",
			Desc:       "check repository state",
			Icon:       "git",
			Source:     string(SourceSpec),
			Confidence: 70,
		}, nil
	}

	if strings.Contains(env.DirSignature, "package.json") {
		return &spec.Suggestion{
			Cmd:        "npm run dev",
			Desc:       "start dev server",
			Icon:       "npm",
			Source:     string(SourceSpec),
			Confidence: 65,
		}, nil
	}

	return nil, nil
}

type ContextCache struct {
	mu               sync.Mutex
	lastSnapshotHash string
	lastSuggestion   *spec.Suggestion
	lastFetchedAt    time.Time
}

func NewContextCache() *ContextCache {
	return &ContextCache{}
}

func (c *ContextCache) ShouldCallAI(snap EnvSnapshot, minInterval time.Duration) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	hash := snap.Hash()
	if hash == c.lastSnapshotHash {
		return false
	}
	if time.Since(c.lastFetchedAt) < minInterval {
		return false
	}
	return true
}

func (c *ContextCache) GetCachedSuggestion(snap EnvSnapshot) *spec.Suggestion {
	c.mu.Lock()
	defer c.mu.Unlock()
	if snap.Hash() == c.lastSnapshotHash {
		return c.lastSuggestion
	}
	return nil
}

func (c *ContextCache) Update(snap EnvSnapshot, sugg *spec.Suggestion) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastSnapshotHash = snap.Hash()
	c.lastSuggestion = sugg
	c.lastFetchedAt = time.Now()
}

func (c *ContextCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastSnapshotHash = ""
	c.lastSuggestion = nil
}

type EmptyLinePredictor struct {
	ruleSuggester ContextSuggester
	aiSuggester   ContextSuggester
	cache         *ContextCache
	minInterval   time.Duration
}

func NewEmptyLinePredictor(rule ContextSuggester, ai ContextSuggester, minInterval time.Duration) *EmptyLinePredictor {
	if minInterval == 0 {
		minInterval = 2 * time.Second
	}
	if rule == nil {
		rule = RuleBasedSuggester{}
	}
	return &EmptyLinePredictor{
		ruleSuggester: rule,
		aiSuggester:   ai,
		cache:         NewContextCache(),
		minInterval:   minInterval,
	}
}

func (p *EmptyLinePredictor) Predict(ctx context.Context, env EnvSnapshot, aiEnabled bool) (*spec.Suggestion, error) {
	if p.ruleSuggester != nil {
		sugg, err := p.ruleSuggester.SuggestOnEmpty(ctx, env)
		if err == nil && sugg != nil && sugg.Confidence >= 70 {
			return sugg, nil
		}
	}

	if !aiEnabled || p.aiSuggester == nil {
		return nil, nil
	}

	if !p.cache.ShouldCallAI(env, p.minInterval) {
		return p.cache.GetCachedSuggestion(env), nil
	}

	sugg, err := p.aiSuggester.SuggestOnEmpty(ctx, env)
	if err != nil || ctx.Err() != nil {
		return nil, err
	}

	p.cache.Update(env, sugg)
	return sugg, nil
}

func (p *EmptyLinePredictor) Cache() *ContextCache {
	return p.cache
}
