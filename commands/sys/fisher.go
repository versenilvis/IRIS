package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "fisher",
		Description: "[Prompt] - 🌊 The ultimate Fish prompt",
		Subcommands: []core.Subcommand{
			{Name: "install", Description: "Install plugin"},
			{Name: "remove", Description: "Remove plugins"},
			{Name: "installed plugins", Description: "The plugin you want to remove"},
			{Name: "update", Description: "Update plugins"},
			{Name: "list", Description: "List plugins"},
			{Name: "RegEx", Description: "Search in list with regular expression"},
		},
		Options: []core.Option{
			{Name: "--help", Description: "Show help for fisher"},
			{Name: "--version", Description: "Show fisher version"},
		},
	})
}
