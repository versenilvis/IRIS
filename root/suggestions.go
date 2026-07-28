package root

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/versenilvis/iris/integration"
	"github.com/versenilvis/iris/internal/ai"
	"github.com/versenilvis/iris/internal/config"
	"github.com/versenilvis/iris/internal/logger"
	"github.com/versenilvis/iris/internal/scoring"
	"github.com/versenilvis/iris/spec"
)

// MergeResults collects and dedupes suggestions for a query and mode
func MergeResults(query string, mode string) []spec.Suggestion {
	maxSugg := config.Get().UI.MaxSuggestions
	seen := make(map[string]bool)
	deduped := []spec.Suggestion{}
	normalizedQuery := strings.TrimSpace(query)

	// add suggestion helper to deduplicate
	addSuggestion := func(s spec.Suggestion) {
		normalizedCmd := strings.TrimSpace(s.Cmd)
		if normalizedCmd == "" || normalizedCmd == normalizedQuery {
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

	// always call lookup to scan aliases and get spec suggestions
	logger.Debugf("Merge Calling Lookup for '%s'", query)
	cmdResults := spec.Lookup(query)

	if mode == "history" {
		aliases := spec.GetAliasesCopy()
		histResults, _ := integration.SearchHistory(query, aliases)

		baseConf := 75
		for i, h := range histResults {
			addSuggestion(spec.Suggestion{
				Cmd:        h.Cmd,
				Desc:       "history",
				Icon:       "history",
				Source:     "history",
				Confidence: max(baseConf-(i*2), 60),
			})
		}

		for _, s := range cmdResults {
			addSuggestion(s)
		}

		if normalizedQuery == "" {
			if len(deduped) > maxSugg {
				return deduped[:maxSugg]
			}
			return deduped
		}

		injectAISuggestion(&deduped, seen, normalizedQuery)

		sort.SliceStable(deduped, func(i, j int) bool {
			return deduped[i].Confidence > deduped[j].Confidence
		})

		if len(deduped) > maxSugg {
			return deduped[:maxSugg]
		}
		return deduped
	}

	for _, s := range cmdResults {
		addSuggestion(s)
	}

	injectAISuggestion(&deduped, seen, normalizedQuery)

	cwd := spec.GetCWD()
	tokens := spec.Tokenize(query)
	rootCmd := ""
	if len(tokens) > 0 {
		rootCmd = tokens[0]
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	store, _ := scoring.GetFrecencyStore()
	signals := scoring.CollectSignals(ctxTimeout, cwd, query, rootCmd, store, getPrevSkeleton())
	scored := scoring.Score(deduped, signals)

	finalResults := make([]spec.Suggestion, 0, len(scored))
	for _, sc := range scored {
		finalResults = append(finalResults, sc.Suggestion)
	}

	if len(finalResults) > maxSugg {
		return finalResults[:maxSugg]
	}
	return finalResults
}

func injectAISuggestion(deduped *[]spec.Suggestion, seen map[string]bool, normalizedQuery string) {
	if aiSugg := GetCurrentAISuggestion(); aiSugg != nil {
		normalizedCmd := strings.TrimSpace(aiSugg.Cmd)
		if normalizedCmd != "" && normalizedCmd != normalizedQuery && strings.HasPrefix(strings.ToLower(normalizedCmd), strings.ToLower(normalizedQuery)) {
			if !seen[aiSugg.Cmd] {
				seen[aiSugg.Cmd] = true
				*deduped = append(*deduped, *aiSugg)
			} else {
				for i, item := range *deduped {
					if item.Cmd == aiSugg.Cmd && aiSugg.Confidence > item.Confidence {
						(*deduped)[i].Confidence = aiSugg.Confidence
						if (*deduped)[i].Source == "" || (*deduped)[i].Source == "spec" || (*deduped)[i].Source == "history" {
							(*deduped)[i].Source = "ai"
						}
						break
					}
				}
			}
		}
	}
}

var (
	aiEngine     *ai.AIEngine
	aiEngineOnce sync.Once
)

func GetAIEngine() *ai.AIEngine {
	aiEngineOnce.Do(func() {
		aiEngine = ai.NewAIEngine(nil)
		for _, p := range ai.DefaultProviders {
			aiEngine.RegisterProvider(p)
		}
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
