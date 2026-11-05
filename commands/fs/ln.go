package fs

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "ln",
		Description: "create links",
		Options: []spec.Option{
			{Name: "-s", Description: "symbolic link"},
			{Name: "-f", Description: "force"},
		},
	})
}
