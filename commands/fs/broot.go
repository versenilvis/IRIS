package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "broot",
		Description: "Show the last modified date of files and directories",
		Options: []core.Option{
			{Name: "--dates", Description: "Show the last modified date of files and directories"},
			{Name: "--no-dates", Description: "Don't show the last modified date"},
			{Name: "--only-folders", Description: "Only show folders"},
			{Name: "--no-only-folders", Description: "Show folders and files alike"},
			{Name: "--show-root-fs", Description: "Show filesystem info on top"},
			{Name: "--show-git-info", Description: "Show git statuses on files and stats of repository"},
			{Name: "--no-show-git-info", Description: "Don't show git statuses on files nor stats"},
			{Name: "--git-status", Description: "Only show files having an interesting git status, including hidden ones"},
			{Name: "--hidden", Description: "Show hidden files"},
			{Name: "--no-hidden", Description: "Don't show hidden files"},
			{Name: "--show-gitignored", Description: "Show files which should be ignored according to git"},
			{Name: "--no-show-gitignored", Description: "Don't show gitignored files"},
			{Name: "--permissions", Description: "Show permissions with owner and group"},
			{Name: "--no-permissions", Description: "Don't show permissions"},
			{Name: "--sizes", Description: "Don't show sizes"},
			{Name: "--sort-by-count", Description: "Sort by count (only show one level of the tree)"},
			{Name: "--sort-by-date", Description: "Sort by date (only show one level of the tree)"},
			{Name: "--sort-by-size", Description: "Sort by size (only show one level of the tree)"},
			{Name: "--whale-spotting", Description: "Sort by size, show ignored and hidden files"},
			{Name: "--no-sort", Description: "Don't sort"},
			{Name: "--trim-root", Description: "Install or reinstall the br shell function"},
			{Name: "--get-root", Description: "Ask for the current root of the remote broot"},
			{Name: "--color", Description: "Prints a help page, with more or less the same content as this man page"},
			{Name: "--version", Description: "Prints the version of broot"},
			{Name: "--outcmd", Description: "Where to write a command if broot produces one"},
			{Name: "--cmd", Description: "Semicolon separated commands to execute on start of broot"},
			{Name: "--conf", Description: "Semicolon separated paths to specific config files"},
			{Name: "--height", Description: "Where to write the produced path, if any"},
			{Name: "--set-install-state", Description: "Listen for commands"},
			{Name: "--send", Description: "Send commands to a remote broot then quits"},
		},
	})
}
