package fs

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "unzip",
		Description: "Extract compressed files in a ZIP archive",
		Options: []spec.Option{
			{Name: "-l", Description: "List the contents of a zip file without extracting"},
		},
	})
}
