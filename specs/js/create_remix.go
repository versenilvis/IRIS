package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "create-remix",
		Description: "Display help for command",
		Options: []core.Option{
			{Name: "-h", Description: "Display help for command"},
			{Name: "-v", Description: "Display version for command"},
		},
	})
}
