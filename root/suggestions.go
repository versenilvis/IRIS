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
		if normalizedCmd == "" {
			return
		}
		if s.Source != "alias" && s.Source != "abbr" && normalizedCmd == normalizedQuery {
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

		for i, h := range histResults {
			var conf int
			switch {
			case h.Tier <= 2:
				conf = max(75-(i*2), 65)
			case h.Tier == 3:
				conf = max(70-(i*2), 60)
			default:
				conf = max(45-(i*2), 30)
			}

			icon := "history"
			if h.Source == "atuin" {
				icon = "atuin"
			}

			addSuggestion(spec.Suggestion{
				Cmd:        h.Cmd,
				Desc:       h.Source,
				Icon:       icon,
				Source:     h.Source,
				Confidence: conf,
			})
		}

		// only boost a spec suggestion above history when the user has typed
		// enough of it to be specific (coverage >= 70%). this avoids short queries
		// like "git" pulling up system binaries above real history commands.
		if normalizedQuery != "" {
			qLen := len(normalizedQuery)
			for _, s := range cmdResults {
				if s.Cmd == "" || s.Icon == "system" || s.Desc == "system command" {
					continue
				}
				cLen := len(s.Cmd)
				if cLen == 0 {
					continue
				}
				coverage := float64(qLen) / float64(cLen)
				if coverage >= 0.70 && strings.HasPrefix(strings.ToLower(s.Cmd), strings.ToLower(normalizedQuery)) {
					s.Confidence = 85
					if s.Source == "" {
						s.Source = "spec"
					}
					addSuggestion(s)
				}
			}
		}
	}

	for _, s := range cmdResults {
		addSuggestion(s)
	}

	injectAISuggestion(&deduped, seen, normalizedQuery)

	if mode == "history" && normalizedQuery == "" {
		sort.SliceStable(deduped, func(i, j int) bool {
			return deduped[i].Confidence > deduped[j].Confidence
		})
		if aiSugg := GetCurrentAISuggestion(); aiSugg != nil {
			for i, s := range deduped {
				if s.Cmd == aiSugg.Cmd && s.Source == "ai" {
					if i > 0 {
						deduped = append([]spec.Suggestion{s}, append(deduped[:i], deduped[i+1:]...)...)
					}
					break
				}
			}
		}
		if len(deduped) > maxSugg {
			return deduped[:maxSugg]
		}
		return deduped
	}

	var finalResults []spec.Suggestion
	if mode == "history" {
		sort.SliceStable(deduped, func(i, j int) bool {
			return deduped[i].Confidence > deduped[j].Confidence
		})
		finalResults = deduped
	} else {
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

		finalResults = make([]spec.Suggestion, 0, len(scored))
		for _, sc := range scored {
			finalResults = append(finalResults, sc.Suggestion)
		}
	}

	if aiSugg := GetCurrentAISuggestion(); aiSugg != nil {
		for i, s := range finalResults {
			if s.Cmd == aiSugg.Cmd && s.Source == "ai" {
				if i > 0 {
					finalResults = append([]spec.Suggestion{s}, append(finalResults[:i], finalResults[i+1:]...)...)
				}
				break
			}
		}
	}

	if len(finalResults) > maxSugg {
		return finalResults[:maxSugg]
	}
	return finalResults
}

func injectAISuggestion(deduped *[]spec.Suggestion, seen map[string]bool, normalizedQuery string) {
	if aiSugg := GetCurrentAISuggestion(); aiSugg != nil {
		normalizedCmd := strings.TrimSpace(aiSugg.Cmd)
		if normalizedCmd == "" || normalizedCmd == normalizedQuery {
			return
		}
		if normalizedQuery != "" && !strings.HasPrefix(strings.ToLower(normalizedCmd), strings.ToLower(normalizedQuery)) {
			// allow in-buffer fix when subcommand or flags differ but root command matches
			qTokens := strings.Fields(normalizedQuery)
			cTokens := strings.Fields(normalizedCmd)
			if len(qTokens) == 0 || len(cTokens) == 0 || !strings.EqualFold(qTokens[0], cTokens[0]) {
				return
			}
		}

		if !seen[aiSugg.Cmd] {
			seen[aiSugg.Cmd] = true
			*deduped = append(*deduped, *aiSugg)
		} else {
			for i, item := range *deduped {
				if item.Cmd == aiSugg.Cmd {
					(*deduped)[i] = *aiSugg
					break
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
