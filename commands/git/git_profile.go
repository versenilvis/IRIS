package git

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "git-profile",
		Description: "Use profile",
		Subcommands: []core.Subcommand{
			{Name: "use", Description: "Use a profile"},
			{Name: "profile", Description: "Profile you want to apply in this repository"},
		},
		Options: []core.Option{
			{Name: "--help", Description: "Help for git-profile script"},
		},
	})
}
