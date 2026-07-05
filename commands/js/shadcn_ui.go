package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "shadcn-ui",
		Description: "Shadcn UI CLI",
		Subcommands: []spec.Subcommand{
			{Name: "add", Description: "Add a component to your project"},
			{Name: "components", Description: "The components to add"},
			{Name: "diff", Description: "Check for updates against the registry"},
			{Name: "component", Description: "The component name"},
			{Name: "init", Description: "Initialize your project and install dependencies"},
		},
		Options: []spec.Option{
			{Name: "-y", Description: "Skip confirmation prompt"},
			{Name: "-o", Description: "Overwrite existing files"},
			{Name: "-c", Description: "The working directory. defaults to the current directory"},
			{Name: "-p", Description: "The path to add the component to"},
		},
	})
}
