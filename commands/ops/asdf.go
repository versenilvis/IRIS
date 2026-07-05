package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "asdf",
		Description: "Plugin name",
		Options: []spec.Option{
			{Name: "--urls", Description: "Show git urls"},
			{Name: "--refs", Description: "Show git refs"},
			{Name: "--all", Description: "Update all plugins to latest commit on default branch"},
			{Name: "--head", Description: "Using HEAD commit"},
			{Name: "--version", Description: "Version for asdf"},
			{Name: "-h", Description: "Help for asdf"},
		},
	})
}
