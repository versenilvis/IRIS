package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "mv",
		Description: "move (rename) files",
		Generator:   core.FileGenerator(),
		Options: []core.Option{
			{Name: "-i", Description: "interactive"},
			{Name: "-f", Description: "force (overwrite)"},
			{Name: "-n", Description: "no clobber (don't overwrite)"},
			{Name: "-v", Description: "verbose"},
		},
	})
}
