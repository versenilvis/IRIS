package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "echo",
		Description: "Environment Variable",
		Options: []core.Option{
			{Name: "-n", Description: "Do not print the trailing newline character"},
			{Name: "-e", Description: "Interpret escape sequences"},
			{Name: "-E", Description: "Disable escape sequences"},
		},
	})
}
