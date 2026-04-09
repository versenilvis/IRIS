package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "ln",
		Description: "create links",
		Generator:   core.FileGenerator(),
		Options: []core.Option{
			{Name: "-s", Description: "symbolic link"},
			{Name: "-f", Description: "force"},
		},
	})
}
