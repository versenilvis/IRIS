package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "build-storybook",
		Description: "Storybook build CLI tools",
		Options: []core.Option{
			{Name: "-o", Description: "Directory where to store built files"},
			{Name: "-w", Description: "Enables watch mode"},
			{Name: "--loglevel", Description: "Controls level of logging during build"},
		},
	})
}
