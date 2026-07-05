package python

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "sqlfluff",
		Description: "A dialect-flexible and configurable SQL linter",
		Subcommands: []core.Subcommand{
			{Name: "lint", Description: "Lint SQL files via passing a list of files or using stdin"},
			{Name: "templater", Description: "Name of templater to use, eg. raw"},
			{Name: "dialect", Description: "Name of dialect, eg. ANSI"},
			{Name: "logger", Description: "Name of logger to limit to, eg. templater"},
			{Name: "annotation-level", Description: "Level of annotation, eg. notice"},
			{Name: "fix", Description: "Fix SQL files"},
			{Name: "parse", Description: "Parse SQL files and just spit out the result"},
			{Name: "dialects", Description: "Show the current dialects available"},
			{Name: "version", Description: "Show the version of sqlfluff"},
			{Name: "rules", Description: "Show the current rules in use"},
		},
		Options: []core.Option{
			{Name: "--version", Description: "Show the version and exit"},
			{Name: "--help", Description: "Show help for sqlfluff"},
			{Name: "--nocolor", Description: "No color - output will be without ANSI color codes"},
			{Name: "--ignore", Description: "Narrow the search to only specific rules"},
			{Name: "--templater", Description: "The templater to use (default=jinja)"},
			{Name: "--dialect", Description: "The dialect of SQL to lint"},
			{Name: "--format", Description: "What format to return the lint result in (default=human)"},
			{Name: "--processes", Description: "Set this flag to ignore inline noqa comments"},
			{Name: "--bench", Description: "Set this flag to engage the benchmarking tool output"},
			{Name: "--logger", Description: "Choose to limit the logging to one of the loggers"},
			{Name: "--encoding", Description: "Specify encoding to use when reading and writing files. Defaults to autodetect"},
			{Name: "--ignore-local-config", Description: "Level of annotation, eg. notice"},
			{Name: "--disregard-sqlfluffignores", Description: "Perform the operation regardless of .sqlfluffignore configurations"},
			{Name: "--disable-progress-bar", Description: "Disables progress bars"},
			{Name: "--nofail", Description: "Show help for sqlfluff"},
			{Name: "--disable-noqa", Description: "Set this flag to ignore inline noqa comments"},
			{Name: "--FIX-EVEN-UNPARSABLE", Description: "Enables fixing of files that have templating or parse errors"},
			{Name: "--show-lint-violations", Description: "Show lint violations"},
			{Name: "--code-only", Description: "Output only the code elements of the parse tree"},
			{Name: "--include-meta", Description: "What format to return the lint result in (default=human)"},
			{Name: "--recurse", Description: "The depth to recursively parse to (0 for unlimited)"},
			{Name: "--verbose", Description: "Show the version of sqlfluff"},
		},
	})
}
