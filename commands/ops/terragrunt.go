package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "terragrunt",
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
			{Name: "-h", Description: "Show this help output, or the help for a specified subcommand"},
			{Name: "-v", Description: "Show the current Terragrunt version"},
			{Name: "--terragrunt-config", Description: "Path to the Terragrunt config file. Default is terragrunt.hcl"},
			{Name: "--terragrunt-tfpath", Description: "Path to the Terraform binary. Default is terraform (on PATH)"},
			{Name: "--terragrunt-no-auto-init", Description: "Don't automatically re-run command in case of transient errors"},
			{Name: "--terragrunt-non-interactive", Description: "Assume 'yes' for all prompts"},
			{Name: "--terragrunt-working-dir", Description: "The path to the Terraform templates. Default is current directory"},
			{Name: "--terragrunt-download-dir", Description: "*-all commands continue processing components even if a dependency fails"},
			{Name: "--terragrunt-iam-role", Description: "Unix-style glob of directories to exclude when running *-all commands"},
			{Name: "--terragrunt-include-dir", Description: "Unix-style glob of directories to include when running *-all commands"},
			{Name: "--terragrunt-strict-include", Description: "*-all commands will be run disregarding the dependencies"},
			{Name: "--terragrunt-parallelism", Description: "*-all commands parallelism set to at most N modules"},
			{Name: "--terragrunt-debug", Description: "Write terragrunt-debug.tfvars to working folder to help root-cause issues"},
			{Name: "--terragrunt-log-level", Description: "Enable check mode in the hclfmt command"},
			{Name: "--terragrunt-hclfmt-file", Description: "The path to a single hcl file that the hclfmt command should run on"},
			{Name: "--terragrunt-override-attr", Description: "Prepare your working directory for other commands"},
			{Name: "-upgrade", Description: "Check whether the configuration is valid"},
			{Name: "-compact-warnings", Description: "If set, generates a plan to destroy all the known resources"},
			{Name: "-detailed-exitcode", Description: "Return a detailed exit code when the command exits"},
			{Name: "-out", Description: "The path to save the generated execution plan"},
			{Name: "-parallelism", Description: "Update the state prior to checking for differences"},
			{Name: "-state", Description: "A Resource Address to target. This flag can be used multiple times"},
			{Name: "-var", Description: "Set variables in the Terraform configuration from a variable file"},
			{Name: "-update", Description: "Generate a Graphviz graph of the steps in an operation"},
			{Name: "-allow-missing", Description: "Remove the 'tainted' state from a resource instance"},
		},
	})
}
