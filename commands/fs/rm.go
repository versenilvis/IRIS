package fs

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "rm",
		Description: "remove files or directories",
		Generator:   spec.FileGenerator(),
		Options: []spec.Option{
			{Name: "-r", Description: "recursive"},
			{Name: "-f", Description: "force"},
			{Name: "-i", Description: "interactive"},
			{Name: "-v", Description: "verbose"},
			{Name: "-d", Description: "remove empty directories"},
		},
	})
}
