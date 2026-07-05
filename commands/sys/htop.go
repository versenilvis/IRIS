package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "htop",
		Description: "Improved top (interactive process viewer)",
		Options: []spec.Option{
			{Name: "--help", Description: "Show help for htop"},
			{Name: "--no-color", Description: "Use a monochrome color scheme"},
			{Name: "--delay", Description: "Delay between updates, in tenths of sec"},
			{Name: "--filter", Description: "Filter commands"},
			{Name: "--highlight-changes", Description: "Highlight new and old processes"},
			{Name: "--no-mouse", Description: "Disable the mouse"},
			{Name: "--pid", Description: "Show only the given PIDs"},
			{Name: "--sort-key", Description: "Sort by COLUMN in list view"},
			{Name: "--tree", Description: "Show the tree view"},
			{Name: "--user", Description: "Show only processes for a given user (or $USER)"},
			{Name: "--no-unicode", Description: "Do not use unicode but plain ASCII"},
			{Name: "--version", Description: "Print version info"},
		},
	})
}
