package fs

import (
	"strings"

	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "cd",
		Description: "change directory",
		MaxArgs:     0,
		Generator: func(tokens []string, prefix string, partial string) []spec.Suggestion {
			fullQuery := strings.Join(tokens[1:], " ")
			return spec.FileGenerator("/")(tokens, prefix, fullQuery)
		},
	})
}
