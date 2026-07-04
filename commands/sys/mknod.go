package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "mknod",
		Description: "Create device special file",
		Subcommands: []core.Subcommand{
			{Name: "c", Description: "Create (c)haracter device"},
			{Name: "b", Description: "Create (b)lock device"},
		},
		Options: []core.Option{
			{Name: "-F", Description: "Format"},
		},
	})
}
