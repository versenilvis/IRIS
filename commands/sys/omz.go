package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "omz",
		Description: "Oh My Zsh",
		Subcommands: []spec.Subcommand{
			{Name: "help", Description: "Print the help message"},
			{Name: "changelog", Description: "Print the changelog"},
			{Name: "plugin", Description: "Manage plugins"},
			{Name: "disable", Description: "Disable plugin(s)"},
			{Name: "enable", Description: "Enable plugin(s)"},
			{Name: "info", Description: "Get information of a plugin"},
			{Name: "list", Description: "List all available Oh My Zsh plugins"},
			{Name: "load", Description: "Load plugin(s)"},
			{Name: "pr", Description: "Manage Oh My Zsh Pull Requests"},
			{Name: "clean", Description: "Delete all PR branches (ohmyzsh/pull-*)"},
			{Name: "test", Description: "Fetch PR #NUMBER and rebase against master"},
			{Name: "reload", Description: "Reload the current zsh session"},
			{Name: "theme", Description: "Manage themes"},
			{Name: "set", Description: "Set a theme in your .zshrc file"},
			{Name: "use", Description: "Load a theme"},
			{Name: "update", Description: "Update Oh My Zsh"},
			{Name: "version", Description: "Show the version"},
		},
	})
}
