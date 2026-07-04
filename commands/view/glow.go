package view

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "glow",
		Description: "Render markdown on the CLI, with pizzazz!",
		Subcommands: []core.Subcommand{
			{Name: "config", Description: "Edit the glow config file"},
			{Name: "help", Description: "Help about any command"},
			{Name: "stash", Description: "Manage your stash of markdown files"},
		},
		Options: []core.Option{
			{Name: "-a", Description: "Show system files and directories (TUI-mode only)"},
			{Name: "--config", Description: "Config file"},
			{Name: "-h", Description: "Help for glow"},
			{Name: "-l", Description: "Show local files only; no network (TUI-mode only)"},
			{Name: "-p", Description: "Display with pager"},
			{Name: "-s", Description: "Style name or JSON path (default 'auto')"},
			{Name: "-v", Description: "Version for glow"},
			{Name: "-w", Description: "Word-wrap at width"},
			{Name: "-m", Description: "Memo/note for stashing"},
		},
	})
}
