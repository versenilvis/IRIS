package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "passwd",
		Description: "Modify a user",
		Options: []core.Option{
			{Name: "-i", Description: "Specify where the password update should be applied"},
			{Name: "-l", Description: "The location of the chosen directory system"},
			{Name: "-u", Description: "Specify the user name to use when authenticating to the directory node"},
		},
	})
}
