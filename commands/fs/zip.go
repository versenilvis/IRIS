package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "zip",
		Description: "Package and compress (archive) files into zip file",
		Options: []core.Option{
			{Name: "-r", Description: "Package and compress a directory and its contents, recursively"},
			{Name: "-e", Description: "Archive a directory and its contents with the highest level [9] of compression"},
		},
	})
}
