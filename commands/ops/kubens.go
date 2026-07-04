package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "kubens",
		Description: "Switch between Kubernetes-namespaces",
		Options: []core.Option{
			{Name: "--help", Description: "Show help for kubens"},
			{Name: "--current", Description: "Show current namespace"},
		},
	})
}
