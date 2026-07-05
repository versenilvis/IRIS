package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "create-vite",
		Description: "Create a new project powered by Vite",
	})
}
