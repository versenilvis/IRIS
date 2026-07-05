package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "nhost",
		Description: "Nhost",
		Subcommands: []core.Subcommand{
			{Name: "deploy", Description: "Deploy local migrations and metadata changes to Nhost production"},
			{Name: "dev", Description: "Start Nhost project for local development"},
			{Name: "down", Description: "Stop and remove local Nhost backend started by `nhost dev`"},
			{Name: "env", Description: "List environment variables"},
			{Name: "env:pull", Description: "Sync remote environment variables to your local environment"},
			{Name: "help", Description: "Display help for nhost"},
			{Name: "init", Description: "Initialize current working directory as a Nhost project"},
			{Name: "link", Description: "Link Nhost Project"},
			{Name: "login", Description: "Login to your Nhost account"},
			{Name: "logout", Description: "Logout from your Nhost account"},
		},
		Options: []core.Option{
			{Name: "--email", Description: "Email address"},
		},
	})
}
