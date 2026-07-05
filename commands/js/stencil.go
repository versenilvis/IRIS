package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "stencil",
		Description: "CLI to build Stencil projects and generate components",
		Subcommands: []spec.Subcommand{
			{Name: "build", Description: "Build components for development or production"},
			{Name: "workers", Description: "Number of workers"},
			{Name: "test", Description: "Run unit and end-to-end tests"},
			{Name: "generate", Description: "Bootstrap components"},
			{Name: "telemetry", Description: "Opt in or out of telemetry"},
			{Name: "off", Description: "Disable sharing anonymous usage data"},
			{Name: "on", Description: "Enable sharing anonymous usage data"},
		},
		Options: []spec.Option{
			{Name: "--ci", Description: "Set stencil config file"},
			{Name: "--debug", Description: "Set the log level to debug"},
			{Name: "--dev", Description: "Development build"},
			{Name: "--docs-readme", Description: "Generate component readme.md docs"},
			{Name: "--es5", Description: "Creates an ES5 compatible build"},
			{Name: "--log", Description: "Write stencil-build.log file"},
			{Name: "--prerender", Description: "Prerender the application"},
			{Name: "--prod", Description: "Runs a production build"},
			{Name: "--max-workers", Description: "Max number of workers the compiler should use"},
			{Name: "--next", Description: "Opt-in to test the 'next' Stencil compiler features"},
			{Name: "--no-cache", Description: "Disables using the cache"},
			{Name: "--no-open", Description: "Do not automatically open the browser window"},
			{Name: "--port", Description: "Port for the Integrated Dev Server"},
			{Name: "--serve", Description: "Start the dev-server"},
			{Name: "--stats", Description: "Write stencil-stats.json file"},
			{Name: "--verbose", Description: "Logs additional information about each step of the build"},
			{Name: "--watch", Description: "Rebuild when files update"},
			{Name: "--spec", Description: "Run unit tests with Jest"},
			{Name: "--e2e", Description: "Run e2e tests with Puppeteer"},
			{Name: "--no-build", Description: "Skips build process before running the tests"},
			{Name: "--help", Description: "Display the help output explaining the different flags"},
			{Name: "--version", Description: "Prints the current Stencil version"},
		},
	})
}
