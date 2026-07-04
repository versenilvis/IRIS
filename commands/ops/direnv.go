package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "direnv",
		Description: "Help for direnv",
		Subcommands: []core.Subcommand{
			{Name: "allow", Description: "Grants direnv to load the given .envrc"},
			{Name: "deny", Description: "Revokes the authorization of a given .envrc"},
			{Name: "exec", Description: "Executes a command after loading the first .envrc found in DIR"},
			{Name: "fetchurl", Description: "Fetches a given URL into direnv's CAS"},
			{Name: "hook", Description: "Used to setup the shell hook"},
			{Name: "prune", Description: "Removes old allowed files"},
			{Name: "reload", Description: "Triggers an env reload"},
			{Name: "status", Description: "Prints some debug status information"},
			{Name: "stdlib", Description: "Displays the stdlib available in the .envrc execution context"},
			{Name: "show_dump", Description: "Show the data inside of a dump for debugging purposes"},
			{Name: "dotenv", Description: "Transforms a .env file to evaluatable `export KEY=PAIR` statements"},
			{Name: "dump", Description: "Used to export the inner bash state at the end of execution"},
			{Name: "FILE", Description: "Overwrites by dump data"},
			{Name: "export", Description: "Loads an .envrc and prints the diff in terms of exports"},
			{Name: "watch", Description: "Adds a path to the list that direnv watches for changes"},
			{Name: "watch-dir", Description: "Recursively adds a directory to the list that direnv watches for changes"},
			{Name: "watch-list", Description: "Pipe pairs of `mtime path` to stdin to build a list of files to watch"},
			{Name: "current", Description: "Reports whether direnv's view of a file is current (or stale)"},
		},
	})
}
