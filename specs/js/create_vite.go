package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "create-vite",
		Description: "Create a new project powered by Vite",
	})
}
