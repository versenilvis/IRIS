package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "vr",
		Description: "The npm-style script runner for Deno",
		Subcommands: []spec.Subcommand{
			{Name: "run", Description: "Run a script"},
			{Name: "export", Description: "Export one or more scripts as standalone executable files"},
			{Name: "upgrade", Description: "Upgrade Velociraptor to the latest version or to a specific one"},
		},
		Options: []spec.Option{
			{Name: "-o", Description: "The folder where the scripts will be exported"},
			{Name: "--help", Description: "Show help for Velociraptor"},
			{Name: "-V", Description: "Show the version number for Velociraptor"},
		},
	})
}
