package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "ansible-config",
		Description: "View ansible configuration",
		Subcommands: []spec.Subcommand{
			{Name: "list", Description: "List and output available configs"},
			{Name: "CONFIG_FILE", Description: "Path to configuration file"},
			{Name: "TYPE", Description: "Plugin type"},
			{Name: "args", Description: "Specific plugin to target, requires type of plugin to be set"},
			{Name: "dump", Description: "Shows the current settings, merges ansible.cfg if specified"},
			{Name: "view", Description: "Displays the current config file"},
			{Name: "init", Description: "Initializes a new config file (to stdout)"},
			{Name: "FORMAT", Description: "Output format"},
		},
		Options: []spec.Option{
			{Name: "--help", Description: "Show help and exit"},
			{Name: "--verbose", Description: "Verbose mode (-vvv for more, -vvvv to enable connection debugging)"},
			{Name: "-v", Description: "Verbose mode (-vvv for more, -vvvv to enable connection debugging)"},
			{Name: "--config", Description: "Path to configuration file, defaults to first file found in precedence"},
			{Name: "--type", Description: "Filter down to a specific plugin type"},
			{Name: "--only-changed", Description: "Only show configurations that have changed from the default"},
			{Name: "--disabled", Description: "Prefixes all entries with a comment character to disable them"},
			{Name: "--format", Description: "Output format for init"},
			{Name: "--version", Description: "Show help and exit"},
		},
	})
}
