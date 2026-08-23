package fs

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "rmdir",
		Description: "Remove directories",
		Generator:   spec.FileGenerator("/"),
		Options: []spec.Option{
			{Name: "-p", Description: "Remove each directory of path"},
		},
	})
}
