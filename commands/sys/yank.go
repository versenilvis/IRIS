package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "yank",
		Description: "Yank terminal output to clipboard",
		Options: []spec.Option{
			{Name: "-i", Description: "Ignore case differences between pattern and the input"},
			{Name: "-l", Description: "Use the default delimiters except for space"},
			{Name: "-x", Description: "Use alternate screen"},
			{Name: "-v", Description: "Print the version"},
			{Name: "-d", Description: "All input characters not present in delim will be recognized as fields"},
			{Name: "-g", Description: "Use pattern to recognize fields"},
		},
	})
}
