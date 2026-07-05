package runner

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "just",
		Description: "command runner",
		MaxArgs:     1,
		Generator: func(tokens []string, prefix string, partial string) []spec.Suggestion {
			file, err := os.Open("justfile")
			if err != nil {
				// try uppercase
				file, err = os.Open("Justfile")
				if err != nil {
					return nil
				}
			}
			defer func() { _ = file.Close() }()

			var suggestions []spec.Suggestion
			seen := make(map[string]bool)
			scanner := bufio.NewScanner(file)
			// matching recipes like `build:`
			recipeRegex := regexp.MustCompile(`^([a-zA-Z0-9_-]+):`)

			lastComment := ""
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())

				// Handle comments
				if after, ok := strings.CutPrefix(line, "#"); ok {
					lastComment = strings.TrimSpace(after)
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

					suggestions = append(suggestions, spec.Suggestion{
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
