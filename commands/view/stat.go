package view

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "stat",
		Description: "display file status",
		Generator:   core.FileGenerator(),
	})
}
