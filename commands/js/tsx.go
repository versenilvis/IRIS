package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "tsx",
		Description: "Run TypeScript file using tsx",
		Subcommands: []spec.Subcommand{
			{Name: "watch", Description: "Run the script and watch for changes"},
		},
		Options: []spec.Option{
			{Name: "--help", Description: "Show help for tsx"},
			{Name: "--no-cache", Description: "Disable caching"},
			{Name: "--clear-screen", Description: "Disable clearing the screen on rerun"},
			{Name: "-v", Description: "Show version"},
			{Name: "--tsconfig", Description: "Custom tsconfig.json path"},
		},
	})
}
