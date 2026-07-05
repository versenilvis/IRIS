package runner

import (
	"bufio"
	"os"
	"regexp"

	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "make",
		Description: "build automation",
		Generator: func(tokens []string, prefix string, partial string) []core.Suggestion {
			// Read the Makefile
			file, err := os.Open("Makefile")
			if err != nil {
				// No Makefile, don't pretend we have one
				return nil
			}
			defer func() { _ = file.Close() }()

			var suggestions []core.Suggestion
			seen := make(map[string]bool)
			scanner := bufio.NewScanner(file)
			// Matches targets like `build:`, `  run :`, etc.
			// Ignores lines starting with whitespace followed by command (tab)
			targetRegex := regexp.MustCompile(`^[ \t]*([a-zA-Z0-9_-]+)[ \t]*:`)

			for scanner.Scan() {
				line := scanner.Text()
				if matches := targetRegex.FindStringSubmatch(line); len(matches) > 1 {
					target := matches[1]
					// Ignore standard non-target keywords
					if target == "PHONY" || seen[target] {
						continue
					}
					seen[target] = true
					
					cmd := target
					if prefix != "" {
						cmd = prefix + " " + target
					}
					suggestions = append(suggestions, core.Suggestion{
						Cmd:  cmd,
						Desc: "make target",
					})
				}
			}
			return suggestions
		},
	})
}
