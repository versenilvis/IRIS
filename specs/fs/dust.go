package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "dust",
		Description: "Like du but more intuitive",
		Options: []core.Option{
			{Name: "--help", Description: "Show help for dust"},
			{Name: "--version", Description: "Print version information"},
			{Name: "--no-percent-bars", Description: "No percent bars or percentages will be displayed"},
			{Name: "--no-colors", Description: "No colors will be printed (Useful for commands like: watch)"},
			{Name: "--depth", Description: "Depth to show"},
			{Name: "--filter", Description: "Directory 'size' is number of child files/dirs not disk size"},
			{Name: "--si", Description: "Print sizes in powers of 1000 (e.g., 1.1G)"},
			{Name: "--ignore_hidden", Description: "Do not display hidden files"},
			{Name: "--number-of-lines", Description: "Number of lines of output to show. (Default is terminal_height - 10)"},
			{Name: "--full-paths", Description: "Subdirectories will not have their path shortened"},
			{Name: "--reverse", Description: "Print tree upside down (biggest highest)"},
			{Name: "--apparent-size", Description: "Use file length instead of blocks"},
			{Name: "--skip-total", Description: "No total row will be displayed"},
			{Name: "--file_types", Description: "Show only these file types"},
			{Name: "--invert-filter", Description: "Specify width of output overriding the auto detection of terminal width"},
			{Name: "--limit-filesystem", Description: "Exclude any file or directory with this name"},
			{Name: "--min-size", Description: "Minimum size file to include in output"},
		},
	})
}
