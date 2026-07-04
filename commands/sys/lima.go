package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "lima",
		Description: "Lima is an alias for",
		Options: []core.Option{
			{Name: "-h", Description: "Help for lima"},
		},
	})
}
