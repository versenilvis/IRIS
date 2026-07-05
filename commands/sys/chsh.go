package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "chsh",
		Description: "Change your login shell",
		Options: []spec.Option{
			{Name: "-s", Description: "Specify login shell"},
			{Name: "-l", Description: "Print list of shells and exit"},
			{Name: "-u", Description: "Print help message and exit"},
			{Name: "-v", Description: "Print version and exit"},
		},
	})
}
