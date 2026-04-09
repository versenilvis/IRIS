package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "cd",
		Description: "change directory",
		Generator:   core.FileGenerator("/"), // Directory only
	})
}
