package text

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "unix2dos",
		Description: "Unix to DOS text file format convertor",
	})
}
