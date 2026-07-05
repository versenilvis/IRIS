package view

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "vi",
		Description: "Print help message for vi and exit",
		Options: []core.Option{
			{Name: "-h", Description: "Print help message for vi and exit"},
		},
	})
}
