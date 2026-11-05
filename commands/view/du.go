package view

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "du",
		Description: "estimate file space usage",
		Options: []spec.Option{
			{Name: "-h", Description: "human readable"},
			{Name: "-s", Description: "summarize"},
			{Name: "-a", Description: "all files"},
		},
	})
}
