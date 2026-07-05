package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "kubectx",
		Description: "Switch between Kubernetes-contexts",
		Options: []spec.Option{
			{Name: "--help", Description: "Show help for kubectx"},
			{Name: "--current", Description: "Show current context"},
			{Name: "--unset", Description: "Unset the current context"},
			{Name: "-d", Description: "Delete context"},
		},
	})
}
