package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "cat",
		Description: "concatenate and print",
		Generator:   core.FileGenerator(),
		Options: []core.Option{
			{Name: "-n", Description: "number lines"},
			{Name: "-b", Description: "number non-blank"},
		},
	})
}
