package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "defaults",
		Description: "Global domain",
		Subcommands: []core.Subcommand{
			{Name: "read", Description: "Shows defaults"},
			{Name: "write", Description: "Writes key for domain"},
			{Name: "delete", Description: "Deletes domain or key in domain"},
			{Name: "rename", Description: "Renames old_key to new_key"},
			{Name: "domains", Description: "Lists all domains"},
			{Name: "find", Description: "Lists all entries containing word"},
			{Name: "word", Description: "The word to search for"},
			{Name: "help", Description: "Show help text"},
			{Name: "read-type", Description: "Shows the type for the given domain, key"},
		},
		Options: []core.Option{
			{Name: "-globalDomain", Description: "Global domain"},
			{Name: "-app", Description: "Application name"},
			{Name: "-string", Description: "Command line interface to a user's defaults"},
		},
	})
}
