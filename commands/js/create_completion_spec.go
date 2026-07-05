package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "create-completion-spec",
		Description: "Setup fig folder and create spec with the given name",
		Subcommands: []spec.Subcommand{
			{Name: "help", Description: "Display help for command"},
		},
		Options: []spec.Option{
			{Name: "--here", Description: "Set if the spec should be created in the current folder"},
			{Name: "-h", Description: "Display help for command"},
		},
	})
}
