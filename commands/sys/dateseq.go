package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "dateseq",
		Description: "Print help and exit",
		Options: []spec.Option{
			{Name: "--help", Description: "Print help and exit"},
			{Name: "--version", Description: "Print version and exit"},
			{Name: "--quiet", Description: "Suppress message about date/time and duration parse errors and fix-ups"},
			{Name: "-f", Description: "Date/time"},
		},
	})
}
