package view

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "cot",
		Description: "Command-line utility for CotEditor",
	})
}
