package view

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "emacs",
		Description: "An extensible, customizable, free/libre text editor - and more",
	})
}
