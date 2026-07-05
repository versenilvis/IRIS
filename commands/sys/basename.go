package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "basename",
		Description: "Return filename portion of pathname",
		Options: []spec.Option{
			{Name: "-a", Description: "Treat every argument as a string"},
			{Name: "-s", Description: "Suffix to remove from string"},
		},
	})
}
