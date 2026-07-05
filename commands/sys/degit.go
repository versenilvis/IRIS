package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "degit",
		Description: "Straightforward project scaffolding",
		Options: []spec.Option{
			{Name: "--help", Description: "Print help"},
			{Name: "-f", Description: "Overwrite existing files"},
			{Name: "-c", Description: "Use a cache"},
			{Name: "-v", Description: "Be verbose?"},
			{Name: "-m", Description: "Clone mode"},
		},
	})
}
