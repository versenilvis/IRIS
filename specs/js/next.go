package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "next",
		Description: "A port number on which to start the application",
		Subcommands: []core.Subcommand{
			{Name: "build", Description: "Create an optimized production build of your application"},
			{Name: "dev", Description: "Start the application in development mode"},
			{Name: "start", Description: "Start the application in production mode"},
			{Name: "export", Description: "Exports the application for production deployment"},
			{Name: "telemetry", Description: "Allows you to control Next.js' telemetry collection"},
			{Name: "status", Description: "Turn Next.js' telemetry collection on or off"},
			{Name: "enable", Description: "Enable Next.js' telemetry collection"},
			{Name: "disable", Description: "Disable Next.js' telemetry collection"},
		},
		Options: []core.Option{
			{Name: "-p", Description: "A port number on which to start the application"},
			{Name: "-H", Description: "Hostname on which to start the application"},
			{Name: "-h", Description: "Output usage information"},
			{Name: "-v", Description: "Output the version number"},
			{Name: "--profile", Description: "Enable production profiling"},
			{Name: "--debug", Description: "Enable more verbose build output"},
			{Name: "-s", Description: "Do not print any messages to console"},
		},
	})
}
