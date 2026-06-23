package view

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "file",
		Description: "determine file type",
		Generator:   core.FileGenerator(),
	})
}
