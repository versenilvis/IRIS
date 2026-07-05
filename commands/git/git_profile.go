package git

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "git-profile",
		Description: "Use profile",
		Subcommands: []spec.Subcommand{
			{Name: "use", Description: "Use a profile"},
			{Name: "profile", Description: "Profile you want to apply in this repository"},
		},
		Options: []spec.Option{
			{Name: "--help", Description: "Help for git-profile script"},
		},
	})
}
