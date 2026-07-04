package text

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "typos",
		Description: "Source code spelling correction",
		Options: []core.Option{
			{Name: "-c", Description: "Custom config file"},
			{Name: "--isolated", Description: "Ignore implicit configuration files"},
			{Name: "--diff", Description: "Print a diff of what would change"},
			{Name: "-w", Description: "Write fixes out"},
			{Name: "--files", Description: "Debug: Print each file that would be spellchecked"},
			{Name: "--identifiers", Description: "Debug: Print each identifier that would be spellchecked"},
			{Name: "--words", Description: "Debug: Print each word that would be spellchecked"},
			{Name: "--dump-config", Description: "Print to stdout"},
			{Name: "--type-list", Description: "Show all supported file types"},
			{Name: "--format", Description: "Set the output format"},
			{Name: "-j", Description: "The approximate number of threads to use"},
			{Name: "--exclude", Description: "Ignore files & directories matching the glob"},
			{Name: "--hidden", Description: "Search hidden files and directories"},
			{Name: "--no-ignore", Description: "Don't respect ignore files"},
			{Name: "--no-ignore-dot", Description: "Don't respect .ignore files"},
			{Name: "--no-ignore-global", Description: "Don't respect global ignore files"},
			{Name: "--no-ignore-parent", Description: "Don't respect ignore files in parent directories"},
			{Name: "--no-ignore-vcs", Description: "Don't respect ignore files in vcs directories"},
			{Name: "--binary", Description: "Search binary files"},
			{Name: "--no-check-filenames", Description: "Skip verifying spelling in file names"},
			{Name: "--no-check-files", Description: "Skip verifying spelling in files"},
			{Name: "--no-unicode", Description: "Only allow ASCII characters in identifiers"},
			{Name: "--locale", Description: "Set the locale to use"},
			{Name: "--color", Description: "Controls when to use color"},
			{Name: "-v", Description: "More output per occurrence"},
			{Name: "-q", Description: "Less output per occurrence"},
			{Name: "-h", Description: "Print help information"},
			{Name: "-V", Description: "Print version information"},
		},
	})
}
