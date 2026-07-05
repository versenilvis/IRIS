package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "atuin",
		Description: "Magical shell history",
		Subcommands: []spec.Subcommand{
			{Name: "history", Description: "Manipulate shell history"},
			{Name: "start", Description: "Begins a new command in the history"},
			{Name: "end", Description: "Finishes a new command in the history (adds time, exit code)"},
			{Name: "list", Description: "List all items in history"},
			{Name: "last", Description: "Get the last command ran"},
			{Name: "import", Description: "Import shell history from file"},
			{Name: "auto", Description: "Import history for the current shell"},
			{Name: "stats", Description: "Calculate statistics for your history"},
			{Name: "search", Description: "Interactive history search"},
			{Name: "sync", Description: "Sync with the configured server"},
			{Name: "login", Description: "Login to the configured server"},
			{Name: "logout", Description: "Log out"},
			{Name: "register", Description: "Register with the configured server"},
			{Name: "key", Description: "Print the encryption key for transfer to another machine"},
			{Name: "server", Description: "Start an atuin server"},
			{Name: "init", Description: "Output shell setup"},
			{Name: "uuid", Description: "Generate a UUID"},
			{Name: "gen-completions", Description: "Generate shell completions"},
		},
		Options: []spec.Option{
			{Name: "--exit", Description: "List all items in history"},
			{Name: "--cwd", Description: "Show only the text of the command"},
			{Name: "--format", Description: "Get the last command ran"},
			{Name: "--human", Description: "Show only the text of the command"},
			{Name: "--count", Description: "How many top commands to list"},
			{Name: "--exclude-cwd", Description: "Exclude directory from results"},
			{Name: "--exclude-exit", Description: "Exclude results with this exit code"},
			{Name: "--before", Description: "Only include results added before this date"},
			{Name: "--after", Description: "Only include results after this date"},
			{Name: "--limit", Description: "How many entries to return at most"},
			{Name: "--offset", Description: "Offset from the start of the results"},
			{Name: "--interactive", Description: "Open interactive search UI"},
			{Name: "--filter-mode", Description: "Allow overriding filter mode over config"},
			{Name: "--search-mode", Description: "Allow overriding search mode over config"},
			{Name: "--cmd-only", Description: "Show only the text of the command"},
			{Name: "--delete", Description: "Delete anything matching this query. Will not print out the match"},
			{Name: "--delete-it-all", Description: "Delete EVERYTHING!"},
			{Name: "--reverse", Description: "Reverse the order of results, oldest first"},
			{Name: "--force", Description: "Force re-download everything"},
			{Name: "--username", Description: "The encryption key for your account"},
			{Name: "--base64", Description: "Switch to base64 output of the key"},
			{Name: "--key", Description: "Start an atuin server"},
			{Name: "--host", Description: "The host address to bind"},
			{Name: "--port", Description: "The port to bind"},
			{Name: "--disable-ctrl-r", Description: "Disable the binding of CTRL-R to atuin"},
			{Name: "--disable-up-arrow", Description: "Disable the binding of the Up Arrow key to atuin"},
			{Name: "--shell", Description: "Set the shell for generating completions"},
			{Name: "--out-dir", Description: "Set the output directory"},
			{Name: "--help", Description: "Print help"},
			{Name: "--version", Description: "Print version"},
		},
	})
}
