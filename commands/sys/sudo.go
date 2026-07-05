package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "sudo",
		Description: "Execute a command as the superuser or another user",
		Options: []spec.Option{
			{Name: "-g", Description: "Run command as the specified group name or ID"},
			{Name: "-h", Description: "Display help message and exit"},
			{Name: "-u", Description: "Run command as specified user name or ID"},
		},
	})
}
