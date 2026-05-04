// please note that zoxide also shows external suggestions
// at the end of the list on command mode
// they are the old directories that you have visited
// this is a feature, not a bug, and I want to keep it
package fs

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/versenilvis/fuzzyvn"
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "z",
		Description: "jump to directory",
		MaxArgs:     1,
		Generator:   zoxideGenerator(),
	})
	core.Register(&core.Spec{
		Name:        "zi",
		Description: "jump to directory interactively",
		MaxArgs:     1,
		Generator:   zoxideGenerator(),
	})
}

func zoxideGenerator() core.GeneratorFunc {
	return func(tokens []string, prefix string, partial string) []core.Suggestion {

		localSuggestions := core.FileGenerator("/")(tokens, prefix, partial)

		var zoxideSuggestions []core.Suggestion
		cmd := exec.Command("zoxide", "query", "-l")
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

			if partial == "" {
				limit := 20
				if len(dirs) < limit {
					limit = len(dirs)
				}
				for i := 0; i < limit; i++ {
					path := dirs[i]
					display := strings.Replace(path, home, "~", 1)
					zoxideSuggestions = append(zoxideSuggestions, core.Suggestion{
						Cmd:  path,
						Desc: display,
					})
				}
			} else if !strings.Contains(partial, "/") {

				searcher := fuzzyvn.NewPlainSearcher(dirs)
				matches := searcher.SearchWithScores(partial, &fuzzyvn.SearchOptions{Limit: 10})
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

		for _, s := range localSuggestions {
			finalResults = append(finalResults, s)

		}

		for _, s := range zoxideSuggestions {
			if !seen[s.Cmd] {
				finalResults = append(finalResults, s)
				seen[s.Cmd] = true
			}
		}

		return finalResults
	}
}
