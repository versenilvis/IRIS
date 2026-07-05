package sys

import (
	"os"
	"strings"

	"github.com/versenilvis/iris/commands/core"
)

func envVarGenerator(tokens []string, _ string, _ string) []core.Suggestion {
	var results []core.Suggestion
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		name := parts[0]
		val := ""
		if len(parts) == 2 {
			val = parts[1]
			// truncate long values
			if len(val) > 60 {
				val = val[:57] + "..."
			}
		}
		results = append(results, core.Suggestion{Cmd: name, Desc: val})
	}
	return results
}

func init() {
	core.Register(&core.Spec{
		Name:        "export",
		Description: "set environment variable",
		Generator:   envVarGenerator,
	})

	core.Register(&core.Spec{
		Name:        "env",
		Description: "print environment",
		Generator:   envVarGenerator,
		Options: []core.Option{
			{Name: "-i", Description: "start with empty env"},
			{Name: "-u", Description: "unset variable"},
			{Name: "-0", Description: "null delimited"},
		},
	})

	core.Register(&core.Spec{
		Name:        "printenv",
		Description: "print environment variables",
		Generator:   envVarGenerator,
	})

	core.Register(&core.Spec{
		Name:        "unset",
		Description: "unset variable",
		Generator:   envVarGenerator,
	})
}
