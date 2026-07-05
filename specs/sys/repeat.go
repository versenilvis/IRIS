package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "repeat",
		Description: "Interpret the result as a number and repeat the commands this many times",
	})
}
