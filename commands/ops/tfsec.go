package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "tfsec",
		Description: "Terraform workspaces",
		Options: []spec.Option{
			{Name: "--concise-output", Description: "Reduce the amount of output and no statistics"},
			{Name: "--config-file", Description: "Config file to use during run"},
			{Name: "--custom-check-dir", Description: "Explicitly the custom checks dir location"},
			{Name: "--debug", Description: "Enable debug logging (same as verbose)"},
			{Name: "-G", Description: "Disable grouping of similar results"},
			{Name: "-e", Description: "Provide comma-separated list of rule IDs to exclude from run"},
			{Name: "--exclude-downloaded-modules", Description: "Remove results for downloaded modules in .terraform folder"},
			{Name: "--exclude-path", Description: "Filter results to return specific checks only (supports comma-delimited input)"},
			{Name: "--force-all-dirs", Description: "Don't search for tf files, include everything below provided directory"},
			{Name: "-f", Description: "Help for tfsec"},
			{Name: "--ignore-hcl-errors", Description: "Stop and report an error if an HCL parse error is encountered"},
			{Name: "--include-ignored", Description: "Include ignored checks in the result output"},
			{Name: "--include-passed", Description: "Include passed checks in the result output"},
			{Name: "--migrate-ignores", Description: "Migrate ignore codes to the new ID structure"},
			{Name: "-m", Description: "The minimum severity to report. One of CRITICAL, HIGH, MEDIUM, LOW"},
			{Name: "--no-color", Description: "Disable colored output (American style!)"},
			{Name: "--no-colour", Description: "Disable coloured output"},
			{Name: "--no-ignores", Description: "Do not apply any ignore rules - normally ignored checks will fail"},
			{Name: "--no-module-downloads", Description: "Do not download remote modules"},
			{Name: "-O", Description: "Print a JSON representation of the input supplied to rego policies"},
			{Name: "--rego-policy-dir", Description: "Directory to load rego policies from (recursively)"},
			{Name: "--run-statistics", Description: "View statistics table of current findings"},
			{Name: "--single-thread", Description: "Run checks using a single thread"},
			{Name: "-s", Description: "Runs checks but suppresses error code"},
			{Name: "--tfvars-file", Description: "Update to latest version"},
			{Name: "--verbose", Description: "Enable verbose logging (same as debug)"},
			{Name: "-v", Description: "Show version information and exit"},
		},
	})
}
