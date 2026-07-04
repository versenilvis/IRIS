package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "terraform",
		Description: "Workspace",
		Subcommands: []core.Subcommand{
			{Name: "show", Description: "Display the current workspace"},
			{Name: "list", Description: "List the workspace"},
			{Name: "delete", Description: "Delete the specified workspace"},
			{Name: "select", Description: "Change the current working workspace"},
		},
		Options: []core.Option{
			{Name: "-lock", Description: "Lock the state file when locking is supported. Defaults to true"},
			{Name: "-force", Description: "Delete the workspace even if its state is not empty. Defaults to false"},
			{Name: "-lock-timeout", Description: "Duration to retry a state lock. Default 0s"},
			{Name: "-input", Description: "Ask for input for variables if not directly set"},
			{Name: "-no-color", Description: "Disables output with coloring"},
			{Name: "-help", Description: "Show this help output, or the help for a specified subcommand"},
			{Name: "-chdir", Description: "Switch to a different working directory before executing the given subcommand"},
			{Name: "-version", Description: "Show the current Terraform version"},
			{Name: "-upgrade", Description: "Check whether the configuration is valid"},
			{Name: "-compact-warnings", Description: "If set, generates a plan to destroy all the known resources"},
			{Name: "-detailed-exitcode", Description: "Return a detailed exit code when the command exits"},
			{Name: "-out", Description: "The path to save the generated execution plan"},
			{Name: "-parallelism", Description: "Update the state prior to checking for differences"},
			{Name: "-state", Description: "A Resource Address to target. This flag can be used multiple times"},
			{Name: "-var", Description: "Set variables in the Terraform configuration from a variable file"},
			{Name: "-update", Description: "Generate a Graphviz graph of the steps in an operation"},
			{Name: "-allow-missing", Description: "Remove the 'tainted' state from a resource instance"},
			{Name: "-install-autocomplete", Description: "Install bash/zsh tab completion"},
			{Name: "-uninstall-autocomplete", Description: "Uninstall bash/zsh tab completion"},
		},
	})
}
