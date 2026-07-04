package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "coda",
		Description: "Execute ${formulaName}",
		Subcommands: []core.Subcommand{
			{Name: "init", Description: "Initialize an empty project with the recommended settings and dependencies"},
			{Name: "execute", Description: "Execute the formula and print the output to the terminal"},
			{Name: "formula", Description: "Formula name to execute"},
			{Name: "params", Description: "Arguments to pass to the formula"},
			{Name: "url", Description: "The URL to sync from"},
			{Name: "apiToken", Description: "The API token to register"},
			{Name: "name", Description: "The desired Pack name"},
			{Name: "description", Description: "The Pack description"},
			{Name: "upload", Description: "Use this command to upload a new version of your Pack based on your latest code"},
			{Name: "release", Description: "Release a Pack version and make it live for your users"},
			{Name: "whoami", Description: "Looks up information about the API token that is registered in this environment"},
			{Name: "build", Description: "Generate a bundle for your Pack"},
			{Name: "validate", Description: "Validate your Pack definition"},
			{Name: "option", Description: "Currently the only supported option is 'timerStrategy'"},
			{Name: "value", Description: "Value to set for the option"},
		},
		Options: []core.Option{
			{Name: "--dynamicUrl", Description: "The URL to sync from"},
			{Name: "--name", Description: "The desired Pack name"},
			{Name: "--description", Description: "The Pack description"},
			{Name: "--version", Description: "Show version number"},
			{Name: "--help", Description: "Show help"},
		},
	})
}
