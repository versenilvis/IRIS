package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "railway",
		Description: "CLI for managing Railway Apps",
		Subcommands: []core.Subcommand{
			{Name: "add", Description: "Add a plugin to your project"},
			{Name: "completion", Description: "Generate shell-completions"},
			{Name: "shell", Description: "The shell to generate completions for"},
			{Name: "connect", Description: "Connect to a plugin"},
			{Name: "plugin", Description: "The plugin to connect to"},
			{Name: "delete", Description: "Delete a project"},
			{Name: "project-id", Description: "The project to delete"},
			{Name: "docs", Description: "Open Railway Documentation in default browser"},
			{Name: "down", Description: "Remove the most recent deployment"},
			{Name: "environment", Description: "Change your environment"},
			{Name: "help", Description: "Get help about any command"},
			{Name: "command", Description: "The command to get help about"},
			{Name: "init", Description: "Create a new railway project"},
			{Name: "link", Description: "Connect to an existing railway project"},
			{Name: "list", Description: "List all railway projects"},
			{Name: "login", Description: "Login to railway"},
			{Name: "logout", Description: "Logout of railway"},
			{Name: "logs", Description: "Show logs for the most-recent deployment"},
			{Name: "Number of lines", Description: "The number of lines to output"},
			{Name: "open", Description: "Open the project"},
			{Name: "live", Description: "Open the live project"},
			{Name: "metrics", Description: "Open project metrics"},
			{Name: "settings", Description: "Open project settings"},
			{Name: "run", Description: "Run a local command using variables from the active environment"},
			{Name: "status", Description: "View the status of railway project"},
			{Name: "unlink", Description: "Disconnects the current directory from a Railway project"},
			{Name: "up", Description: "Deploy to railway"},
			{Name: "path", Description: "Path to deploy to"},
			{Name: "variables", Description: "Work with environment variables"},
			{Name: "variable", Description: "The name of the variable you want to delete"},
			{Name: "get", Description: "Get variable value"},
			{Name: "set", Description: "Set variable value"},
			{Name: "value", Description: "Value of the variable"},
			{Name: "version", Description: "Get the version of railway's CLI"},
			{Name: "whoami", Description: "Get the logged in user"},
		},
		Options: []core.Option{
			{Name: "-e", Description: "Environment to delete from"},
			{Name: "--browserless", Description: "Login without opening a browser"},
			{Name: "-n", Description: "Output a specific number of lines"},
			{Name: "-s", Description: "Define specific service"},
			{Name: "-d", Description: "Detach from build logs"},
			{Name: "--help", Description: "Show help for railway"},
			{Name: "--version", Description: "Show railway version"},
			{Name: "--verbose", Description: "Enable verbose output"},
		},
	})
}
