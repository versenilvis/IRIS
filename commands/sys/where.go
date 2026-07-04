package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "where",
		Description: "For each name, indicate how it should be interpreted",
		Options: []core.Option{
			{Name: "-w", Description: "For each name, print 'name: word', where 'word' is the kind of command"},
			{Name: "-p", Description: "Do a path search for the name, even if it's an alias/function/builtin"},
			{Name: "-m", Description: "The arguments are taken as patterns (pattern characters must be quoted)"},
			{Name: "-s", Description: "If the pathname contains symlinks, print the symlink-free name as well"},
			{Name: "-S", Description: "Print intermediate symlinks and the resolved name"},
			{Name: "-x", Description: "Expand tabs when outputting shell function"},
		},
	})
}
