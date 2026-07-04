package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "limactl",
		Description: "Lima: Linux virtual machines, with a focus on running containers",
		Subcommands: []core.Subcommand{
			{Name: "completion", Description: "Generate the autocompletion script for the specified shell"},
			{Name: "copy", Description: "Copy files between host and guest"},
			{Name: "delete", Description: "Delete an instance of Lima"},
			{Name: "info", Description: "Show diagnostic information"},
			{Name: "list", Description: "List instances of Lima"},
			{Name: "prune", Description: "Prune garbage objects"},
			{Name: "show-ssh", Description: "Show the ssh command line"},
			{Name: "cmd", Description: "Full ssh command line"},
			{Name: "options", Description: "Ssh option key value pairs"},
			{Name: "config", Description: "~/.ssh/config format"},
			{Name: "stop", Description: "Stop an instance"},
			{Name: "validate", Description: "Validate YAML files"},
		},
		Options: []core.Option{
			{Name: "-h", Description: "Debug mode"},
			{Name: "--no-descriptions", Description: "Disable completion descriptions"},
			{Name: "-r", Description: "Copy directories recursively"},
			{Name: "-f", Description: "Forcibly kill the processes"},
			{Name: "--json", Description: "JSONify output"},
			{Name: "-q", Description: "Only show names"},
			{Name: "--workdir", Description: "Working directory"},
			{Name: "--tty", Description: "Stop an instance"},
			{Name: "--check", Description: "Validate YAML files"},
			{Name: "-v", Description: "Version for limactl"},
		},
	})
}
