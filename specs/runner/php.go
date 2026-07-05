package runner

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "php",
		Description: "Run the PHP interpreter",
	})
}
