package git

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "github",
		Description: "Open a git repository in GitHub Desktop",
		Subcommands: []spec.Subcommand{
			{Name: "clone", Description: "Clone a repository"},
			{Name: "open", Description: "Open a git repository in GitHub Desktop"},
			{Name: "help", Description: "Show the help page for a command"},
		},
		Options: []spec.Option{
			{Name: "--help", Description: "Show the help page for a command"},
			{Name: "--branch", Description: "The branch to checkout after cloning"},
		},
	})
}
