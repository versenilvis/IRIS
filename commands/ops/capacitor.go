package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "capacitor",
		Description: "Add a native platform project to your app",
		Subcommands: []spec.Subcommand{
			{Name: "add", Description: "Add a native platform project to your app"},
			{Name: "ls", Description: "List all installed Cordova and Capacitor plugins"},
			{Name: "sync", Description: "This command runs copy and then update"},
			{Name: "update", Description: "Updates the native plugins and dependencies referenced in package.json"},
		},
		Options: []spec.Option{
			{Name: "--list", Description: "Print a list of target devices available to the given platform"},
			{Name: "--target", Description: "Run on a specific target device"},
			{Name: "--deployment", Description: "Podfile.lock won't be deleted and pod install will use --deployment option"},
			{Name: "--inline", Description: "Updates the native plugins and dependencies referenced in package.json"},
			{Name: "--help", Description: "Output usage information. Can be used with individual commands too"},
			{Name: "--version", Description: "Output the version number"},
		},
	})
}
