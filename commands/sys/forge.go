package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "forge",
		Description: "A command line interface for managing Atlassian-hosted apps",
		Subcommands: []spec.Subcommand{
			{Name: "autocomplete", Description: "Configures autocomplete for the Forge CLI"},
			{Name: "create", Description: "Create an app"},
			{Name: "deploy", Description: "Deploy your app to an environment"},
			{Name: "feedback", Description: "Let us know what you think about Forge"},
			{Name: "install", Description: "Manage app installations"},
			{Name: "list", Description: "List app installations"},
			{Name: "help", Description: "Display help for command"},
			{Name: "lint", Description: "Check the source files for common errors"},
			{Name: "login", Description: "Log in to your Atlassian account"},
			{Name: "logout", Description: "Log out of your Atlassian account"},
			{Name: "logs", Description: "View app logs"},
			{Name: "providers", Description: "Manage external providers"},
			{Name: "configure", Description: "Configure provider credentials"},
			{Name: "register", Description: "Register an app you didn't create so you can run commands for it"},
			{Name: "settings", Description: "Manage Forge CLI settings"},
			{Name: "set", Description: "Update Forge CLI setting (choices: usage-analytics)"},
			{Name: "uninstall", Description: "Uninstall the app from an Atlassian site"},
			{Name: "variables", Description: "Manage app environment variables"},
			{Name: "unset", Description: "Remove an environment variable"},
			{Name: "webtrigger", Description: "Get a web trigger URL"},
			{Name: "whoami", Description: "Display the account information of the logged in user"},
		},
		Options: []spec.Option{
			{Name: "--verbose", Description: "Enable verbose mode"},
			{Name: "-h", Description: "Display help for command"},
			{Name: "-t", Description: "Specify the template to use"},
			{Name: "-d", Description: "Specify the directory to create (uses the template name by default)"},
			{Name: "-f", Description: "Disable pre-deployment checks"},
			{Name: "-e", Description: "Run the command without input prompts"},
			{Name: "-p", Description: "Product (Jira, Confluence, Compass)"},
			{Name: "--upgrade", Description: "Upgrade an existing installation"},
			{Name: "--confirm-scopes", Description: "Skip confirmation of scopes for the app before installing or upgrading the app"},
			{Name: "--non-interactive", Description: "Run the command without input prompts"},
			{Name: "--fix", Description: "Attempt to automatically fix any issues encountered"},
			{Name: "-u", Description: "Specify the email to use"},
			{Name: "-n", Description: "Number of invocations to return"},
			{Name: "-s", Description: "Group logs by invocation ID"},
			{Name: "--encrypt", Description: "Encrypt variable"},
			{Name: "-V", Description: "Output the version number"},
		},
	})
}
