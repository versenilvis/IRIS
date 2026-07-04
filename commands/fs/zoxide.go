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

	"github.com/versenilvis/fuzzy"
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "z",
		Description: "jump to directory",
		MaxArgs:     0,
		Generator:   ZoxideGenerator(),
	})
	core.Register(&core.Spec{
		Name:        "zi",
		Description: "jump to directory interactively",
		MaxArgs:     0,
		Generator:   ZoxideGenerator(),
	})
}

func ZoxideGenerator() core.GeneratorFunc {
	return func(tokens []string, prefix string, partial string) []core.Suggestion {
		fullQuery := strings.Join(tokens[1:], " ")
		localSuggestions := core.FileGenerator("/")(tokens, prefix, fullQuery)

		var zoxideSuggestions []core.Suggestion
		cmd := exec.CommandContext(context.Background(), "zoxide", "query", "-l")
		out, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(bytes.TrimSpace(out)), "\n")
			var dirs []string
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					dirs = append(dirs, line)
				}
			}

			home, _ := os.UserHomeDir()

			if fullQuery == "" {
				limit := min(len(dirs), 20)
				for i := 0; i < limit; i++ {
					path := dirs[i]
					display := strings.Replace(path, home, "~", 1)
					zoxideSuggestions = append(zoxideSuggestions, core.Suggestion{
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
					zoxideSuggestions = append(zoxideSuggestions, core.Suggestion{
						Cmd:  path,
						Desc: display,
					})
				}
			}
		}

		var finalResults []core.Suggestion
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
