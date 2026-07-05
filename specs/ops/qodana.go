package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "qodana",
		Description: "Run Qodana as fast as possible, with minimum effort required",
		Subcommands: []core.Subcommand{
			{Name: "init", Description: "Configure project for Qodana"},
			{Name: "scan", Description: "Scan a project with Qodana"},
			{Name: "show", Description: "Show Qodana report"},
			{Name: "view", Description: "View SARIF files in CLI"},
		},
		Options: []core.Option{
			{Name: "-h", Description: "Show help page for command"},
			{Name: "-v", Description: "Version for Qodana"},
			{Name: "-i", Description: "Scan a project with Qodana"},
			{Name: "-a", Description: "Unique report identifier (GUID) to be used by Qodana Cloud"},
			{Name: "-b", Description: "Override cache directory (default <userCacheDir>/JetBrains/<linter>/cache)"},
			{Name: "-c", Description: "Override the docker image to be used for the analysis"},
			{Name: "--clear-cache", Description: "Clear the local Qodana cache before running the analysis"},
			{Name: "--disable-sanity", Description: "Skip running the inspections configured by the sanity profile"},
			{Name: "-e", Description: "Override linter to use"},
			{Name: "--port", Description: "Port to serve the report on (default 8080)"},
			{Name: "--print-problems", Description: "Print all found problems by Qodana in the CLI output"},
			{Name: "-n", Description: "Profile name defined in the project"},
			{Name: "-p", Description: "Path to the profile file"},
			{Name: "--script", Description: "Serve HTML report on port"},
			{Name: "--skip-pull", Description: "Skip pulling the latest Qodana container"},
			{Name: "-d", Description: "Show Qodana report"},
		},
	})
}
