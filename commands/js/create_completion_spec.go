package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "create-completion-spec",
		Description: "Setup fig folder and create spec with the given name",
		Subcommands: []core.Subcommand{
			{Name: "help", Description: "Display help for command"},
		},
		Options: []core.Option{
			{Name: "--here", Description: "Set if the spec should be created in the current folder"},
			{Name: "-h", Description: "Display help for command"},
		},
	})
}
