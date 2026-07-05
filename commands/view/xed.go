package view

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "xed",
		Description: "Xcode text editor invocation tool",
		Options: []spec.Option{
			{Name: "--launch", Description: "Launches Xcode, opening a new empty unsaved file"},
			{Name: "--create", Description: "Selects the given line in the last file opened"},
			{Name: "--background", Description: "Opens Xcode without activating it; the process that invoked xed remains in front"},
			{Name: "--help", Description: "Show help for xed"},
			{Name: "--version", Description: "Prints the version number of xed"},
		},
	})
}
