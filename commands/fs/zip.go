package fs

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "zip",
		Description: "Package and compress (archive) files into zip file",
		Options: []spec.Option{
			{Name: "-r", Description: "Package and compress a directory and its contents, recursively"},
			{Name: "-e", Description: "Archive a directory and its contents with the highest level [9] of compression"},
		},
	})
}
