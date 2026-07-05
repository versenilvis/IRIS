package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "mknod",
		Description: "Create device special file",
		Subcommands: []spec.Subcommand{
			{Name: "c", Description: "Create (c)haracter device"},
			{Name: "b", Description: "Create (b)lock device"},
		},
		Options: []spec.Option{
			{Name: "-F", Description: "Format"},
		},
	})
}
