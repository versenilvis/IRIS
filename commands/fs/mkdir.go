package fs

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "mkdir",
		Description: "make directories",
		Generator:   spec.FileGenerator("/"),
		Options: []spec.Option{
			{Name: "-p", Description: "create parent dirs"},
			{Name: "-v", Description: "verbose"},
		},
	})
}
