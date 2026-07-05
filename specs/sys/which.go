package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "which",
		Description: "Executable file",
		Options: []core.Option{
			{Name: "-s", Description: "No output, return 0 if all the executables are found, 1 if not"},
			{Name: "-a", Description: "List all instances of executables found, instead of just the first"},
		},
	})
}
