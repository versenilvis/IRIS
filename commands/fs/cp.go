package fs

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "cp",
		Description: "copy files and directories",
		Generator:   spec.FileGenerator(),
		Options: []spec.Option{
			{Name: "-r", Description: "recursive"},
			{Name: "-a", Description: "archive mode"},
			{Name: "-v", Description: "verbose"},
			{Name: "-i", Description: "interactive"},
			{Name: "-u", Description: "update (copy only if newer)"},
			{Name: "-p", Description: "preserve attributes"},
		},
	})
}
