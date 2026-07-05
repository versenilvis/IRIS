package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "lima",
		Description: "Lima is an alias for",
		Options: []spec.Option{
			{Name: "-h", Description: "Help for lima"},
		},
	})
}
