package git

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "github",
		Description: "Open a git repository in GitHub Desktop",
		Subcommands: []core.Subcommand{
			{Name: "clone", Description: "Clone a repository"},
			{Name: "open", Description: "Open a git repository in GitHub Desktop"},
			{Name: "help", Description: "Show the help page for a command"},
		},
		Options: []core.Option{
			{Name: "--help", Description: "Show the help page for a command"},
			{Name: "--branch", Description: "The branch to checkout after cloning"},
		},
	})
}
