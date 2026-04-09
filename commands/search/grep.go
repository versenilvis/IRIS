package search

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "grep",
		Description: "search text in files",
		Generator:   core.FileGenerator(),
		Options: []core.Option{
			{Name: "-r", Description: "recursive"},
			{Name: "-i", Description: "ignore case"},
			{Name: "-n", Description: "show line numbers"},
			{Name: "-v", Description: "invert match"},
			{Name: "-l", Description: "show filenames only"},
		},
	})
}
