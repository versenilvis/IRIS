package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "pry",
		Description: "Interactive Ruby",
		Options: []spec.Option{
			{Name: "-e", Description: "A line of code to execute in context before the session starts"},
			{Name: "--no-pager", Description: "Disable pager for long output"},
			{Name: "--no-history", Description: "Disable history loading"},
			{Name: "--no-color", Description: "Disable syntax highlighting for session"},
			{Name: "-f", Description: "Suppress loading of pryrc"},
			{Name: "-s", Description: "Only load specified plugin (and no others)"},
			{Name: "-d", Description: "Disable a specific plugin"},
			{Name: "--no-plugins", Description: "Suppress loading of plugins"},
			{Name: "--plugins", Description: "List installed plugins"},
			{Name: "--simple-prompt", Description: "Enable simple prompt mode"},
			{Name: "--noprompt", Description: "No prompt mode"},
			{Name: "-r", Description: "`require` a Ruby script at startup"},
			{Name: "-I", Description: "Add a path to the $LOAD_PATH"},
			{Name: "--gem", Description: "Shorthand for -I./lib -rgemname"},
			{Name: "-v", Description: "Display the Pry version"},
			{Name: "-c", Description: "Display this help message"},
		},
	})
}
