package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "deta",
		Description: "Runtime",
		Subcommands: []core.Subcommand{
			{Name: "login", Description: "Trigger the login process for the Deta CLI"},
			{Name: "version", Description: "Print the Deta version"},
			{Name: "upgrade", Description: "Upgrade Deta CLI version"},
			{Name: "projects", Description: "List Deta projects"},
			{Name: "new", Description: "Create a new Deta Micro"},
			{Name: "path", Description: "Path to new directory for the micro"},
			{Name: "name", Description: "Name of the new micro"},
			{Name: "project", Description: "Name of the existing project"},
			{Name: "runtime", Description: "The selected runtime"},
			{Name: "deploy", Description: "Deploy new code to a Deta Micro"},
			{Name: "details", Description: "Get detailed information about a specific Deta micro"},
			{Name: "watch", Description: "Auto-deploy locally saved changes in real time to your Deta micro"},
			{Name: "auth", Description: "Change auth settings for a Deta Micro"},
			{Name: "disable", Description: "Disable HTTP Auth for a Deta Micro"},
			{Name: "enable", Description: "Enable HTTP Auth for a Deta Micro"},
			{Name: "create-api-key", Description: "Create an API key for a Deta Micro"},
			{Name: "description", Description: "The api-key description"},
			{Name: "outfile", Description: "The api-key output file"},
			{Name: "delete-api-key", Description: "Delete an API key for a Deta Micro"},
			{Name: "pull", Description: "Pull the latest deployed code of a Deta Micro to your local machine"},
			{Name: "clone", Description: "Clone a Deta Micro"},
			{Name: "update", Description: "Update a Deta Micro's name or environment variables"},
			{Name: "env", Description: "Path to env file"},
			{Name: "visor", Description: "Change the Visor settings for a Deta Micro"},
			{Name: "open", Description: "Open Micro's visor page in the browser"},
			{Name: "run", Description: "Run a Deta Micro from the CLI"},
			{Name: "action", Description: "The action to be performed on the micro. See docs for full examples and details"},
			{Name: "cron", Description: "Change cron settings for a Deta Micro"},
			{Name: "set", Description: "Set Deta Micro to run on a schedule"},
			{Name: "expression", Description: "The cron expression to be set"},
			{Name: "remove", Description: "Remove a schedule from a Deta Micro"},
		},
		Options: []core.Option{
			{Name: "-h", Description: "Show help for login"},
			{Name: "-v", Description: "Upgrade CLI to specific version"},
			{Name: "-n", Description: "Create a micro with Node (node14.x) runtime"},
			{Name: "-p", Description: "Create a micro with Python (python 3.9) runtime"},
			{Name: "--name", Description: "Set the name of the new micro"},
			{Name: "--project", Description: "Set the project under which the micro is created"},
			{Name: "--runtime", Description: "Create a micro with a specified runtime"},
			{Name: "-d", Description: "Set the api-key description"},
			{Name: "-o", Description: "Set the api-key output file"},
			{Name: "-f", Description: "Force the overwrite of existing files"},
			{Name: "-r", Description: "The new runtime of the micro"},
			{Name: "-e", Description: "The new env file of the micro"},
			{Name: "-l", Description: "Show the micro logs"},
		},
	})
}
