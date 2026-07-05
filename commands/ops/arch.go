package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "arch",
		Description: "32-bit intel",
		Options: []spec.Option{
			{Name: "-arch", Description: "Print architecture type or run select architecture"},
			{Name: "-32", Description: "Add the native 32-bit architecture to the list of architectures"},
			{Name: "-64", Description: "Add the native 64-bit architecture to the list of architectures"},
			{Name: "-c", Description: "Clear the environment that will be passed to the command"},
			{Name: "-d", Description: "Delete the named environment variable from the command's environment"},
			{Name: "-e", Description: "Assign the given value to the variable in the command's environment"},
			{Name: "-h", Description: "Print a help message and exit"},
		},
	})
}
