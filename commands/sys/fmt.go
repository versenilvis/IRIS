package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "fmt",
		Description: "Simple text formatter",
		Options: []core.Option{
			{Name: "-c", Description: "File(s) to format"},
		},
	})
}
