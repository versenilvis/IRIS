package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "okta",
		Description: "The Okta CLI is the easiest way to get started with Okta!",
		Subcommands: []spec.Subcommand{
			{Name: "register", Description: "Sign up for a new Okta account"},
			{Name: "login", Description: "Authorizes the Okta CLI tool"},
			{Name: "apps", Description: "Manage Okta apps"},
			{Name: "config", Description: "Show an Okta app's configuration"},
			{Name: "create", Description: "Create a new Okta app"},
			{Name: "delete", Description: "Deletes an Okta app"},
			{Name: "appIds", Description: "List of application IDs to be deleted"},
			{Name: "start", Description: "Creates an Okta Sample Application"},
			{Name: "name", Description: "The name of the sample app to create"},
			{Name: "help", Description: "Displays help information about the specified command"},
			{Name: "generate-completion", Description: "Generate bash/zsh completion script for Okta"},
		},
		Options: []spec.Option{
			{Name: "--company", Description: "Company/organization used when registering a new Okta account"},
			{Name: "--email", Description: "Email used when registering a new Okta account"},
			{Name: "--first-name", Description: "First name used when registering a new Okta account"},
			{Name: "--last-name", Description: "Last name used when registering a new Okta account"},
			{Name: "--app", Description: "The App ID"},
			{Name: "--app-name", Description: "Application name to be created, defaults to current directory name"},
			{Name: "--authorization-server-id", Description: "Okta Authorization Server Id"},
			{Name: "--config-file", Description: "Application config file"},
			{Name: "--redirect-uri", Description: "OIDC Redirect URI"},
			{Name: "-f", Description: "Deactivate and delete applications without confirmation"},
			{Name: "--help", Description: "Show help for Okta CLI"},
			{Name: "--version", Description: "Print version information"},
			{Name: "--verbose", Description: "Verbose logging"},
			{Name: "--batch", Description: "Batch mode, will not prompt for user input"},
		},
	})
}
