package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "chown",
		Description: "change file owner",
		Generator:   core.FileGenerator(),
		Options: []core.Option{
			{Name: "-R", Description: "recursive"},
		},
	})
}
