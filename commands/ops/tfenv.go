package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "tfenv",
		Description: "Version",
		Subcommands: []spec.Subcommand{
			{Name: "install", Description: "Install a specific version of Terraform"},
			{Name: "version", Description: "Possible Terraform Version"},
			{Name: "use", Description: "Switch to a version to use"},
			{Name: "uninstall", Description: "Uninstall a specific version of Terraform"},
			{Name: "list", Description: "List all installed versions"},
			{Name: "list-remote", Description: "List all installable versions"},
			{Name: "version-name", Description: "Print current version"},
			{Name: "init", Description: "Update environment to use tfenv correctly"},
		},
		Options: []spec.Option{
			{Name: "-v", Description: "View your current tfenv version"},
			{Name: "--help", Description: "View commands"},
		},
	})
}
