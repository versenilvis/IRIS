package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "command",
		Description: "Run an external command",
		Options: []core.Option{
			{Name: "-v", Description: "Print the location of the command"},
		},
	})
}
