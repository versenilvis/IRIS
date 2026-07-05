package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "echo",
		Description: "Environment Variable",
		Options: []spec.Option{
			{Name: "-n", Description: "Do not print the trailing newline character"},
			{Name: "-e", Description: "Interpret escape sequences"},
			{Name: "-E", Description: "Disable escape sequences"},
		},
	})
}
