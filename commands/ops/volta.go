package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "volta",
		Description: "Enables verbose diagnostics",
		Subcommands: []spec.Subcommand{
			{Name: "fetch", Description: "Fetches a tool to the local machine"},
			{Name: "install", Description: "Installs a tool in your toolchain"},
			{Name: "uninstall", Description: "Uninstalls a tool from your toolchain"},
			{Name: "pin", Description: "Pins your project's runtime or package manager"},
			{Name: "list", Description: "Displays the current toolchain"},
			{Name: "completions", Description: "Generates Volta completions"},
			{Name: "shell", Description: "Shell to generate completions for"},
			{Name: "which", Description: "Locates the actual binary that will be called by Volta"},
			{Name: "setup", Description: "Enables Volta for the current user"},
			{Name: "run", Description: "Run a command with custom Node, npm, and/or Yarn versions"},
			{Name: "help", Description: "Prints this message or the help of the given subcommand(s)"},
		},
		Options: []spec.Option{
			{Name: "--verbose", Description: "Enables verbose diagnostics"},
			{Name: "--quiet", Description: "Prevents unnecessary output"},
			{Name: "-h", Description: "Prints help information"},
			{Name: "-c", Description: "Show the currently-active tool(s)"},
			{Name: "-d", Description: "Show your default tool(s)"},
			{Name: "--format", Description: "Specify output format"},
			{Name: "-f", Description: "Write over an existing file, if any"},
			{Name: "-o", Description: "File to write generated completions to"},
			{Name: "--bundle", Description: "Forces npm to be the version bundled with Node"},
			{Name: "--no-yarn", Description: "Disables Yarn"},
			{Name: "--node", Description: "Set the custom Node version"},
			{Name: "--npm", Description: "Set the custom npm version"},
			{Name: "--yarn", Description: "Set the custom Yarn version"},
			{Name: "--env", Description: "Set an environment variable (can be used multiple times)"},
			{Name: "-v", Description: "Prints the current version of Volta"},
		},
	})
}
