package integration

import (
	"bufio"
	"os"
	"path/filepath"
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
	ID  int
	Cmd string
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

		seen := make(map[string]bool)
		counter := 1
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			cmd := line

			if shellName == "zsh" {
				// zsh history format: : <timestamp>;<command>
				parts := strings.SplitN(line, ";", 2)
				if len(parts) == 2 {
					cmd = parts[1]
				}
			} else if shellName == "bash" {
				// bash history format can include timestamps starting with #
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
				// simple parse for fish history (YAML-like)
				if strings.HasPrefix(line, "- cmd: ") {
					cmd = strings.TrimPrefix(line, "- cmd: ")
				} else {
					continue
				}
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
