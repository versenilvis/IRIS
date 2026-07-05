package view

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "wc",
		Description: "word, line, character count",
		Generator:   spec.FileGenerator(),
		Options: []spec.Option{
			{Name: "-l", Description: "count lines"},
			{Name: "-w", Description: "count words"},
			{Name: "-c", Description: "count bytes"},
		},
	})
}
