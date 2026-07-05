package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "remix",
		Description: "Represent the directory of the Remix application",
		Subcommands: []core.Subcommand{
			{Name: "build", Description: "Create an optimized production build of your application"},
			{Name: "dev", Description: "Start the application in development mode"},
			{Name: "setup", Description: "Prepare node_modules/remix folder (after installation of packages)"},
			{Name: "routes", Description: "Generate the route config of the application"},
		},
		Options: []core.Option{
			{Name: "--help", Description: "Output usage information"},
			{Name: "-v", Description: "Output the version number"},
			{Name: "--sourcemap", Description: "Enables production sourcemap"},
			{Name: "--json", Description: "Print the route config as JSON"},
		},
	})
}
