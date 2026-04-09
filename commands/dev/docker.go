package dev

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "docker",
		Description: "container engine",
		Subcommands: []core.Subcommand{
			{
				Name: "ps",
				Description: "list containers",
				Options: []core.Option{
					{Name: "-a", Description: "show all"},
					{Name: "-q", Description: "only show IDs"},
				},
			},
			{
				Name: "build",
				Description: "build image",
				Options: []core.Option{
					{Name: "-t", Description: "tag name"},
					{Name: "-f", Description: "dockerfile path"},
					{Name: "--no-cache", Description: "no build cache"},
				},
			},
			{
				Name: "run",
				Description: "run container",
				Options: []core.Option{
					{Name: "-d", Description: "detached mode"},
					{Name: "-p", Description: "port mapping"},
					{Name: "-v", Description: "volume mount"},
					{Name: "--rm", Description: "auto remove"},
					{Name: "--name", Description: "container name"},
					{Name: "-it", Description: "interactive tty"},
					{Name: "-e", Description: "set env variable"},
				},
			},
			{
				Name: "pull",
				Description: "pull image",
			},
			{
				Name: "push",
				Description: "push image",
			},
			{
				Name: "exec",
				Description: "exec in container",
				Options: []core.Option{
					{Name: "-it", Description: "interactive tty"},
				},
			},
			{
				Name: "stop",
				Description: "stop container",
			},
			{
				Name: "rm",
				Description: "remove container",
				Options: []core.Option{
					{Name: "-f", Description: "force remove"},
				},
			},
			{
				Name: "rmi",
				Description: "remove image",
			},
			{
				Name: "images",
				Description: "list images",
			},
			{
				Name: "logs",
				Description: "view logs",
				Options: []core.Option{
					{Name: "-f", Description: "follow output"},
					{Name: "--tail", Description: "last n lines"},
				},
			},
			{
				Name: "compose",
				Description: "multi-container",
				Subcommands: []core.Subcommand{
					{Name: "up", Description: "start services"},
					{Name: "down", Description: "stop services"},
					{Name: "build", Description: "build services"},
					{Name: "logs", Description: "view logs"},
					{Name: "ps", Description: "list services"},
					{Name: "exec", Description: "execute command"},
					{Name: "restart", Description: "restart services"},
				},
			},
			{
				Name: "network",
				Description: "manage networks",
				Subcommands: []core.Subcommand{
					{Name: "ls", Description: "list networks"},
					{Name: "create", Description: "create network"},
					{Name: "rm", Description: "remove network"},
					{Name: "inspect", Description: "show details"},
				},
			},
			{
				Name: "volume",
				Description: "manage volumes",
				Subcommands: []core.Subcommand{
					{Name: "ls", Description: "list volumes"},
					{Name: "create", Description: "create volume"},
					{Name: "rm", Description: "remove volume"},
					{Name: "prune", Description: "remove unused"},
				},
			},
		},
	})
}
