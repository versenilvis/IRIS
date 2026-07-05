package runner

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "php",
		Description: "Run the PHP interpreter",
	})
}
