package view

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "cot",
		Description: "Command-line utility for CotEditor",
	})
}
