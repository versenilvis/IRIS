package view

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "tail",
		Description: "output last lines of file",
		Generator:   spec.FileGenerator(),
		Options: []spec.Option{
			{Name: "-n", Description: "number of lines"},
			{Name: "-f", Description: "follow (live updates)"},
		},
	})
}
