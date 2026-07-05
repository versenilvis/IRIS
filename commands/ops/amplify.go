package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "amplify",
		Description: "Environment",
		Subcommands: []spec.Subcommand{
			{Name: "push", Description: "Provisions cloud resources with the latest local developments"},
			{Name: "delete", Description: "Deletes all of the resources tied to the project from the cloud"},
			{Name: "add", Description: "Adds a resource for an Amplify category in your local backend"},
			{Name: "update", Description: "Update resource for an Amplify category in your local backend"},
			{Name: "remove", Description: "Removes a resource for an Amplify category in your local backend"},
			{Name: "upgrade", Description: "Download and install the latest version of the Amplify CLI"},
			{Name: "migrate", Description: "Migrates GraphQL schemas to the latest GraphQL transformer version"},
			{Name: "override", Description: "Generates overrides file to apply custom modifications to CloudFormation"},
			{Name: "mock", Description: "Run mock server for testing categories locally"},
			{Name: "console", Description: "Opens the web console for the selected cloud resource"},
			{Name: "logout", Description: "If using temporary cloud provider credentials, this logs out of the account"},
			{Name: "env", Description: "Display all commands available for new Amplify project"},
			{Name: "pull", Description: "Pulls the current env from the cloud"},
			{Name: "checkout", Description: "Switches to selected environment"},
			{Name: "env-name", Description: "Env name"},
			{Name: "list", Description: "Displays a list of all the environments"},
			{Name: "get", Description: "Displays the environment details"},
			{Name: "import", Description: "Imports an environment"},
		},
		Options: []spec.Option{
			{Name: "-y", Description: "Answer all question as 'Yes'"},
			{Name: "-v", Description: "Deletes all of the resources tied to the project from the cloud"},
			{Name: "--restore", Description: "Overwrite your local changes"},
			{Name: "--details", Description: "See more details"},
			{Name: "--json", Description: "Format the output"},
			{Name: "--name", Description: "Mandatory flag"},
			{Name: "--config", Description: "Specify provider configs"},
			{Name: "--awsInfo", Description: "Specify AWS configs"},
		},
	})
}
