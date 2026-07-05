package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "hop",
		Description: "Interact with Hop in your terminal",
		Subcommands: []spec.Subcommand{
			{Name: "auth", Description: "Authenticate with Hop"},
			{Name: "help", Description: "Prints this message or the help of the given subcommand(s)"},
			{Name: "login", Description: "Login to Hop"},
			{Name: "logout", Description: "Logout the current user"},
			{Name: "deploy", Description: "Deploy a new container"},
			{Name: "ignite", Description: "Interact with Ignite containers"},
			{Name: "ls", Description: "List all deployments"},
			{Name: "rm", Description: "Delete a deployment"},
			{Name: "link", Description: "Link an existing deployment to a hopfile"},
			{Name: "projects", Description: "Interact with projects"},
			{Name: "info", Description: "Get information about a project"},
			{Name: "new", Description: "Create a new project"},
			{Name: "switch", Description: "Switch to a different project"},
			{Name: "whoami", Description: "Get information about the current user"},
		},
		Options: []spec.Option{
			{Name: "--type", Description: "Type of the container, defaults to `ephemeral`"},
			{Name: "--containers", Description: "Number of containers to use, defaults to 1 if `scaling` is manual"},
			{Name: "--cpu", Description: "The number of CPUs to use, defaults to 1"},
			{Name: "--env", Description: "Environment variables to set, in the form of KEY=VALUE"},
			{Name: "--name", Description: "Name of the deployment, defaults to the directory name"},
			{Name: "--project", Description: "Namespace or ID of the project to use"},
			{Name: "--ram", Description: "Amount of RAM to use, defaults to 512MB"},
			{Name: "--scaling", Description: "Scaling strategy, defaults to `manual`"},
			{Name: "--help", Description: "Prints help information"},
			{Name: "--version", Description: "Prints version information"},
		},
	})
}
