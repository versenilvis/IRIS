package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "tfenv",
		Description: "Version",
		Subcommands: []core.Subcommand{
			{Name: "install", Description: "Install a specific version of Terraform"},
			{Name: "version", Description: "Possible Terraform Version"},
			{Name: "use", Description: "Switch to a version to use"},
			{Name: "uninstall", Description: "Uninstall a specific version of Terraform"},
			{Name: "list", Description: "List all installed versions"},
			{Name: "list-remote", Description: "List all installable versions"},
			{Name: "version-name", Description: "Print current version"},
			{Name: "init", Description: "Update environment to use tfenv correctly"},
		},
		Options: []core.Option{
			{Name: "-v", Description: "View your current tfenv version"},
			{Name: "--help", Description: "View commands"},
		},
	})
}
