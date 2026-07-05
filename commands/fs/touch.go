package fs

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "touch",
		Description: "create or update file timestamp",
		Generator:   spec.FileGenerator(),
	})
}
