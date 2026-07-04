package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "goto",
		Description: "Goto",
		Options: []core.Option{
			{Name: "--help", Description: "Show help for goto"},
			{Name: "--register", Description: "Registers an alias"},
			{Name: "--unregister", Description: "Unregister an alias"},
			{Name: "--push", Description: "Pushes the current directory onto the stack, then performs goto"},
			{Name: "--pop", Description: "Pops the top directory from the stack, then changes to that directory"},
			{Name: "--list", Description: "Pops the top directory from the stack, then changes to that directory"},
			{Name: "--expand", Description: "Expands an alias"},
			{Name: "--cleanup", Description: "Cleans up non existent directory aliases"},
			{Name: "--version", Description: "Displays the version of the goto script"},
		},
	})
}
