package intergration

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/versenilvis/fuzzyvn"
)

var (
	historyCache  []string
	idMapCache    map[string]int
	searcherCache *fuzzyvn.Searcher
	mu            sync.Mutex
)

type HistResult struct {
	ID  int
	Cmd string
}

func init() {
	idMapCache = make(map[string]int)
}

// SearchHistory reads the zsh history and searches using fuzzyvn (my own pkg)
func SearchHistory(query string) ([]HistResult, error) {
	mu.Lock()
	defer mu.Unlock()

	// lazy load history if cache is empty
	if len(historyCache) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		histFile := filepath.Join(home, ".zsh_history")
		file, err := os.Open(histFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		seen := make(map[string]bool)
		counter := 1
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			// zsh history format: : <timestamp>;<command>
			parts := strings.SplitN(line, ";", 2)
			cmd := line
			if len(parts) == 2 {
				cmd = parts[1]
			}

			cmd = strings.TrimSpace(cmd)
			if cmd != "" {
				if !seen[cmd] {
					historyCache = append(historyCache, cmd)
					seen[cmd] = true
				}
				idMapCache[cmd] = counter
			}
			counter++
		}
		searcherCache = fuzzyvn.NewPlainSearcher(historyCache)
	}

	if query == "" {
		var results []HistResult
		limit := 100
		if len(historyCache) < limit {
			limit = len(historyCache)
		}

		for i := 0; i < limit; i++ {
			cmd := historyCache[len(historyCache)-1-i]
			results = append(results, HistResult{
				ID:  idMapCache[cmd],
				Cmd: cmd,
			})
		}
		return results, nil
	}

	matches := searcherCache.Search(query, &fuzzyvn.SearchOptions{Limit: 100})

	var results []HistResult
	for _, m := range matches {
		results = append(results, HistResult{
			ID:  idMapCache[m],
			Cmd: m,
		})
	}

	return results, nil
}
