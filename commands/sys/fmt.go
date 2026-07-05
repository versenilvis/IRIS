package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "fmt",
		Description: "Simple text formatter",
		Options: []spec.Option{
			{Name: "-c", Description: "File(s) to format"},
		},
	})
}
