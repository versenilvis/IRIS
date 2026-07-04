package text

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "vale",
		Description: "A syntax-aware linter for prose built with speed and extensibility in mind",
		Subcommands: []core.Subcommand{
			{Name: "ls-config", Description: "Print the current configuration to stdout"},
			{Name: "ls-metrics", Description: "Print the given file's internal metrics to stdout"},
			{Name: "file", Description: "The path to a file you want to analyze"},
		},
		Options: []core.Option{
			{Name: "--help", Description: "Show help for vale"},
			{Name: "--version", Description: "Print the current version"},
			{Name: "--ignore-syntax", Description: "Lint all files line-by-line"},
			{Name: "--no-exit", Description: "Don't return a nonzero exit code on errors"},
			{Name: "--no-wrap", Description: "Don't wrap CLI output"},
			{Name: "--ext", Description: "An extension to associate with stdin"},
			{Name: "--glob", Description: "A glob pattern"},
			{Name: "--minAlertLevel", Description: "The minimum level to display"},
			{Name: "--output", Description: "The alert output style to use"},
			{Name: "--config", Description: "A path to a .vale.ini file"},
		},
	})
}
