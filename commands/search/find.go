package search

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "find",
		Description: "search for files",
		Generator:   core.FileGenerator("/"),
		Options: []core.Option{
			{Name: "-name", Description: "match by name"},
			{Name: "-type", Description: "match by type (f/d)"},
			{Name: "-size", Description: "match by size"},
			{Name: "-mtime", Description: "match by modified time"},
		},
	})
}
