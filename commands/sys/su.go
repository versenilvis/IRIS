package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "su",
		Description: "(no letter) The same as -l",
		Options: []core.Option{
			{Name: "-f", Description: "(no letter) The same as -l"},
		},
	})
}
