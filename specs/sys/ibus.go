package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "ibus",
		Description: "Set or get engine",
		Subcommands: []core.Subcommand{
			{Name: "engine", Description: "Set or get engine"},
			{Name: "exit", Description: "Exit ibus-daemon"},
			{Name: "list-engine", Description: "Show available engines"},
			{Name: "watch", Description: "Not implemented"},
			{Name: "version", Description: "Show version"},
			{Name: "read-cache", Description: "Show the content of registry cache"},
			{Name: "write-cache", Description: "Create registry cache"},
			{Name: "address", Description: "Print the D-Bus address of ibus-daemon"},
			{Name: "read-config", Description: "Show the configuration values"},
			{Name: "reset-config", Description: "Reset the configuration values"},
			{Name: "emoji", Description: "Save emoji on dialog to clipboard"},
			{Name: "help", Description: "Show this information"},
		},
	})
}
