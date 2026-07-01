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

	// search history if in history mode
	var histResults []integration.HistResult
	if mode == "history" {
		aliases := core.GetAliasesCopy()
		histResults, _ = integration.SearchHistory(query, aliases)
	}

	normalizedQuery := strings.TrimSpace(query)

	// add suggestion helper to deduplicate
	addSuggestion := func(s core.Suggestion) {
		normalizedCmd := strings.TrimSpace(s.Cmd)
		if normalizedCmd == "" {
			return
		}
		// filter exact match to avoid loops and redundant suggestions
		if normalizedCmd == normalizedQuery {
			return
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
		// spec mode: spec/alias only
		for _, s := range cmdResults {
			addSuggestion(s)
		}
	}

	if len(deduped) > maxSugg {
		deduped = deduped[:maxSugg]
	}
	return deduped
}
