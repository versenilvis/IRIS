package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "dotnet",
		Description: "The dotnet cli",
		Options: []core.Option{
			{Name: "--version", Description: "Prints out a list of the installed .NET SDKs"},
			{Name: "-?", Description: "Prints out a list of available commands"},
			{Name: "-d", Description: "Enables diagnostic output"},
			{Name: "-v", Description: "Path containing probing policy and assemblies to probe"},
			{Name: "--additional-deps", Description: "Version of the .NET runtime to use to run the application"},
		},
	})
}
