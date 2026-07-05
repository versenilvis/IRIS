package view

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "cat",
		Description: "concatenate and print",
		Generator:   spec.FileGenerator(),
		Options: []spec.Option{
			{Name: "-n", Description: "number lines"},
			{Name: "-b", Description: "number non-blank"},
		},
	})
}
