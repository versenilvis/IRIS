package text

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "seq",
		Description: "Print sequences of numbers. (Defaults to increments of 1)",
		Options: []spec.Option{
			{Name: "-w", Description: "Equalize the widths of all numbers by padding with zeros as necessary"},
			{Name: "-s", Description: "String separator between numbers. Default is newline"},
			{Name: "-f", Description: "Use a printf(3) style format to print each number"},
		},
	})
}
