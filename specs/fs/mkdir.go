package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "mkdir",
		Description: "make directories",
		Options: []core.Option{
			{Name: "-p", Description: "create parent dirs"},
			{Name: "-v", Description: "verbose"},
		},
	})
}
