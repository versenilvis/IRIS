package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "elm",
		Description: "Fig spec for the Elm language cli",
		Subcommands: []core.Subcommand{
			{Name: "reactor", Description: "Start an Elm development server"},
			{Name: "port", Description: "The port number"},
			{Name: "package", Description: "The name of the package to install"},
			{Name: "make", Description: "Build your Elm code"},
			{Name: "output file", Description: "Name and location of output"},
			{Name: "output json", Description: "Name and location of output"},
			{Name: "bump", Description: "Bump the version of your package"},
			{Name: "diff", Description: "See what changed between versions of a package"},
		},
		Options: []core.Option{
			{Name: "--help", Description: "Show help for elm init"},
			{Name: "--no-colors", Description: "Turn off colors in the repl"},
			{Name: "--interpreter", Description: "Path to an alternate JS interpreter, such as Node or Deno"},
			{Name: "--port", Description: "The port to access the development server on"},
			{Name: "--debug", Description: "Compile in debug mode"},
			{Name: "--optimize", Description: "Compile in production mode"},
			{Name: "--output", Description: "Where to output the compiled code"},
			{Name: "--docs", Description: "Generate a JSON file of documentation"},
		},
	})
}
