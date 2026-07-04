package view

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "vimr",
		Description: "VimR — Neovim GUI for macOS in Swift",
		Options: []core.Option{
			{Name: "-h", Description: "Show help"},
			{Name: "--dry-run", Description: "Just print the 'open' command"},
			{Name: "--cwd", Description: "Set the working directory"},
			{Name: "--line", Description: "Go to line"},
			{Name: "--wait", Description: "This command line tool will exit when the corresponding UI window is closed"},
			{Name: "--nvim", Description: "Open files in tabs in a new window"},
			{Name: "-s", Description: "Open files in separate windows"},
		},
	})
}
