package view

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "emacs",
		Description: "An extensible, customizable, free/libre text editor - and more",
	})
}
