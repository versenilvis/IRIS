package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "create-video",
		Description: "CLI used to create remotion video project",
	})
}
