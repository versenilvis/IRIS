package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "find",
		Description: "Walk a file hierarchy",
		Options: []core.Option{
			{Name: "-E", Description: "Permit find to be safely used in conjunction with xargs"},
			{Name: "-d", Description: "Cause find to perform a depth-first traversal"},
			{Name: "-f", Description: "Specify a file hierarch for find to traverse"},
			{Name: "-s", Description: "Cause find to traverse the file hierarchies in lexicographical order"},
		},
	})
}
