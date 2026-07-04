package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "dtm",
		Description: "Plugin",
		Subcommands: []core.Subcommand{
			{Name: "apply", Description: "Create or update DevOps tools according to DevStream configuration file"},
			{Name: "completion", Description: "Generate the autocompletion script for dtm for the specified shell"},
			{Name: "bash", Description: "Generate autocompletion script for bash"},
			{Name: "fish", Description: "Generate autocompletion script for fish"},
			{Name: "powershell", Description: "Generate autocompletion script for powershell"},
			{Name: "zsh", Description: "Generate autocompletion script for zsh"},
			{Name: "delete", Description: "Delete DevOps tools according to DevStream configuration file"},
			{Name: "develop", Description: "Develop is used for develop a new plugin"},
			{Name: "create-plugin", Description: "Create a new plugin"},
			{Name: "validate-plugin", Description: "Validate a plugin"},
			{Name: "init", Description: "Download needed plugins according to the config file"},
			{Name: "list", Description: "This command only supports listing plugins now"},
			{Name: "plugins", Description: "List all plugins"},
			{Name: "show", Description: "Show is used to print plugins' configuration templates or status"},
			{Name: "config", Description: "Show configuration information"},
			{Name: "status", Description: "Show status information"},
			{Name: "upgrade", Description: "Upgrade dtm to the latest release version"},
			{Name: "verify", Description: "Verify DevOps tools according to DevStream config file and state"},
			{Name: "version", Description: "Print the version number of DevStream"},
			{Name: "help", Description: "Help about any command"},
		},
		Options: []core.Option{
			{Name: "--config-file", Description: "Config file"},
			{Name: "--plugin-dir", Description: "Plugins directory"},
			{Name: "--yes", Description: "Apply directly without confirmation"},
			{Name: "--force", Description: "Force delete by config"},
			{Name: "--name", Description: "Specify name of the plugin to be created"},
			{Name: "--all", Description: "Validate all plugins"},
			{Name: "--arch", Description: "Download plugins for specific arch"},
			{Name: "--download-only", Description: "Download plugins only"},
			{Name: "--os", Description: "Download plugins for specific os"},
			{Name: "--plugins", Description: "The plugins to be downloaded"},
			{Name: "--filter", Description: "Filter plugin by regex"},
			{Name: "--plugin", Description: "Specify name with the plugin"},
			{Name: "--template", Description: "Print a template config, e.g. quickstart/gitops/"},
			{Name: "--id", Description: "Specify id with the plugin instance"},
			{Name: "--debug", Description: "Debug level log"},
			{Name: "--help", Description: "Display help"},
		},
	})
}
