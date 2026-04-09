package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "rm",
		Description: "remove files or directories",
		Generator:   core.FileGenerator(),
		Options: []core.Option{
			{Name: "-r", Description: "recursive"},
			{Name: "-f", Description: "force"},
			{Name: "-i", Description: "interactive"},
			{Name: "-v", Description: "verbose"},
			{Name: "-d", Description: "remove empty directories"},
		},
	})
}
