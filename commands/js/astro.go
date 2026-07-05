package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "astro",
		Description: "Add an integration",
		Subcommands: []spec.Subcommand{
			{Name: "add", Description: "Add an integration"},
			{Name: "build", Description: "Build your project and write it to disk"},
			{Name: "check", Description: "Check your project for errors"},
			{Name: "dev", Description: "Starts the development server"},
			{Name: "docs", Description: "Open documentation in your web browser"},
			{Name: "preview", Description: "Preview your build locally"},
			{Name: "sync", Description: "Generate content collection types"},
			{Name: "telemetry", Description: "Configure telemetry settings"},
			{Name: "enable", Description: "Enable anonymous data collection"},
			{Name: "disable", Description: "Disable anonymous data collection"},
			{Name: "reset", Description: "Reset anonymous data collection settings"},
		},
		Options: []spec.Option{
			{Name: "--yes", Description: "Accept all prompts"},
			{Name: "--drafts", Description: "Include Markdown draft pages in the build"},
			{Name: "--watch", Description: "Watch Astro files for changes and re-run checks"},
			{Name: "--port", Description: "Specify a port to listen on"},
			{Name: "--host", Description: "Listen on all addresses, including LAN and public addresses"},
			{Name: "--open", Description: "Automatically open the app in the browser on server start"},
			{Name: "--config", Description: "Specify your config file"},
			{Name: "--root", Description: "Specify your project root folder"},
			{Name: "--site", Description: "Specify your project site"},
			{Name: "--base", Description: "Specify your project base"},
			{Name: "--verbose", Description: "Enable verbose logging"},
			{Name: "--silent", Description: "Disable all logging"},
			{Name: "--version", Description: "Show the version number and exit"},
			{Name: "--help", Description: "Show help for astro"},
		},
	})
}
