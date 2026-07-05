package cc

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "typst",
		Description: "The Typst compiler",
		Subcommands: []spec.Subcommand{
			{Name: "compile", Description: "Compiles an input file into a supported output format"},
			{Name: "input", Description: "Path to input Typst file"},
			{Name: "output", Description: "Path to output file (PDF, PNG, or SVG)"},
			{Name: "watch", Description: "Watches an input file and recompiles on changes"},
			{Name: "query", Description: "Processes an input file to extract provided metadata"},
			{Name: "selector", Description: "Defines which elements to retrieve"},
			{Name: "fonts", Description: "Lists all discovered fonts in system and custom font paths"},
			{Name: "update", Description: "Self update the Typst CLI (disabled)"},
			{Name: "version", Description: "Which version to update to (defaults to latest)"},
			{Name: "help", Description: "Print this message or the help of the given subcommand(s)"},
		},
		Options: []spec.Option{
			{Name: "--root", Description: "Configures the project root (for absolute paths)"},
			{Name: "--font-path", Description: "Adds additional directories to search for fonts"},
			{Name: "--diagnostic-format", Description: "The format to emit diagnostics in"},
			{Name: "-f", Description: "The format of the output file, inferred from the extension by default"},
			{Name: "--open", Description: "Opens the output file using the default viewer after compilation"},
			{Name: "--ppi", Description: "The PPI (pixels per inch) to use for PNG export"},
			{Name: "--flamegraph", Description: "Produces a flamegraph of the compilation process"},
			{Name: "-h", Description: "Print help"},
			{Name: "--field", Description: "Extracts just one field from all retrieved elements"},
			{Name: "--format", Description: "The format to serialize in"},
			{Name: "--one", Description: "Expects and retrieves exactly one element"},
			{Name: "--variants", Description: "Also lists style variants of each font family"},
			{Name: "--force", Description: "Forces a downgrade to an older version (required for downgrading)"},
			{Name: "--cert", Description: "Path to a custom CA certificate to use when making network requests"},
			{Name: "-v", Description: "Print help"},
			{Name: "-V", Description: "Print version"},
		},
	})
}
