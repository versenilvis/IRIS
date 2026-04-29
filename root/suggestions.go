package root

import (
	"fmt"
	"strings"

	"github.com/versenilvis/iris/commands/core"
	"github.com/versenilvis/iris/integration"
)

// mergeResults collects and dedupes suggestions for a query and mode
// example: mergeResults("git ", "spec")
func mergeResults(query string, mode string) []core.Suggestion {
	if query == "" && mode != "history" {
		debugLog("[Merge] Query empty, returning nil")
		return nil
	}

	normalizedQuery := strings.TrimSpace(query)
	seen := make(map[string]bool)
	deduped := []core.Suggestion{}

	if mode == "history" {
		histResults, _ := integration.SearchHistory(query)
		for _, h := range histResults {
			normalizedCmd := strings.TrimSpace(h.Cmd)
			if seen[normalizedCmd] {
				continue
			}
			seen[normalizedCmd] = true
			deduped = append(deduped, core.Suggestion{
				Cmd:  h.Cmd,
				Desc: " history",
				Icon: fmt.Sprintf("%d", h.ID),
			})
			if len(deduped) >= 100 {
				break
			}
		}
		debugLog("[Merge] History mode found %d items", len(deduped))
		return deduped
	}

	debugLog("[Merge] Calling Lookup for '%s'", query)
	cmdResults := core.Lookup(query)
	debugLog("[Merge] Lookup returned %d raw items", len(cmdResults))

	for _, s := range cmdResults {
		normalizedCmd := strings.TrimSpace(s.Cmd)
		if normalizedCmd == normalizedQuery { // filter exact matches to avoid loops
			debugLog("[Merge] Filtered EXACT MATCH: '%s'", normalizedCmd)
			continue
		}

		if !seen[s.Cmd] {
			seen[s.Cmd] = true
			deduped = append(deduped, s)
		}
	}
	if len(deduped) > 100 {
		deduped = deduped[:100]
	}
	return deduped
}
