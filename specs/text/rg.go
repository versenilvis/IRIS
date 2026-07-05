package text

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "rg",
		Description: "ripgrep (fast search)",
		Generator:   core.FileGenerator(),
		Options: []core.Option{
			{Name: "-i", Description: "case insensitive"},
			{Name: "-s", Description: "case sensitive"},
			{Name: "-n", Description: "show line numbers"},
			{Name: "-N", Description: "hide line numbers"},
			{Name: "-l", Description: "filenames only"},
			{Name: "-L", Description: "follow symlinks"},
			{Name: "-c", Description: "count matches"},
			{Name: "-w", Description: "word boundary match"},
			{Name: "-x", Description: "whole line match"},
			{Name: "-v", Description: "invert match"},
			{Name: "-o", Description: "only matching text"},
			{Name: "-A", Description: "lines after match"},
			{Name: "-B", Description: "lines before match"},
			{Name: "-C", Description: "lines around match"},
			{Name: "-m", Description: "max matches per file"},
			{Name: "-e", Description: "additional pattern"},
			{Name: "-f", Description: "patterns from file"},
			{Name: "--type", Description: "file type filter (e.g. go, py, js)"},
			{Name: "--type-not", Description: "exclude file type"},
			{Name: "--glob", Description: "include glob pattern"},
			{Name: "--iglob", Description: "include glob (case insensitive)"},
			{Name: "--hidden", Description: "search hidden files"},
			{Name: "--no-ignore", Description: "don't use .gitignore"},
			{Name: "--no-ignore-vcs", Description: "don't use VCS ignore"},
			{Name: "-g", Description: "include/exclude glob"},
			{Name: "--json", Description: "output as JSON"},
			{Name: "--stats", Description: "print match statistics"},
			{Name: "--color", Description: "colorize output (auto/always/never)"},
			{Name: "--no-heading", Description: "no file name headings"},
			{Name: "-H", Description: "always show filename"},
			{Name: "--max-depth", Description: "max directory depth"},
			{Name: "-p", Description: "smart heading/colors (for piping)"},
			{Name: "-U", Description: "multiline matching"},
			{Name: "--multiline", Description: "enable multiline"},
			{Name: "-F", Description: "treat pattern as literal string"},
			{Name: "--fixed-strings", Description: "literal string match"},
			{Name: "-z", Description: "search compressed files"},
			{Name: "--encoding", Description: "file encoding"},
			{Name: "--trim", Description: "trim whitespace"},
		},
	})

	core.Register(&core.Spec{
		Name:        "fd",
		Description: "fast find alternative",
		Generator:   core.FileGenerator("/"),
		Options: []core.Option{
			{Name: "-t", Description: "type (f/d/l/x/e/s)"},
			{Name: "-e", Description: "file extension"},
			{Name: "-H", Description: "include hidden"},
			{Name: "-I", Description: "no gitignore"},
			{Name: "-d", Description: "max depth"},
			{Name: "-x", Description: "execute command"},
			{Name: "-X", Description: "exec batch"},
			{Name: "-l", Description: "long listing"},
			{Name: "-L", Description: "follow symlinks"},
			{Name: "-p", Description: "full path match"},
			{Name: "-g", Description: "glob pattern"},
			{Name: "-c", Description: "color mode"},
			{Name: "-0", Description: "null-separated output"},
			{Name: "-j", Description: "number of threads"},
			{Name: "--changed-within", Description: "recently changed"},
			{Name: "--changed-before", Description: "changed before date"},
		},
	})
}
