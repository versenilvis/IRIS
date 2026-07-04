package view

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "zed",
		Description: "A lightning-fast, collaborative code editor written in Rust",
		Options: []core.Option{
			{Name: "-h", Description: "Print help information"},
			{Name: "-v", Description: "Print Zed's version and the app path"},
			{Name: "-w", Description: "Wait for all of the given paths to be closed before exiting"},
			{Name: "-b", Description: "Custom Zed.app path"},
		},
	})
}
