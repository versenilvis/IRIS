package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "publish",
		Description: "Set up a new website in the current folder",
		Subcommands: []spec.Subcommand{
			{Name: "new", Description: "Set up a new website in the current folder"},
			{Name: "deploy", Description: "Generate and deploy the website in the current folder"},
			{Name: "generate", Description: "Generate the website in the current folder"},
		},
		Options: []spec.Option{
			{Name: "-p", Description: "Customize the port"},
			{Name: "--help", Description: "Show help for publish"},
		},
	})
}
