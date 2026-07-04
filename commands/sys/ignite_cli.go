package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "ignite-cli",
		Description: "Output usage information",
		Subcommands: []core.Subcommand{
			{Name: "new", Description: "Create a new React Native app"},
			{Name: "g", Description: "Generate components and other app features"},
			{Name: "update", Description: "Update installed generators"},
			{Name: "doctor", Description: "Check your environment & display versions of installed dependencies"},
		},
		Options: []core.Option{
			{Name: "-h", Description: "Output usage information"},
			{Name: "-v", Description: "Output the version number"},
			{Name: "--expo", Description: "Use Expo"},
			{Name: "--bundle", Description: "Set the bundle ID of the app"},
			{Name: "--list", Description: "List installed generators"},
			{Name: "--update", Description: "Update installed generators"},
			{Name: "--all", Description: "Update all installed generators"},
		},
	})
}
