package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "tac",
		Description: "Concatenate and print files in reverse",
		Options: []spec.Option{
			{Name: "--help", Description: "Display this help and exit"},
			{Name: "--before", Description: "Attach the separator before instead of after"},
			{Name: "--regex", Description: "Interpret the separator as a regular expression"},
			{Name: "--separator", Description: "Use STRING as the separator instead of newline"},
			{Name: "--version", Description: "Output version information and exit"},
		},
	})
}
