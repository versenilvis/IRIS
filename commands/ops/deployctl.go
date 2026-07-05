package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "deployctl",
		Description: "Command line tool for Deno Deploy",
		Subcommands: []spec.Subcommand{
			{Name: "deploy", Description: "Deploy a script with static files to Deno Deploy"},
			{Name: "upgrade", Description: "Upgrade deployctl"},
		},
		Options: []spec.Option{
			{Name: "--exclude", Description: "Exclude files that match this pattern"},
			{Name: "--include", Description: "Only upload files that match this pattern"},
			{Name: "--no-static", Description: "Don't include the files in the CWD as static files"},
			{Name: "--prod", Description: "Create a production deployment (default is preview deployment)"},
			{Name: "-p", Description: "The project to deploy to"},
			{Name: "--token", Description: "The API token to use"},
			{Name: "--help", Description: "Show help"},
			{Name: "-V", Description: "Show the version"},
		},
	})
}
