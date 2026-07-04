package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "shortcuts",
		Description: "Run a shortcut",
		Subcommands: []core.Subcommand{
			{Name: "help", Description: "Show help information"},
		},
		Options: []core.Option{
			{Name: "-i", Description: "The input to provide to the shortcut"},
			{Name: "-o", Description: "Where to write the shortcut output, if applicable"},
			{Name: "--output-type", Description: "JavaScript Object Notation (JSON)"},
			{Name: "--folder-name", Description: "The name of the folder to list"},
			{Name: "--folders", Description: "List folders instead of shortcuts"},
			{Name: "--input", Description: "The shortcut file to sign"},
			{Name: "--output", Description: "Output path for the signed shortcut file"},
			{Name: "--mode", Description: "The signing mode. (default: people-who-know-me)"},
		},
	})
}
