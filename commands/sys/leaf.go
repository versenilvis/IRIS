package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "leaf",
		Description: "Create and interact with your leaf projects",
		Subcommands: []spec.Subcommand{
			{Name: "completion", Description: "Dump the shell completion script"},
			{Name: "create", Description: "[init] Create a new Leaf PHP project"},
			{Name: "project name", Description: "Name of the project"},
			{Name: "deploy", Description: "[publish] Deploy your leaf project"},
			{Name: "help", Description: "Display help for a command"},
			{Name: "install", Description: "Add a new package to your leaf app"},
			{Name: "package name", Description: "Name of the package"},
			{Name: "interact", Description: "Interact with your application"},
			{Name: "list", Description: "List commands"},
			{Name: "run", Description: "Run a script in your composer.json"},
			{Name: "command name", Description: "Name of the command"},
			{Name: "serve", Description: "Run your Leaf app"},
			{Name: "port number", Description: "The port number to run the server on"},
			{Name: "test", Description: "Test your leaf application through leaf alchemy"},
			{Name: "test:setup", Description: "Add tests to your application"},
			{Name: "uninstall", Description: "Uninstall a package"},
			{Name: "update", Description: "Update leaf cli to the latest version"},
		},
		Options: []spec.Option{
			{Name: "-h", Description: "Do not output any message"},
			{Name: "-V", Description: "Display this application version"},
			{Name: "--ansi", Description: "Force ANSI output"},
			{Name: "--no-ansi", Description: "Disable ANSI output"},
			{Name: "-n", Description: "Do not ask any interactive question"},
			{Name: "-v", Description: "Dump the shell completion script"},
			{Name: "--port", Description: "The port number to run the server on"},
			{Name: "--watch", Description: "Watch for changes and restart the server"},
		},
	})
}
