package fs

import (
	"strings"

	"github.com/versenilvis/iris/internal/config"
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "cd",
		Description: "change directory",
		MaxArgs:     0,
		Generator: func(tokens []string, prefix string, partial string) []spec.Suggestion {
			// The file generator only reads the current directory, so `cd proj`
			// comes back empty unless its parent is where you already are --
			// which is not how cd is typically used. zoxide.extend-cd folds in
			// the frecency database that `z` already queries, making any
			// directory zoxide has seen reachable by name from anywhere.
			if config.Get().Zoxide.ExtendCd {
				return ZoxideGenerator()(tokens, prefix, partial)
			}
			fullQuery := strings.Join(tokens[1:], " ")
			return spec.FileGenerator("/")(tokens, prefix, fullQuery)
		},
	})
}
