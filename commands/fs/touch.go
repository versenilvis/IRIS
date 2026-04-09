package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "touch",
		Description: "create or update file timestamp",
		Generator:   core.FileGenerator(),
	})
}
