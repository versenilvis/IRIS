package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "vr",
		Description: "The npm-style script runner for Deno",
		Subcommands: []core.Subcommand{
			{Name: "run", Description: "Run a script"},
			{Name: "export", Description: "Export one or more scripts as standalone executable files"},
			{Name: "upgrade", Description: "Upgrade Velociraptor to the latest version or to a specific one"},
		},
		Options: []core.Option{
			{Name: "-o", Description: "The folder where the scripts will be exported"},
			{Name: "--help", Description: "Show help for Velociraptor"},
			{Name: "-V", Description: "Show the version number for Velociraptor"},
		},
	})
}
