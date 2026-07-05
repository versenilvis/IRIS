package cc

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "xcode-select",
		Description: "Active developer directory for Xcode tools",
		Options: []core.Option{
			{Name: "-h", Description: "Help message"},
			{Name: "-p", Description: "Display path to active developer directory"},
			{Name: "-s", Description: "Set path to active developer directory"},
			{Name: "--install", Description: "Install Xcode command line tools"},
			{Name: "-v", Description: "Display version"},
			{Name: "-r", Description: "Reset to the default CLT path"},
		},
	})
}
