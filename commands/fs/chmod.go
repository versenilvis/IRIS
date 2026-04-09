package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "chmod",
		Description: "change file permissions",
		Generator:   core.FileGenerator(),
		Options: []core.Option{
			{Name: "-R", Description: "recursive"},
		},
	})
}
