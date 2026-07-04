package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "readlink",
		Description: "Display file status",
		Options: []core.Option{
			{Name: "-f", Description: "Do not force a newline to appear at the end of each piece of output"},
		},
	})
}
