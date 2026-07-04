package view

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "xed",
		Description: "Xcode text editor invocation tool",
		Options: []core.Option{
			{Name: "--launch", Description: "Launches Xcode, opening a new empty unsaved file"},
			{Name: "--create", Description: "Selects the given line in the last file opened"},
			{Name: "--background", Description: "Opens Xcode without activating it; the process that invoked xed remains in front"},
			{Name: "--help", Description: "Show help for xed"},
			{Name: "--version", Description: "Prints the version number of xed"},
		},
	})
}
