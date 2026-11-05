package view

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "subl",
		Description: "Sublime Text",
		Options: []spec.Option{
			{Name: "--project", Description: "Load the given project"},
			{Name: "--command", Description: "Run the given command"},
			{Name: "-n", Description: "Open a new window"},
			{Name: "-a", Description: "Add folders to the current window"},
			{Name: "--launch-or-new-window", Description: "Only open a new window if the application is open"},
			{Name: "-w", Description: "Wait for the files to be closed before returning"},
			{Name: "-b", Description: "Don't activate the application"},
			{Name: "-s", Description: "Keep the application activated after closing the file"},
			{Name: "--safe-mode", Description: "Launch using a clean environment"},
			{Name: "-h", Description: "Show a help message and exit"},
			{Name: "-v", Description: "Show the version and exit"},
		},
	})
}
