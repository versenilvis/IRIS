package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "pwd",
		Description: "Return working directory name",
		Options: []core.Option{
			{Name: "-L", Description: "Display the logical current working directory"},
			{Name: "-P", Description: "Display the physical current working directory"},
		},
	})
}
