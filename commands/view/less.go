package view

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "less",
		Description: "view file contents (scrollable)",
		Generator:   spec.FileGenerator(),
	})
}
