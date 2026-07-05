package runner

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "vapor",
		Description: "Vapor Toolbox (Server-side Swift web framework)",
		Subcommands: []spec.Subcommand{
			{Name: "build", Description: "Builds an app in the console"},
			{Name: "new", Description: "Generates a new app"},
			{Name: "name", Description: "Name of project and folder"},
			{Name: "clean", Description: "Cleans temporary files"},
			{Name: "heroku", Description: "Commands for working with Heroku"},
			{Name: "init", Description: "Configures app for deployment to Heroku"},
			{Name: "push", Description: "Deploys app to Heroku"},
			{Name: "supervisor", Description: "Commands for working with Supervisord"},
			{Name: "restart", Description: "Restarts current project in Supervisor"},
			{Name: "update", Description: "Updates Supervisor entry for current project"},
			{Name: "xcode", Description: "Opens an app in Xcode"},
		},
		Options: []spec.Option{
			{Name: "--template", Description: "The URL of a Git repository to use as a template"},
			{Name: "--branch", Description: "Template repository branch to use"},
			{Name: "--output", Description: "The directory to place the new project in"},
			{Name: "--no-commit", Description: "Skips adding a first commit to the newly created repo"},
			{Name: "--update", Description: "Delete Package.resolved file if it exists"},
			{Name: "--global", Description: "Clean Xcode's global DerivedData cache"},
			{Name: "--swiftpm", Description: "Delete .swiftpm folder"},
			{Name: "--help", Description: "Show help for vapor"},
		},
	})
}
