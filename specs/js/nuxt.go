package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "nuxt",
		Description: "Launch the development server",
		Subcommands: []core.Subcommand{
			{Name: "dev", Description: "Launch the development server"},
			{Name: "build", Description: "Build and optimize your application with webpack for production"},
			{Name: "webpack", Description: "Inspect the webpack config"},
		},
		Options: []core.Option{
			{Name: "--name", Description: "Bundle name to inspect. (client, server, modern)"},
			{Name: "--dev", Description: "Inspect webpack config for dev mode"},
			{Name: "--depth", Description: "Inspection depth. Defaults to 2 to prevent verbose output"},
		},
	})
}
