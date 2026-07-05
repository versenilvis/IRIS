package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "aliases",
		Description: "Prints help information",
		Subcommands: []spec.Subcommand{
			{Name: "add", Description: "Add an alias via the cli"},
			{Name: "name", Description: "The name of the alias"},
			{Name: "command", Description: "The command you want to run"},
			{Name: "clone", Description: "Clone external aliases"},
			{Name: "username", Description: "The username of the aliases you want to clone"},
			{Name: "repo_url", Description: "The git repo url of the aliases (defaults to github/<username>/dot-aliases)"},
			{Name: "directories", Description: "List all directories initialized with aliases"},
			{Name: "exec", Description: "Execute an alias for a given directory"},
			{Name: "directory", Description: "Directory where the alias is defined"},
			{Name: "help", Description: "Prints help information"},
			{Name: "init", Description: "Initialize a directory for aliases"},
			{Name: "list", Description: "List the aliases available"},
			{Name: "pull", Description: "Pull a cloned user's aliases"},
			{Name: "rehash", Description: "Update the aliases"},
			{Name: "remove", Description: "Remove an alias via the cli"},
			{Name: "users", Description: "List the users"},
			{Name: "disable", Description: "Disable a user's aliases"},
			{Name: "enable", Description: "Enable a user's aliases"},
			{Name: "move", Description: "Move a user up or down the prioritization list"},
			{Name: "use", Description: "Assign a user to the top of the priority list"},
		},
		Options: []spec.Option{
			{Name: "--help", Description: "Prints help information"},
			{Name: "--version", Description: "Prints version information"},
			{Name: "-E", Description: "Whether to enable the user if they are not currently enabled"},
			{Name: "-g", Description: "Returns the global initialization for the app"},
			{Name: "-u", Description: "Initialize aliases for a specific user"},
			{Name: "-l", Description: "List only local aliases"},
			{Name: "-d", Description: "List aliases for a specific directory"},
		},
	})
}
