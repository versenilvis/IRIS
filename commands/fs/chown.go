package fs

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "chown",
		Description: "change file owner",
		Generator:   spec.FileGenerator(),
		Options: []spec.Option{
			{Name: "-R", Description: "recursive"},
		},
	})
}
