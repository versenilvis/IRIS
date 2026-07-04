package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "shell-config",
		Description: "Display help for command",
		Subcommands: []core.Subcommand{
			{Name: "install", Description: "Install MacOS setup with Multi-Selection"},
			{Name: "shell", Description: "Setup a shell configuration with a robust set of tools and architecture"},
			{Name: "update", Description: "Update the CLI to the latest version"},
			{Name: "assets", Description: "Configure your shell assets, such as `gitprofile` and `npmrc`"},
			{Name: "external", Description: "Install and manage your external shell"},
			{Name: "list", Description: "Show list of all exists externals shells"},
			{Name: "file_path", Description: "File to set as external shell"},
			{Name: "external_name", Description: "Name to register under the registry"},
			{Name: "delete", Description: "Delete external shell"},
			{Name: "help", Description: "Display help for command"},
		},
		Options: []core.Option{
			{Name: "-h", Description: "Display help for command"},
			{Name: "-t", Description: "Select update version"},
			{Name: "-m", Description: "Daemon check for update notification. When specified is true"},
			{Name: "-V", Description: "Output the version number"},
		},
	})
}
