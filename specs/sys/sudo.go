package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "sudo",
		Description: "Execute a command as the superuser or another user",
		Options: []core.Option{
			{Name: "-g", Description: "Run command as the specified group name or ID"},
			{Name: "-h", Description: "Display help message and exit"},
			{Name: "-u", Description: "Run command as specified user name or ID"},
		},
	})
}
