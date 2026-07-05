package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "exec",
		Description: "Replace the current shell with a program",
	})
}
