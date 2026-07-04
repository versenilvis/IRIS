package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "nrm",
		Description: "Use the right package manage - remove",
		Options: []core.Option{
			{Name: "-g", Description: "Package will be removed from your `devDependencies`"},
			{Name: "-P", Description: "Remove package from your `peerDependencies`"},
			{Name: "-O", Description: "Remove package from your `optionalDependencies`"},
			{Name: "--frozen", Description: "Don't generate a lockfile and fail if an update is needed"},
			{Name: "-h", Description: "Output usage information"},
		},
	})
}
