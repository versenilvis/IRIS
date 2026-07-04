package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "esbuild",
		Description: "An extremely fast JavaScript bundler",
		Options: []core.Option{
			{Name: "--bundle", Description: "Bundle all dependencies into the output files"},
			{Name: "--define", Description: "Replace variable names with a literal value, eg. --define:DEBUG=true"},
			{Name: "--external", Description: "Exclude modules from the build"},
			{Name: "--format", Description: "The output format"},
			{Name: "--loader", Description: "For a given file extension, specify a loader"},
			{Name: "--minify", Description: "Minify the output (sets all the --minify-* options)"},
			{Name: "--outdir", Description: "The output directory for multiple entrypoints"},
			{Name: "--outfile", Description: "The output file for one entrypoint"},
			{Name: "--platform", Description: "The platform target"},
			{Name: "--serve", Description: "Start a local HTTP server on this host:port"},
			{Name: "--splitting", Description: "Enable code splitting"},
			{Name: "--target", Description: "Rebuild on file system changes"},
			{Name: "--allow-overwrite", Description: "Allow output files to overwrite input files"},
			{Name: "--analyze", Description: "Print a report about the contents of the bundle"},
			{Name: "--asset-names", Description: "Path template for 'file' loader files"},
			{Name: "--banner", Description: "Text to be prepended to each output file type"},
			{Name: "--charset", Description: "Use UTF-8 instead of escaped codepoints in ASCII"},
			{Name: "--chunk-names", Description: "Path template to use for code splitting chunks"},
			{Name: "--color", Description: "Force use of terminal colors"},
			{Name: "--drop", Description: "Remove certain constructs"},
			{Name: "--entry-names", Description: "Path template to use for entry point output paths"},
			{Name: "--footer", Description: "Text to be appended to each file type"},
			{Name: "--global-name", Description: "The name of the global if using --format=iife"},
			{Name: "--ignore-annotations", Description: "Enable this to work with packages that have incorrect tree-shaking annotations"},
			{Name: "--inject", Description: "Import the file into all input files, automatically replace matching globals"},
			{Name: "--jsx-factory", Description: "What to use for the JSX factory"},
			{Name: "--jsx-fragment", Description: "What to use for the JS Fragment factory"},
			{Name: "--jsx", Description: "Preserve JSX instead of transforming"},
			{Name: "--jsx-dev", Description: "Toggles development mode for the automatic runtime"},
			{Name: "--jsx-import-source", Description: "Overrides the root import for runtime functions (default: react)"},
			{Name: "--keep-names", Description: "Preserve 'name' on functions and classes"},
			{Name: "--legal-comments", Description: "Where to place legal comments"},
			{Name: "--log-level", Description: "Set the log level"},
			{Name: "--log-limit", Description: "Maximum message count, 0 to disable"},
			{Name: "--log-override", Description: "For a particular identifier, set the log level"},
			{Name: "--main-fields", Description: "Override the main file order in package.json"},
			{Name: "--mangle-cache", Description: "Save 'mangle props' decisions to a JSON file"},
			{Name: "--mangle-props", Description: "Rename all properties matching a regular expression"},
			{Name: "--mangle-quoted", Description: "Enable mangling (renaming) quoted properties"},
			{Name: "--metafile", Description: "Write metadata about the build to a JSON file"},
		},
	})
}
