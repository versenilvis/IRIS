package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "lsd",
		Description: "An ls command with a lot of pretty colors and some other stuff",
		Options: []core.Option{
			{Name: "-1", Description: "Display one entry per line"},
			{Name: "-A", Description: "Do not list implied . and"},
			{Name: "-a", Description: "Do not ignore entries starting with"},
			{Name: "-d", Description: "Append indicator (one of */=>@|) at the end of the file names"},
			{Name: "-h", Description: "For ls compatibility purposes ONLY, currently set by default"},
			{Name: "-i", Description: "Display the index number of each file"},
			{Name: "-L", Description: "Display the extended file metadata as a table"},
			{Name: "-R", Description: "Recurse into directories"},
			{Name: "-r", Description: "Reverse the order of the sort"},
			{Name: "-S", Description: "Sort by size"},
			{Name: "-t", Description: "Sort by time modified"},
			{Name: "-v", Description: "Natural sort of (version) numbers within text"},
			{Name: "--classic", Description: "Enable classic mode (no colors or icons)"},
			{Name: "-X", Description: "Sort by file extension"},
			{Name: "--help", Description: "Prints help information"},
			{Name: "--ignore-config", Description: "Ignore the configuration file"},
			{Name: "--no-symlink", Description: "Do not display symlink target"},
			{Name: "--total-size", Description: "Display the total size of directories"},
			{Name: "--tree", Description: "Recurse into directories and present the result as a tree"},
			{Name: "-V", Description: "Prints version information"},
			{Name: "--blocks", Description: "Specify the blocks that will be displayed and in what order"},
			{Name: "--color", Description: "When to use terminal colours"},
			{Name: "--date", Description: "How to display date"},
			{Name: "--depth", Description: "Stop recursing into directories after reaching depth"},
			{Name: "--group-dirs", Description: "Sort the directories then the files"},
			{Name: "--icon", Description: "When to print the icons"},
			{Name: "--icon-theme", Description: "Whether to use fancy or unicode icons"},
			{Name: "--ignore-glob", Description: "How to display size"},
			{Name: "--sort", Description: "Sort by WORD instead of name"},
		},
	})
}
