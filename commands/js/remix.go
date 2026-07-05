package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "remix",
		Description: "Represent the directory of the Remix application",
		Subcommands: []spec.Subcommand{
			{Name: "build", Description: "Create an optimized production build of your application"},
			{Name: "dev", Description: "Start the application in development mode"},
			{Name: "setup", Description: "Prepare node_modules/remix folder (after installation of packages)"},
			{Name: "routes", Description: "Generate the route config of the application"},
		},
		Options: []spec.Option{
			{Name: "--help", Description: "Output usage information"},
			{Name: "-v", Description: "Output the version number"},
			{Name: "--sourcemap", Description: "Enables production sourcemap"},
			{Name: "--json", Description: "Print the route config as JSON"},
		},
	})
}
