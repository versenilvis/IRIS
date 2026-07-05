package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "elm-review",
		Description: "Prints a single JSON object",
		Subcommands: []spec.Subcommand{
			{Name: "init", Description: "Initialize Elm Review in this directory"},
			{Name: "suppress", Description: "Generate suppression files for rules that report many errors"},
			{Name: "new-package", Description: "Creates an new project aimed to contain rules and to be published later"},
			{Name: "new-rule", Description: "Adds a new rule to your review configuration or review package"},
		},
		Options: []spec.Option{
			{Name: "--help", Description: "Show help for elm-review init"},
			{Name: "--config", Description: "Specify the path to the elm compiler"},
			{Name: "--check-after-tests", Description: "Checks whether there are uncommitted suppression files"},
			{Name: "--unsuppress", Description: "Include suppressed errors in the error report for all rules"},
			{Name: "--unsuppress-rules", Description: "Include suppressed errors in the error report for all rules"},
			{Name: "--compiler", Description: "Specify the path to the elm compiler"},
			{Name: "--rules", Description: "Run with a subsection of the rules in the configuration"},
			{Name: "--watch", Description: "Re-run elm-review automatically when your project or configuration changes"},
			{Name: "--watch-code", Description: "Re-run elm-review automatically when your project changes"},
			{Name: "--elmjson", Description: "Specify the path to the elm.json file of the project"},
			{Name: "--template", Description: "Use the review configuration from a GitHub repository"},
			{Name: "--version", Description: "Print the version of the elm-review CLI"},
			{Name: "--debug", Description: "Add helpful information to debug your configuration or rules"},
			{Name: "--report", Description: "Error reports will be in JSON format"},
			{Name: "--no-details", Description: "Hide the details from error reports for a more compact view"},
			{Name: "--ignore-dirs", Description: "Ignore the reports of all rules for the specified directories"},
			{Name: "--ignore-files", Description: "Ignore the reports of all rules for the specified files"},
			{Name: "--fix", Description: "Specify the path to elm-format"},
			{Name: "--fix-limit", Description: "Limit the number of fixes applied in a single batch to N"},
			{Name: "--extract", Description: "Disable colors in the output"},
		},
	})
}
