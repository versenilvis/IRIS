package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "ls",
		Description: "list directory contents",
		Generator:   core.FileGenerator(),
		Options: []core.Option{
			{Name: "-a", Description: "all files"},
			{Name: "-l", Description: "long format"},
			{Name: "-h", Description: "human readable"},
			{Name: "-R", Description: "recursive"},
			{Name: "-t", Description: "sort by time"},
		},
	})
}
