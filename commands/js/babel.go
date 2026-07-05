package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "babel",
		Description: "A comma-separated list of preset names",
		Options: []spec.Option{
			{Name: "-f", Description: "A comma-separated list of preset names"},
			{Name: "--plugins", Description: "A comma-separated list of plugin names"},
			{Name: "--config-file", Description: "Path to a .babelrc file to use"},
			{Name: "--env-name", Description: "The project-root resolution mode"},
			{Name: "--source-type", Description: "Whether or not to look up .babelrc and .babelignore files"},
			{Name: "--ignore", Description: "List of glob paths to **not** compile"},
			{Name: "--only", Description: "List of glob paths to **only** compile"},
			{Name: "--no-highlight-code", Description: "Enable or disable ANSI syntax highlighting of code frames"},
			{Name: "--no-comments", Description: "Write comments to generated output"},
			{Name: "--retain-lines", Description: "Retain line numbers. This will result in really ugly code"},
			{Name: "--compact", Description: "Do not include superfluous whitespace characters and line terminators"},
			{Name: "--minified", Description: "Save as many bytes when printing. (false by default)"},
			{Name: "--auxiliary-comment-before", Description: "Print a comment before any injected non-user code"},
			{Name: "--auxiliary-comment-after", Description: "Print a comment after any injected non-user code"},
			{Name: "-s", Description: "Set `file` on returned source map"},
			{Name: "--source-file-name", Description: "Set `sources[0]` on returned source map"},
			{Name: "--source-root", Description: "The root from which all sources are relative"},
			{Name: "-x", Description: "Preserve the file extensions of the input files"},
			{Name: "-w", Description: "Recompile files on changes"},
			{Name: "--skip-initial-build", Description: "Do not compile files before watching"},
			{Name: "-o", Description: "Compile all input files into a single file"},
			{Name: "-d", Description: "Compile an input directory of modules into an output directory"},
			{Name: "--relative", Description: "Compile into an output directory relative to input directory or file"},
			{Name: "-D", Description: "When compiling a directory copy over non-compilable files"},
			{Name: "--include-dotfiles", Description: "Include dotfiles when compiling and copying non-compilable files"},
			{Name: "--no-copy-ignored", Description: "Exclude ignored files when copying non-compilable files"},
			{Name: "--verbose", Description: "Log everything. This option conflicts with --quiet"},
			{Name: "--quiet", Description: "Don't log anything. This option conflicts with --verbose"},
			{Name: "--delete-dir-on-start", Description: "Delete the out directory before compilation"},
			{Name: "--out-file-extension", Description: "Use a specific extension for the output files"},
			{Name: "-V", Description: "Output the version number"},
			{Name: "-h", Description: "Output usage information"},
		},
	})
}
