package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "pwd",
		Description: "Return working directory name",
		Options: []spec.Option{
			{Name: "-L", Description: "Display the logical current working directory"},
			{Name: "-P", Description: "Display the physical current working directory"},
		},
	})
}
