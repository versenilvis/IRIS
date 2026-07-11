package ai

import (
	"context"
	"time"

	"github.com/versenilvis/iris/spec"
)

type RuleBasedSuggester struct{}

type EmptyLineRule struct {
	Name    string
	Match   func(env EnvSnapshot) bool
	Suggest func(env EnvSnapshot) *spec.Suggestion
}

func (s RuleBasedSuggester) SuggestOnEmpty(ctx context.Context, env EnvSnapshot) (*spec.Suggestion, error) {
	for _, rule := range DefaultEmptyLineRules {
		if rule.Match(env) {
			return rule.Suggest(env), nil
		}
	}
	return nil, nil
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
