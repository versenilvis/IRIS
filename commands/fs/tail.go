package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "tail",
		Description: "output last lines of file",
		Generator:   core.FileGenerator(),
		Options: []core.Option{
			{Name: "-n", Description: "number of lines"},
			{Name: "-f", Description: "follow (live updates)"},
		},
	})
}
