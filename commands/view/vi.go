package view

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "vi",
		Description: "Print help message for vi and exit",
		Options: []spec.Option{
			{Name: "-h", Description: "Print help message for vi and exit"},
		},
	})
}
