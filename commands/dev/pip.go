package dev

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "pip",
		Description: "python packages",
		Subcommands: []core.Subcommand{
			{Name: "install", Description: "install package"},
			{Name: "uninstall", Description: "remove package"},
			{Name: "freeze", Description: "list packages"},
			{Name: "list", Description: "list installed"},
			{Name: "show", Description: "show package info"},
		},
	})
}
