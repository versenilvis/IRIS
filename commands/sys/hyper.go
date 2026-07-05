package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "hyper",
		Description: "Hyper is an Electron-based terminal",
		Subcommands: []spec.Subcommand{
			{Name: "install", Description: "Install a plugin"},
			{Name: "docs", Description: "Open the npm page of a plugin"},
			{Name: "help", Description: "Display help"},
			{Name: "list", Description: "List installed plugins"},
			{Name: "list-remote", Description: "List plugins available on npm"},
			{Name: "search", Description: "Search for plugins on npm"},
			{Name: "uninstall", Description: "Uninstall plugin"},
			{Name: "plugin", Description: "Plugin to uninstall"},
			{Name: "version", Description: "Show version"},
		},
	})
}
