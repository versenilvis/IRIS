package runner

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "just",
		Description: "command runner",
		MaxArgs:     1,
		Generator: func(tokens []string, prefix string, partial string) []core.Suggestion {
			file, err := os.Open("justfile")
			if err != nil {
				// try uppercase
				file, err = os.Open("Justfile")
				if err != nil {
					return nil
				}
			}
			defer file.Close()

			var suggestions []core.Suggestion
			seen := make(map[string]bool)
			scanner := bufio.NewScanner(file)
			// matching recipes like `build:`
			recipeRegex := regexp.MustCompile(`^([a-zA-Z0-9_-]+):`)

			lastComment := ""
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				
				// Handle comments
				if strings.HasPrefix(line, "#") {
					lastComment = strings.TrimSpace(strings.TrimPrefix(line, "#"))
					continue
				}

				// Skip blank lines or group attributes
				if line == "" || strings.HasPrefix(line, "[") {
					continue
				}

				if matches := recipeRegex.FindStringSubmatch(line); len(matches) > 1 {
					recipe := matches[1]
					if seen[recipe] {
						lastComment = ""
						continue
					}
					seen[recipe] = true

					cmd := recipe
					
					desc := "just recipe"
					if lastComment != "" {
						desc = lastComment
					}

					suggestions = append(suggestions, core.Suggestion{
						Cmd:  cmd,
						Desc: desc,
					})
					lastComment = ""
				} else {
					// If it's something else (like a command line), reset comment
					// unless it's just indented (which we already skip or handle)
					if !strings.HasPrefix(scanner.Text(), " ") && !strings.HasPrefix(scanner.Text(), "\t") {
						lastComment = ""
					}
				}
			}
			return suggestions
		},
	})
}
