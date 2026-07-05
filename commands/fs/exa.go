package fs

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "exa",
		Description: "A modern replacement for ls",
	})
}
