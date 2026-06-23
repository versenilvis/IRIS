package root

import (
	"fmt"
	"strings"

	"github.com/versenilvis/iris/commands/core"
	"github.com/versenilvis/iris/config"
	"github.com/versenilvis/iris/integration"
)

// mergeResults collects and dedupes suggestions for a query and mode
// example: mergeResults("git ", "spec")
func MergeResults(query string, mode string) []core.Suggestion {
	maxSugg := config.Get().UI.MaxSuggestions
	seen := make(map[string]bool)
	deduped := []core.Suggestion{}

	// always call lookup to scan aliases and get spec suggestions
	var cmdResults []core.Suggestion
	if query != "" {
		debugLog("[Merge] Calling Lookup for '%s'", query)
		cmdResults = core.Lookup(query)
	}

	// search history if in history mode or if query is not empty
	var histResults []integration.HistResult
	if mode == "history" || query != "" {
		aliases := core.GetAliasesCopy()
		histResults, _ = integration.SearchHistory(query, aliases)
	}

	// add suggestion helper to deduplicate
	addSuggestion := func(s core.Suggestion) {
		normalizedCmd := strings.TrimSpace(s.Cmd)
		if normalizedCmd == "" {
			return
		}
		// filter exact match for spec commands to avoid loops
		if s.Desc != " history" && !strings.HasPrefix(s.Desc, "alias:") {
			normalizedQuery := strings.TrimSpace(query)
			if normalizedCmd == normalizedQuery {
				return
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
			addSuggestion(core.Suggestion{
				Cmd:  h.Cmd,
				Desc: " history",
				Icon: fmt.Sprintf("%d", h.ID),
			})
		}
		for _, s := range cmdResults {
			addSuggestion(s)
		}
	} else {
		// spec mode: spec/alias first, then history
		for _, s := range cmdResults {
			addSuggestion(s)
		}
		for _, h := range histResults {
			addSuggestion(core.Suggestion{
				Cmd:  h.Cmd,
				Desc: " history",
				Icon: fmt.Sprintf("%d", h.ID),
			})
		}
	}

	if len(deduped) > maxSugg {
		deduped = deduped[:maxSugg]
	}
	return deduped
}
