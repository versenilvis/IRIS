package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "electron",
		Description: "Build cross platform desktop apps with JavaScript, HTML and CSS",
		Options: []core.Option{
			{Name: "-i", Description: "Open a REPL to the main process"},
			{Name: "-r", Description: "Module to preload"},
			{Name: "-v", Description: "Print the version"},
			{Name: "-a", Description: "Print the Node ABI version"},
		},
	})
}
