package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "mkfifo",
		Description: "Make FIFOs (first-in, first-out)",
		Options: []core.Option{
			{Name: "-m", Description: "FIFO(s) to create"},
		},
	})
}
