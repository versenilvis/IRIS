package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "chsh",
		Description: "Change your login shell",
		Options: []core.Option{
			{Name: "-s", Description: "Specify login shell"},
			{Name: "-l", Description: "Print list of shells and exit"},
			{Name: "-u", Description: "Print help message and exit"},
			{Name: "-v", Description: "Print version and exit"},
		},
	})
}
