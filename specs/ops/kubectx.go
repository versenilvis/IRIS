package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "kubectx",
		Description: "Switch between Kubernetes-contexts",
		Options: []core.Option{
			{Name: "--help", Description: "Show help for kubectx"},
			{Name: "--current", Description: "Show current context"},
			{Name: "--unset", Description: "Unset the current context"},
			{Name: "-d", Description: "Delete context"},
		},
	})
}
