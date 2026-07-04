package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "exa",
		Description: "A modern replacement for ls",
	})
}
