package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "fnm",
		Description: "Fast and simple Node.js version manager",
		Subcommands: []core.Subcommand{
			{Name: "install", Description: "Install a new Node.js version"},
			{Name: "uninstall", Description: "Uninstall a Node.js version"},
			{Name: "use", Description: "Change Node.js version"},
			{Name: "exec", Description: "Run a command within fnm context"},
			{Name: "current", Description: "Print the current Node.js version"},
			{Name: "list", Description: "List all locally installed Node.js versions"},
			{Name: "list-remote", Description: "List all remote Node.js versions"},
			{Name: "alias", Description: "Alias a version to a common name"},
			{Name: "name", Description: "Alias name"},
			{Name: "unalias", Description: "Deletes the alias named <name>"},
			{Name: "requested-alias", Description: "Alias name"},
			{Name: "completions", Description: "Print shell completions to stdout"},
			{Name: "env", Description: "Print and set up required environment variables for fnm"},
			{Name: "help", Description: "Prints the help page or the help of the given subcommand(s)"},
		},
		Options: []core.Option{
			{Name: "--lts", Description: "Install latest LTS"},
			{Name: "--install-if-missing", Description: "Install the version if it isn't installed yet"},
			{Name: "--help", Description: "Prints help information"},
			{Name: "--version", Description: "Prints version information"},
			{Name: "--arch", Description: "The root directory of fnm installations"},
			{Name: "--log-level", Description: "The log level of fnm commands"},
			{Name: "--node-dist-mirror", Description: "Mirror of https://nodejs.org/dist"},
			{Name: "--version-file-strategy", Description: "Strategy for how to resolve the Node version"},
			{Name: "--using", Description: "Either an explicit version, or a filename with the version written in it"},
			{Name: "--use-on-cd", Description: "Print the script to change Node versions every directory change"},
		},
	})
}
