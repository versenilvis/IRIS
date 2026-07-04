package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "login",
		Description: "Begin session on the system",
		Options: []core.Option{
			{Name: "-p", Description: "Preserve environment"},
			{Name: "-r", Description: "Perform autologin protocol for rlogin"},
			{Name: "-h", Description: "Specify host"},
			{Name: "-f", Description: "Don't authenticate user, user is preauthenticated"},
		},
	})
}
