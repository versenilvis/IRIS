package view

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "idea",
		Description: "IntelliJ IDEA CLI",
		Subcommands: []core.Subcommand{
			{Name: "diff", Description: "Open the diff viewer to see the differences between two specified files"},
			{Name: "merge", Description: "Open the Merge dialog to merge the specified files"},
			{Name: "format", Description: "Apply code style formatting to the specified files"},
			{Name: "inspect", Description: "Perform code inspection on the specified project"},
		},
		Options: []core.Option{
			{Name: "-h", Description: "Show help for format command"},
			{Name: "-m", Description: "Process specified directories recursively"},
			{Name: "-s", Description: "Specify the code style settings file to use for formatting"},
			{Name: "-allowDefaults", Description: "Preserve encoding and enforce the charset for reading and writing source files"},
			{Name: "-d", Description: "Run the formatter in the validation mode"},
			{Name: "-changes", Description: "Run inspections only on local uncommitted changes"},
			{Name: "-v", Description: "Set the verbosity level of the output"},
			{Name: "--wait", Description: "Wait for the files to be closed before returning to the command prompt"},
		},
	})
}
