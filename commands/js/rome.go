package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "rome",
		Description: "Rome CLI",
		Subcommands: []core.Subcommand{
			{Name: "check", Description: "Run the linter on a set of files"},
			{Name: "ci", Description: "Run the linter and formatter check on a set of files"},
			{Name: "format", Description: "Run the formatter on a set of files"},
			{Name: "init", Description: "Bootstraps a new rome project"},
			{Name: "start", Description: "Start the Rome daemon server process"},
			{Name: "stop", Description: "Stop the Rome daemon server process"},
			{Name: "lsp-proxy", Description: "Acts as a server for the Language Server Protocol over stdin/stdout"},
			{Name: "rage", Description: "Prints information for debugging"},
			{Name: "version", Description: "Shows the Rome version information and quit"},
			{Name: "help", Description: "Prints help message"},
		},
		Options: []core.Option{
			{Name: "--colors", Description: "Set the formatting mode for markup"},
			{Name: "--use-server", Description: "Connect to a running instance of the Rome daemon server"},
			{Name: "--version", Description: "Show the Rome version information and quit"},
			{Name: "--files-max-size", Description: "The maximum allowed size for source code files in bytes"},
			{Name: "--apply", Description: "Apply safe fixes"},
			{Name: "--apply-unsafe", Description: "Apply safe and unsafe fixes"},
			{Name: "--max-diagnostics", Description: "Cap the amount of diagnostics displayed"},
			{Name: "--config-path", Description: "Set the filesystem path to the config dir of the rome.json file"},
			{Name: "--verbose", Description: "Print additional verbose advices on diagnostics"},
			{Name: "--formatter-enabled", Description: "Allow to enable or disable the formatter check"},
			{Name: "--linter-enabled", Description: "Allow to enable or disable the linter check"},
			{Name: "--organize-imports-enabled", Description: "Allow to enable or disable the organize imports"},
			{Name: "--indent-style", Description: "Change the indention character"},
			{Name: "--indent-size", Description: "How many spaces should be used for indentation"},
			{Name: "--line-width", Description: "How many characters the formatter is allowed to print in a single line"},
			{Name: "--quote-style", Description: "Changes the quotation character for strings"},
			{Name: "--quote-properties", Description: "Changes when properties in object should be quoted"},
			{Name: "--trailing-comma", Description: "Changes trailing commas in multi-line comma-separated syntactic structures"},
			{Name: "--semicolons", Description: "Changes when to print semicolons for statements"},
			{Name: "--write", Description: "Edit the files in place (beware!) instead of printing the diff to the console"},
			{Name: "--skip-errors", Description: "Skip over files containing syntax errors instead of emitting an error diagnostic"},
			{Name: "--stdin-file-path", Description: "A file name with its extension to pass when reading from standard in"},
		},
	})
}
