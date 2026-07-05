package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "eleventy",
		Description: "Eleventy is a simpler static site generator",
	})
}
