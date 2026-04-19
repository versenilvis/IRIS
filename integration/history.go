package integration

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/versenilvis/fuzzyvn"
	"github.com/versenilvis/iris/integration/shell"
)

var (
	historyCache  []string
	idMapCache    map[string]int
	searcherCache *fuzzyvn.Searcher
	mu            sync.Mutex
	lastModTime   int64
)

type HistResult struct {
	ID         int
	Cmd        string
	FuzzyScore int
}

func init() {
	idMapCache = make(map[string]int)
}

func SearchHistory(query string) ([]HistResult, error) {
	mu.Lock()
	defer mu.Unlock()

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	shellName := "bash"
	if shell.Current != nil {
		shellName = shell.Current.GetName()
	}

	var histFile string
	switch shellName {
	case "zsh":
		histFile = filepath.Join(home, ".zsh_history")
	case "fish":
		histFile = filepath.Join(home, ".local", "share", "fish", "fish_history")
	default:
		histFile = filepath.Join(home, ".bash_history")
	}

	if info, err := os.Stat(histFile); err == nil {
		if info.ModTime().UnixNano() > lastModTime {
			historyCache = nil // force reload
			idMapCache = make(map[string]int)
			lastModTime = info.ModTime().UnixNano()
		}
	}

	// lazy load history if cache is empty
	if len(historyCache) == 0 {
		file, err := os.Open(histFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		var allCmds []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			cmd := line

			if shellName == "zsh" {
				parts := strings.SplitN(line, ";", 2)
				if len(parts) == 2 {
					cmd = parts[1]
				}
			} else if shellName == "bash" {
				if strings.HasPrefix(line, "#") && len(line) > 1 {
					isTimestamp := true
					for _, c := range line[1:] {
						if c < '0' || c > '9' {
							isTimestamp = false
							break
						}
					}
					if isTimestamp {
						continue
					}
				}
			} else if shellName == "fish" {
				if strings.HasPrefix(line, "- cmd: ") {
					cmd = strings.TrimPrefix(line, "- cmd: ")
				} else {
					continue
				}
			}

			cmd = strings.TrimSpace(cmd)
			if cmd != "" {
				allCmds = append(allCmds, cmd)
			}
		}

		// build historyCache backwards so newest commands come first
		seen := make(map[string]bool)
		for i := len(allCmds) - 1; i >= 0; i-- {
			cmd := allCmds[i]
			if !seen[cmd] {
				historyCache = append(historyCache, cmd)
				seen[cmd] = true
				// we assign the ID as the original line number (1-indexed based on allCmds length)
				idMapCache[cmd] = i + 1
			}
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
			cmd := historyCache[i]
			results = append(results, HistResult{
				ID:  idMapCache[cmd],
				Cmd: cmd,
			})
		}
		return results, nil
	}

	matches := searcherCache.SearchWithScores(query, &fuzzyvn.SearchOptions{Limit: 100})

	var results []HistResult
	for _, m := range matches {
		results = append(results, HistResult{
			ID:         idMapCache[m.Str],
			Cmd:        m.Str,
			FuzzyScore: m.Score,
		})
	}

	// within the same tier, we sort by ID
	// tier 1: exact match
	// tier 2: prefix match
	// tier 3: substring match
	// tier 4: fuzzy match
	getTier := func(cmd, q string) int {
		cmdLow := strings.ToLower(cmd)
		qLow := strings.ToLower(q)
		if cmdLow == qLow {
			return 1
		}
		if strings.HasPrefix(cmdLow, qLow) {
			return 2
		}
		if strings.Contains(cmdLow, qLow) {
			return 3
		}
		return 4
	}

	sort.SliceStable(results, func(i, j int) bool {
		tI := getTier(results[i].Cmd, query)
		tJ := getTier(results[j].Cmd, query)
		if tI != tJ {
			return tI < tJ // lower tier is better
		}
		
		// If both are fuzzy matches (Tier 4), prioritize fuzzyvn score first!
		if tI == 4 && results[i].FuzzyScore != results[j].FuzzyScore {
			return results[i].FuzzyScore > results[j].FuzzyScore
		}
		
		// if same tier (and same fuzzy score), sort by ID descending (most recent first)
		return results[i].ID > results[j].ID
	})

	return results, nil
}
