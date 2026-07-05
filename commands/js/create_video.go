package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "create-video",
		Description: "CLI used to create remotion video project",
	})
}
