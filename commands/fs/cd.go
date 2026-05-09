package fs

import (
	"strings"

	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "cd",
		Description: "change directory",
		MaxArgs:     0,
		Generator: func(tokens []string, prefix string, partial string) []core.Suggestion {
			fullQuery := strings.Join(tokens[1:], " ")
			return core.FileGenerator("/")(tokens, prefix, fullQuery)
		},
	})
}
