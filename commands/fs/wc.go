package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "wc",
		Description: "word, line, character count",
		Generator:   core.FileGenerator(),
		Options: []core.Option{
			{Name: "-l", Description: "count lines"},
			{Name: "-w", Description: "count words"},
			{Name: "-c", Description: "count bytes"},
		},
	})
}
