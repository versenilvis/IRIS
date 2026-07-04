package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "btop",
		Description: "Beautifuler htop (interactive process viewer)",
		Options: []core.Option{
			{Name: "--help", Description: "Shows help for btop"},
			{Name: "--low-color", Description: "Disables truecolor, converts 24-bit colors to 256-color"},
			{Name: "--tty_on", Description: "Forces ON tty mode, max 16 colors and tty friendly graph symbol"},
			{Name: "--tty_off", Description: "Forces OFF tty mode"},
			{Name: "--preset", Description: "Start with preset"},
			{Name: "--utf-force", Description: "Force start even if no UTF-8 locale was detected"},
			{Name: "--debug", Description: "Shows the version of btop"},
		},
	})
}
