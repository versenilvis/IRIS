package view

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "du",
		Description: "estimate file space usage",
		Generator:   core.FileGenerator(),
		Options: []core.Option{
			{Name: "-h", Description: "human readable"},
			{Name: "-s", Description: "summarize"},
			{Name: "-a", Description: "all files"},
		},
	})
}
