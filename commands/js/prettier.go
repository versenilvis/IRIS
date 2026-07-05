package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "prettier",
		Description: "Run Prettier from the command line",
		Options: []spec.Option{
			{Name: "-c", Description: "Check if your files are formatted"},
			{Name: "-l", Description: "Print the names of files that are different from Prettier's formatting"},
			{Name: "-w", Description: "Edit files in-place"},
			{Name: "--arrow-parens", Description: "Include parentheses around a sole arrow function parameter"},
			{Name: "--no-bracket-spacing", Description: "Do not print spaces between brackets"},
			{Name: "--end-of-line", Description: "Which end of line characters to apply"},
			{Name: "--html-whitespace-sensitivity", Description: "How to handle whitespaces in HTML"},
			{Name: "--jsx-bracket-same-line", Description: "Put > on the last line instead of at a new line"},
			{Name: "--jsx-single-quote", Description: "Use single quotes in JSX"},
			{Name: "--parser", Description: "Which parser to use"},
			{Name: "--print-width", Description: "The line length where Prettier will try wrap"},
			{Name: "--prose-wrap", Description: "How to wrap prose"},
			{Name: "--quote-props", Description: "Change when properties in objects are quoted"},
			{Name: "--no-semi", Description: "Do not print semicolons, except at the beginning of lines which may need them"},
			{Name: "--single-quote", Description: "Use single quotes instead of double quotes"},
			{Name: "--tab-width", Description: "Number of spaces per indentation level"},
			{Name: "--trailing-comma", Description: "Print trailing commas wherever possible when multi-line"},
			{Name: "--use-tabs", Description: "Indent with tabs instead of spaces"},
			{Name: "--vue-indent-script-and-style", Description: "Indent script and style tags in Vue files"},
			{Name: "--config", Description: "Do not look for a configuration file"},
			{Name: "--config-precedence", Description: "Define in which order config files and CLI options should be evaluated"},
			{Name: "--no-editorconfig", Description: "Don't take .editorconfig into account when parsing configuration"},
			{Name: "--find-config-path", Description: "Finds a path to the configuration file for the given input file"},
			{Name: "--ignore-path", Description: "Path to a file with patterns describing files to ignore"},
			{Name: "--plugin", Description: "Add a plugin"},
			{Name: "--plugin-search-dir", Description: "Custom directory that contains prettier plugins in node_modules subdirectory"},
			{Name: "--with-node-modules", Description: "Process files inside 'node_modules' directory"},
			{Name: "--cursor-offset", Description: "Format code ending at a given character offset (exclusive)"},
			{Name: "--range-start", Description: "Format code starting at a given character offset"},
			{Name: "--no-color", Description: "Do not colorize error messages"},
			{Name: "--file-info", Description: "Extract the following info (as JSON) for a given file path"},
			{Name: "-h", Description: "Show CLI usage, or details about the given flag"},
			{Name: "-u", Description: "Ignore unknown files"},
			{Name: "--insert-pragma", Description: "Insert @format pragma into file's first docblock comment"},
			{Name: "--loglevel", Description: "What level of logs to report"},
			{Name: "--require-pragma", Description: "Path to the file to pretend that stdin comes from"},
			{Name: "--support-info", Description: "Print support information as JSON"},
			{Name: "-v", Description: "Print Prettier version"},
			{Name: "--debug-check", Description: "Prevent errors when pattern is unmatched"},
		},
	})
}
