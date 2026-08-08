// please note that zoxide also shows external suggestions
// at the end of the list on command mode
// they are the old directories that you have visited
// this is a feature, not a bug, and I want to keep it
package fs

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/versenilvis/fuzzy"
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "z",
		Description: "jump to directory",
		MaxArgs:     0,
		Generator:   ZoxideGenerator(),
	})
	spec.Register(&spec.Spec{
		Name:        "zi",
		Description: "jump to directory interactively",
		MaxArgs:     0,
		Generator:   ZoxideGenerator(),
	})
}

// zoxideCacheTTL bounds how stale the directory list may be. Generators run on
// every keystroke and `zoxide query -l` is a subprocess costing on the order of
// ten milliseconds, which is the entire latency budget for a suggestion redraw.
// The list only changes when the user changes directory, so a short window
// costs nothing in practice.
const zoxideCacheTTL = 2 * time.Second

var zoxideCache struct {
	sync.Mutex
	dirs    []string
	err     error
	fetched time.Time
}

// zoxideDirs returns zoxide's known directories, most frecent first.
func zoxideDirs() ([]string, error) {
	zoxideCache.Lock()
	defer zoxideCache.Unlock()

	if !zoxideCache.fetched.IsZero() && time.Since(zoxideCache.fetched) < zoxideCacheTTL {
		return zoxideCache.dirs, zoxideCache.err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	out, err := exec.CommandContext(ctx, "zoxide", "query", "-l").Output()
	zoxideCache.fetched = time.Now()
	zoxideCache.err = err
	zoxideCache.dirs = nil
	if err != nil {
		return nil, err
	}

	for line := range strings.SplitSeq(string(bytes.TrimSpace(out)), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			zoxideCache.dirs = append(zoxideCache.dirs, line)
		}
	}
	return zoxideCache.dirs, nil
}

func ZoxideGenerator() spec.GeneratorFunc {
	return func(tokens []string, prefix string, partial string) []spec.Suggestion {
		fullQuery := strings.Join(tokens[1:], " ")
		localSuggestions := spec.FileGenerator("/")(tokens, prefix, fullQuery)

		var zoxideSuggestions []spec.Suggestion
		dirs, err := zoxideDirs()
		if err == nil {
			home, _ := os.UserHomeDir()

			if fullQuery == "" {
				limit := min(len(dirs), 20)
				for i := range limit {
					path := dirs[i]
					display := strings.Replace(path, home, "~", 1)
					zoxideSuggestions = append(zoxideSuggestions, spec.Suggestion{
						Cmd:  path,
						Desc: display,
					})
				}
			} else if !strings.Contains(fullQuery, "/") {
				searcher := fuzzy.NewPlainSearcher(dirs)
				matches := searcher.SearchWithScores(fullQuery, &fuzzy.SearchOptions{Limit: 10})
				for _, m := range matches {
					path := m.Str
					display := strings.Replace(path, home, "~", 1)
					zoxideSuggestions = append(zoxideSuggestions, spec.Suggestion{
						Cmd:  path,
						Desc: display,
					})
				}
			}
		}

		var finalResults []spec.Suggestion
		seen := make(map[string]bool)

		finalResults = append(finalResults, localSuggestions...)

		for _, s := range zoxideSuggestions {
			if !seen[s.Cmd] {
				finalResults = append(finalResults, s)
				seen[s.Cmd] = true
			}
		}

		return finalResults
	}
}
