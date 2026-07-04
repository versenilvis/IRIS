package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "basename",
		Description: "Return filename portion of pathname",
		Options: []core.Option{
			{Name: "-a", Description: "Treat every argument as a string"},
			{Name: "-s", Description: "Suffix to remove from string"},
		},
	})
}
