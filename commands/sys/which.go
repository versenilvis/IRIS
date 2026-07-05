package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "which",
		Description: "Executable file",
		Options: []spec.Option{
			{Name: "-s", Description: "No output, return 0 if all the executables are found, 1 if not"},
			{Name: "-a", Description: "List all instances of executables found, instead of just the first"},
		},
	})
}
