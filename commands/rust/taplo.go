package rust

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "taplo",
		Description: "Set color values for the output",
		Subcommands: []spec.Subcommand{
			{Name: "config", Description: "Operations with the Taplo config file"},
			{Name: "default", Description: "Print the default `.taplo.toml` configuration file"},
			{Name: "help", Description: "Print this message or the help of the given subcommand(s)"},
			{Name: "schema", Description: "Print the JSON schema of the `.taplo.toml` configuration file"},
			{Name: "format", Description: "Format TOML documents"},
			{Name: "FILES ...", Description: "Paths or glob patterns to TOML documents"},
			{Name: "f", Description: "Force formatting of files"},
			{Name: "get", Description: "Extract a value from the given TOML document"},
			{Name: "lint", Description: "Lint a TOML documents"},
			{Name: "lsp", Description: "Language server operations"},
			{Name: "stdio", Description: "Run the language server over the standard input and output"},
			{Name: "tcp", Description: "Run the language server and listen on a TCP address"},
		},
		Options: []spec.Option{
			{Name: "--colors", Description: "Set color values for the output"},
			{Name: "--verbose", Description: "Enable verbose logging format"},
			{Name: "--log-spans", Description: "Enable logging spans"},
			{Name: "--help", Description: "Print help information for config"},
			{Name: "--config", Description: "Path to the Taplo configuration file"},
			{Name: "--cache-path", Description: "Set a cache path"},
			{Name: "--check", Description: "Report any files that are not correctly formatted"},
			{Name: "--diff", Description: "Print the differences in patch formatting to `stdout`"},
			{Name: "--no-auto-config", Description: "Do not search for a configuration file"},
			{Name: "--option", Description: "A formatter option given as a 'key=value', can be set multiple times"},
			{Name: "--stdin-filepath", Description: "A path to the file that the taplo will treat like stdin"},
			{Name: "--file-path", Description: "Path to the TOML document"},
			{Name: "-o", Description: "The format specifying how the output is printed"},
			{Name: "--strip-newline", Description: "Strip the trailing newline from the output"},
			{Name: "--default-schema-catalogs", Description: "Use the default online catalogs for schemas"},
			{Name: "--no-schema", Description: "Disable all schema validation"},
			{Name: "--schema", Description: "URL to the schema to be used for validation"},
			{Name: "--schema-catalog", Description: "URL to the schema catalog to be used for validation"},
			{Name: "--version", Description: "Print version information for taplo"},
		},
	})
}
