package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "eleventy",
		Description: "Eleventy is a simpler static site generator",
	})
}
