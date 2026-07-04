package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "ag",
		Description: "Recursively search for PATTERN in PATH. Like grep or ack, but faster",
		Options: []core.Option{
			{Name: "--ackmate", Description: "Set thread affinity (if platform supports it)"},
			{Name: "--noaffinity", Description: "Don't set thread affinity (if platform supports it)"},
			{Name: "-a", Description: "Print lines after match"},
			{Name: "-B", Description: "Print lines before match"},
			{Name: "--nobreak", Description: "Print a newline between matches in different files. Enabled by default"},
			{Name: "-c", Description: "Print color codes in results"},
			{Name: "--nocolor", Description: "Don't print color codes in results"},
			{Name: "--color-line-number", Description: "Color codes for line numbers. Default is 1;33"},
			{Name: "--color-match", Description: "Color codes for result match numbers. Default is 30;43"},
			{Name: "--color-path", Description: "Color codes for path names. Default is 1;32"},
			{Name: "--column", Description: "Print column numbers in results"},
			{Name: "-C", Description: "Print lines before and after matches"},
			{Name: "-D", Description: "Search up to NUM directories deep, -1 for unlimited"},
			{Name: "--filename", Description: "Print file names"},
			{Name: "--nofilename", Description: "Don't print file names"},
			{Name: "-f", Description: "Follow symlinks"},
			{Name: "--nofollow", Description: "Don't follow symlinks"},
			{Name: "-F", Description: "Alias for --literal for compatibility with grep"},
			{Name: "--group", Description: "Print filenames matching PATTERN"},
			{Name: "-G", Description: "Only search files whose names match PATTERN"},
			{Name: "-H", Description: "Print filenames above matching contents"},
			{Name: "--noheading", Description: "Don't print filenames above matching contents"},
			{Name: "--hidden", Description: "Search hidden files. This option obeys ignored files"},
			{Name: "--ignore", Description: "The pattern to look for"},
			{Name: "--ignore-dir", Description: "Alias for --ignore for compatibility with ack"},
			{Name: "-i", Description: "Match case-insensitively"},
			{Name: "-l", Description: "Only print the names of files that don't contain matches"},
			{Name: "--list-file-types", Description: "See FILE TYPES below"},
			{Name: "-m", Description: "Skip the rest of a file after NUM matches. Default is 0, which never skips"},
			{Name: "--mmap", Description: "Match regexes across newlines"},
			{Name: "--nomultiline", Description: "Don't match regexes across newlines"},
			{Name: "-n", Description: "Don't recurse into directories"},
			{Name: "--numbers", Description: "Print line numbers"},
			{Name: "--nonumbers", Description: "Don't print line numbers"},
			{Name: "-o", Description: "Print only the matching part of the lines"},
			{Name: "--one-device", Description: "Provide a path to a specific .ignore file"},
			{Name: "--pager", Description: "The pager"},
			{Name: "--parallel", Description: "Print matches on very long lines (> 2k characters by default)"},
			{Name: "--passthrough", Description: "When searching a stream, print all lines even if they don't match"},
			{Name: "-Q", Description: "Do not parse PATTERN as a regular expression. Try to match it literally"},
		},
	})
}
