package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "degit",
		Description: "Straightforward project scaffolding",
		Options: []core.Option{
			{Name: "--help", Description: "Print help"},
			{Name: "-f", Description: "Overwrite existing files"},
			{Name: "-c", Description: "Use a cache"},
			{Name: "-v", Description: "Be verbose?"},
			{Name: "-m", Description: "Clone mode"},
		},
	})
}
