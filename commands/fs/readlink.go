package fs

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "readlink",
		Description: "Display file status",
		Options: []spec.Option{
			{Name: "-f", Description: "Do not force a newline to appear at the end of each piece of output"},
		},
	})
}
