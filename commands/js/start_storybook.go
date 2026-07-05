package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "start-storybook",
		Description: "Display usage information",
		Options: []spec.Option{
			{Name: "--help", Description: "Display usage information"},
			{Name: "-V", Description: "Display the version number"},
			{Name: "-s", Description: "Directory where to load static files from, comma-separated list"},
			{Name: "-c", Description: "Directory where to load Storybook configurations from"},
			{Name: "--https", Description: "Provide an SSL certificate. (Required with --https)"},
			{Name: "--ssl-key", Description: "Provide an SSL key. (Required with --https)"},
			{Name: "--smoke-test", Description: "Exit after successful start"},
			{Name: "--ci", Description: "CI mode (skip interactive prompts, don't open browser)"},
			{Name: "--quiet", Description: "Suppress verbose build output"},
			{Name: "--no-dll", Description: "Do not use dll reference (no-op)"},
			{Name: "--debug-webpack", Description: "Display final webpack configurations for debugging purposes"},
			{Name: "--webpack-stats-json", Description: "Write Webpack Stats JSON to disk"},
			{Name: "--docs", Description: "Starts Storybook in documentation mode"},
			{Name: "--no-manager-cache", Description: "Storybook start CLI tools"},
			{Name: "-p", Description: "Port to run Storybook"},
			{Name: "-h", Description: "Host to run Storybook"},
		},
	})
}
