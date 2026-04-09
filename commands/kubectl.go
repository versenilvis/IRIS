package commands

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "kubectl",
		Description: "kubernetes cli",
		Subcommands: []core.Subcommand{
			{Name: "get", Description: "display resources",
				Subcommands: []core.Subcommand{
					{Name: "pods", Description: "list pods"},
					{Name: "services", Description: "list services"},
					{Name: "deployments", Description: "list deployments"},
					{Name: "nodes", Description: "list nodes"},
					{Name: "namespaces", Description: "list namespaces"},
				},
			},
			{Name: "apply", Description: "apply config",
				Options: []core.Option{
					{Name: "-f", Description: "filename"},
				},
			},
			{Name: "describe", Description: "show details"},
			{Name: "logs", Description: "view pod logs",
				Options: []core.Option{
					{Name: "-f", Description: "follow logs"},
				},
			},
			{Name: "delete", Description: "delete resource"},
			{Name: "exec", Description: "execute in pod"},
			{Name: "port-forward", Description: "forward ports"},
			{Name: "scale", Description: "scale resource"},
		},
	})
}
