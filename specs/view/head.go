package view

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "head",
		Description: "output first lines of file",
		Generator:   core.FileGenerator(),
		Options: []core.Option{
			{Name: "-n", Description: "number of lines"},
		},
	})
}
