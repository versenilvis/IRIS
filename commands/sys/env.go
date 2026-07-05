package sys

import (
	"os"
	"strings"

	"github.com/versenilvis/iris/spec"
)

func envVarGenerator(tokens []string, _ string, _ string) []spec.Suggestion {
	var results []spec.Suggestion
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
		results = append(results, spec.Suggestion{Cmd: name, Desc: val})
	}
	return results
}

func init() {
	spec.Register(&spec.Spec{
		Name:        "export",
		Description: "set environment variable",
		Generator:   envVarGenerator,
	})

	spec.Register(&spec.Spec{
		Name:        "env",
		Description: "print environment",
		Generator:   envVarGenerator,
		Options: []spec.Option{
			{Name: "-i", Description: "start with empty env"},
			{Name: "-u", Description: "unset variable"},
			{Name: "-0", Description: "null delimited"},
		},
	})

	spec.Register(&spec.Spec{
		Name:        "printenv",
		Description: "print environment variables",
		Generator:   envVarGenerator,
	})

	spec.Register(&spec.Spec{
		Name:        "unset",
		Description: "unset variable",
		Generator:   envVarGenerator,
	})
}
