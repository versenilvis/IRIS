package text

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "unix2dos",
		Description: "Unix to DOS text file format convertor",
	})
}
