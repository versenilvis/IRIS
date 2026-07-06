package root

import (
	"strings"
	"sync"

	"github.com/versenilvis/iris/ai"
	"github.com/versenilvis/iris/config"
	"github.com/versenilvis/iris/integration"
	"github.com/versenilvis/iris/logger"
	"github.com/versenilvis/iris/spec"
)

// MergeResults collects and dedupes suggestions for a query and mode
func MergeResults(query string, mode string) []spec.Suggestion {
	maxSugg := config.Get().UI.MaxSuggestions
	seen := make(map[string]bool)
	deduped := []spec.Suggestion{}

	// always call lookup to scan aliases and get spec suggestions
	logger.Debugf("Merge Calling Lookup for '%s'", query)
	cmdResults := spec.Lookup(query)

	// search history if in history mode
	var histResults []integration.HistResult
	if mode == "history" {
		aliases := spec.GetAliasesCopy()
		histResults, _ = integration.SearchHistory(query, aliases)
	}

	normalizedQuery := strings.TrimSpace(query)

	// add suggestion helper to deduplicate
	addSuggestion := func(s spec.Suggestion) {
		normalizedCmd := strings.TrimSpace(s.Cmd)
		if normalizedCmd == "" {
			return
		}
		// filter exact match to avoid loops and redundant suggestions
		if normalizedCmd == normalizedQuery {
			return
		}
		if s.Source == "" {
			s.Source = "spec"
			if s.Confidence == 0 {
				s.Confidence = 50
			}
		}
		if !seen[s.Cmd] {
			seen[s.Cmd] = true
			deduped = append(deduped, s)
		}
	}

	if mode == "history" {
		// history mode: history first, then spec/alias
		for _, h := range histResults {
			addSuggestion(spec.Suggestion{
				Cmd:        h.Cmd,
				Desc:       "history",
				Icon:       "history",
				Source:     "history",
				Confidence: 70,
			})
		}
		for _, s := range cmdResults {
			addSuggestion(s)
		}
	} else {
		// spec mode: spec/alias only
		for _, s := range cmdResults {
			addSuggestion(s)
		}
	}

	if aiSugg := GetCurrentAISuggestion(); aiSugg != nil {
		if strings.HasPrefix(strings.ToLower(aiSugg.Cmd), strings.ToLower(normalizedQuery)) {
			addSuggestion(*aiSugg)
		}
	}

	if len(deduped) > maxSugg {
		deduped = deduped[:maxSugg]
	}
	return deduped
}

var (
	aiEngine     *ai.AIEngine
	aiEngineOnce sync.Once
)

func GetAIEngine() *ai.AIEngine {
	aiEngineOnce.Do(func() {
		aiEngine = ai.NewAIEngine(nil)
	})
	return aiEngine
}

var (
	currentAISugg *spec.Suggestion
	aiSuggMu      sync.RWMutex
)

func SetCurrentAISuggestion(sugg *spec.Suggestion) {
	aiSuggMu.Lock()
	defer aiSuggMu.Unlock()
	currentAISugg = sugg
}

func GetCurrentAISuggestion() *spec.Suggestion {
	aiSuggMu.RLock()
	defer aiSuggMu.RUnlock()
	return currentAISugg
}
