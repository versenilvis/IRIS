package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "fold",
		Description: "Fold long lines for finite width output device",
		Options: []core.Option{
			{Name: "-b", Description: "File(s) to fold"},
		},
	})
}
