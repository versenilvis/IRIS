package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "unzip",
		Description: "Extract compressed files in a ZIP archive",
		Options: []core.Option{
			{Name: "-l", Description: "List the contents of a zip file without extracting"},
		},
	})
}
