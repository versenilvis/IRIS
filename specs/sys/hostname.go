package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "hostname",
		Description: "Set or print name of current host system",
		Options: []core.Option{
			{Name: "-f", Description: "Include domain information in the printed name"},
			{Name: "-s", Description: "Trim off any domain information from the printed name"},
			{Name: "-d", Description: "Only print domain information"},
		},
	})
}
