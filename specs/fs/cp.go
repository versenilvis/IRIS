package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "cp",
		Description: "copy files and directories",
		Generator:   core.FileGenerator(),
		Options: []core.Option{
			{Name: "-r", Description: "recursive"},
			{Name: "-a", Description: "archive mode"},
			{Name: "-v", Description: "verbose"},
			{Name: "-i", Description: "interactive"},
			{Name: "-u", Description: "update (copy only if newer)"},
			{Name: "-p", Description: "preserve attributes"},
		},
	})
}
