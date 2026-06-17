package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "less",
		Description: "view file contents (scrollable)",
		Generator:   core.FileGenerator(),
	})
}
