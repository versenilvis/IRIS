package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "oxlint",
		Description: "All lints (except nursery)",
		Options: []core.Option{
			{Name: "-c", Description: "Path to Oxlint JSON configuration file"},
			{Name: "--tsconfig", Description: "Specify the file to use as your .eslintignore"},
			{Name: "--ignore-pattern", Description: "Specify file patterns to ignore (in addition to those in .eslintignore)"},
			{Name: "--no-ignore", Description: "Follow symlinks when linting, which are ignored by default"},
			{Name: "-D", Description: "Deny a lint rule or category (enable a lint as an error)"},
			{Name: "-W", Description: "Warn about a lint rule or category (enable a lint as a warning)"},
			{Name: "-A", Description: "Allow a lint rule or category (suppress a lint)"},
			{Name: "--fix", Description: "Fix as many issues as possible. Only unfixed issues are reported in the output"},
			{Name: "--silent", Description: "Do not display any diagnostics"},
			{Name: "--deny-warnings", Description: "Exit with a non-zero code if there are any warnings"},
			{Name: "--max-warnings", Description: "Exit with a non-zero code if there are more than `max` warnings"},
			{Name: "-f", Description: "Use a specific output format"},
			{Name: "--threads", Description: "Number of threads to use. Set to 1 to use only 1 CPU core"},
			{Name: "--rules", Description: "List all available rules"},
			{Name: "-h", Description: "Show help"},
			{Name: "-V", Description: "Show version"},
		},
	})
}
