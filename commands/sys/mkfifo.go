package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "mkfifo",
		Description: "Make FIFOs (first-in, first-out)",
		Options: []spec.Option{
			{Name: "-m", Description: "FIFO(s) to create"},
		},
	})
}
