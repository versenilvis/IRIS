package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "ng",
		Description: "Project name",
		Subcommands: []core.Subcommand{
			{Name: "new", Description: "Create a new Angular app"},
			{Name: "generate", Description: "Generate new files"},
			{Name: "application", Description: "Generates a new application"},
			{Name: "name", Description: "Name of the new app"},
			{Name: "component", Description: "Generate a new component"},
			{Name: "library", Description: "Generates a new library"},
			{Name: "class", Description: "Generates a class"},
			{Name: "version", Description: "View your Angular CLI version (update for Angular 14+)"},
		},
		Options: []core.Option{
			{Name: "--project", Description: "Project name"},
			{Name: "--create-application", Description: "Create a default application?"},
			{Name: "--style", Description: "Generate a new component"},
			{Name: "--change-detection", Description: "The change detection strategy to use"},
			{Name: "--display-block", Description: "Add :host block to styles"},
			{Name: "--flat", Description: "Create at the top level"},
			{Name: "--version", Description: "View your Angular CLI version"},
		},
	})
}
