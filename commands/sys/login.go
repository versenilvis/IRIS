package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "login",
		Description: "Begin session on the system",
		Options: []spec.Option{
			{Name: "-p", Description: "Preserve environment"},
			{Name: "-r", Description: "Perform autologin protocol for rlogin"},
			{Name: "-h", Description: "Specify host"},
			{Name: "-f", Description: "Don't authenticate user, user is preauthenticated"},
		},
	})
}
