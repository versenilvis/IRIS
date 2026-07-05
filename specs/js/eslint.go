package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "eslint",
		Description: "Pluggable JavaScript linter",
		Options: []core.Option{
			{Name: "--no-eslintrc", Description: "Disable use of configuration from .eslintrc.*"},
			{Name: "-c", Description: "Use this configuration, overriding .eslintrc.* config options if present"},
			{Name: "--env", Description: "Specify environments"},
			{Name: "--ext", Description: "Specify JavaScript file extensions"},
			{Name: "--global", Description: "Define global variables"},
			{Name: "--parser", Description: "Specify the parser to be used"},
			{Name: "--parser-options", Description: "Specify parser options"},
			{Name: "--resolve-plugins-relative-to", Description: "A folder where plugins should be resolved from"},
			{Name: "--rulesdir", Description: "Use additional rules from this directory"},
			{Name: "--plugin", Description: "Specify plugins"},
			{Name: "--rule", Description: "Specify rules"},
			{Name: "--fix", Description: "Automatically fix problems"},
			{Name: "--fix-dry-run", Description: "Automatically fix problems without saving the changes to the file system"},
			{Name: "--fix-type", Description: "Specify the types of fixes to apply"},
			{Name: "--ignore-path", Description: "Specify path of ignore file"},
			{Name: "--no-ignore", Description: "Disable use of ignore files and patterns"},
			{Name: "--ignore-pattern", Description: "Pattern of files to ignore (in addition to those in .eslintignore)"},
			{Name: "--stdin", Description: "Lint code provided on <STDIN>"},
			{Name: "--stdin-filename", Description: "Specify filename to process STDIN as"},
			{Name: "--quiet", Description: "Report errors only"},
			{Name: "--max-warnings", Description: "Number of warnings to trigger nonzero exit code"},
			{Name: "-o", Description: "Specify file to write report to"},
			{Name: "-f", Description: "Use a specific output format"},
			{Name: "--color", Description: "Force enabling of color"},
			{Name: "--no-color", Description: "Force disabling of color"},
			{Name: "--no-inline-config", Description: "Prevent comments from changing config or rules"},
			{Name: "--cache", Description: "Only check changed files"},
			{Name: "--cache-location", Description: "Path to the cache file or directory"},
			{Name: "--cache-strategy", Description: "Strategy to use for detecting changed files"},
			{Name: "--init", Description: "Run config initialization wizard"},
			{Name: "--env-info", Description: "Output execution environment information"},
			{Name: "--debug", Description: "Output debugging information"},
			{Name: "-h", Description: "Show help"},
			{Name: "-v", Description: "Output the version number"},
			{Name: "--print-config", Description: "Print the configuration for the give file"},
		},
	})
}
