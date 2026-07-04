package rust

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "tokei",
		Description: "Count your code, quickly",
		Options: []core.Option{
			{Name: "-f", Description: "Will print out statistics on individual files"},
			{Name: "-h", Description: "Prints help information"},
			{Name: "--hidden", Description: "Count hidden files"},
			{Name: "-l", Description: "Prints out supported languages and their extensions"},
			{Name: "--no-ignore", Description: "Don't respect ignore files (.gitignore, .ignore, etc.)"},
			{Name: "--no-ignore-dot", Description: "Don't respect ignore files (.gitignore, .ignore, etc.) in parent directories"},
			{Name: "--no-ignore-vcs", Description: "Prints version information"},
			{Name: "-v", Description: "Set log output level:"},
			{Name: "-c", Description: "Sets a strict column width of the output, only available for terminal output"},
			{Name: "-e", Description: "Ignore all files & directories matching the pattern"},
			{Name: "-i", Description: "Outputs Tokei in a specific format"},
			{Name: "-s", Description: "Sort languages based on column"},
			{Name: "-t", Description: "Filters output by language type, separated by a comma. i.e. -t=Rust,Markdown"},
		},
	})
}
