package view

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "file",
		Description: "determine file type",
		Generator:   spec.FileGenerator(),
	})
}
