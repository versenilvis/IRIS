package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "trap",
		Description: "Prints all defined signal handlers",
		Options: []core.Option{
			{Name: "--print", Description: "Prints all defined signal handlers"},
			{Name: "--help", Description: "Displays help about using this command"},
		},
	})
}
