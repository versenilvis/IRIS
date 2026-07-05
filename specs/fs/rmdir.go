package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "rmdir",
		Description: "Remove directories",
		Options: []core.Option{
			{Name: "-p", Description: "Remove each directory of path"},
		},
	})
}
