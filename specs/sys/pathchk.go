package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "pathchk",
		Description: "Check pathnames for POSIX portability",
		Options: []core.Option{
			{Name: "-p", Description: "Pathname(s) to check"},
		},
	})
}
