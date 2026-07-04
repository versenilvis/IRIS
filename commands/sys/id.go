package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "id",
		Description: "Display the full name of the user",
		Options: []core.Option{
			{Name: "-A", Description: "Display the full name of the user"},
			{Name: "-G", Description: "Display the MAC label of the current process"},
			{Name: "-P", Description: "Display the id as a password file entry"},
			{Name: "-g", Description: "Display the effective group ID as a number"},
			{Name: "-n", Description: "Make the output human-readable"},
			{Name: "-u", Description: "Display the effective user ID as a number"},
		},
	})
}
