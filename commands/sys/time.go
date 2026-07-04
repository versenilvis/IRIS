package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "time",
		Description: "Time how long a command takes!",
	})
}
