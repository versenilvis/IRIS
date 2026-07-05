package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "eza",
		Description: "A modern replacement for ls",
		Options: []core.Option{
			{Name: "-?", Description: "Show list of command-line options"},
			{Name: "-v", Description: "Show version of eza"},
			{Name: "-1", Description: "Display one entry per line"},
			{Name: "-l", Description: "Display extended file metadata as a table"},
			{Name: "-G", Description: "Display entries as a grid (default)"},
			{Name: "-x", Description: "Sort the grid across, rather than downward"},
			{Name: "-R", Description: "Recurse into directories"},
			{Name: "-T", Description: "Recurse into directories as a tree"},
			{Name: "-X", Description: "Dereference symbolic links when displaying information"},
			{Name: "-F", Description: "Display type indicator by file names"},
			{Name: "--color", Description: "When to use terminal colours"},
			{Name: "--color-scale", Description: "Highlight levels of 'field' distinctly"},
			{Name: "--color-scale-mode", Description: "Use gradient or fixed colors in --color-scale"},
			{Name: "--icons", Description: "When to display icons"},
			{Name: "--no-quotes", Description: "Don't quote file names with spaces"},
			{Name: "--hyperlink", Description: "Display entries as hyperlinks"},
			{Name: "--absolute", Description: "Display entries with their absolute path"},
			{Name: "--follow-symlinks", Description: "Drill down into symbolic links that point to directories"},
			{Name: "-w", Description: "Set screen width in columns"},
			{Name: "-a", Description: "Show hidden and 'dot' files"},
			{Name: "-A", Description: "Equivalent to '--all'"},
			{Name: "-d", Description: "List directories like regular files"},
			{Name: "-D", Description: "List only directories"},
			{Name: "-f", Description: "List only files"},
			{Name: "--show-symlinks", Description: "Explicitly show symbolic links"},
			{Name: "--no-symlinks", Description: "Do not show symbolic links"},
			{Name: "-L", Description: "Limit the depth of recursion"},
			{Name: "-r", Description: "Reverse the sort order"},
			{Name: "-s", Description: "Which field to sort by"},
			{Name: "--group-directories-first", Description: "List directories before other files"},
			{Name: "--group-directories-last", Description: "List directories after other files"},
			{Name: "-I", Description: "Glob patterns (pipe-separated) of files to ignore"},
			{Name: "--git-ignore", Description: "Ignore files mentioned in '.gitignore'"},
			{Name: "-b", Description: "List file sizes with binary prefixes"},
			{Name: "-B", Description: "List file sizes in bytes, without any prefixes"},
			{Name: "-g", Description: "List each file's group"},
			{Name: "--smart-group", Description: "Only show group if it has a different name from owner"},
			{Name: "-h", Description: "Add a header row to each column"},
			{Name: "-H", Description: "List each file's number of hard links"},
			{Name: "-i", Description: "List each file's inode number"},
		},
	})
}
