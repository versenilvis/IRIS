package fs

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "tree",
		Description: "Display directories as trees (with optional color/HTML output)",
		Options: []spec.Option{
			{Name: "-a", Description: "All files are listed"},
			{Name: "-d", Description: "List directories only"},
			{Name: "-l", Description: "Follow symbolic links like directories"},
			{Name: "-f", Description: "Print the full path prefix for each file"},
			{Name: "-x", Description: "Stay on current filesystem only"},
			{Name: "-L", Description: "Descend only level directories deep"},
			{Name: "-R", Description: "Rerun tree when max dir level reached"},
			{Name: "-P", Description: "List only those files that match the pattern given"},
			{Name: "-I", Description: "Do not list files that match the given pattern"},
			{Name: "--ignore-case", Description: "Ignore case when pattern matching"},
			{Name: "--matchdirs", Description: "Include directory names in -P pattern matching"},
			{Name: "--noreport", Description: "Turn off file/directory count at end of tree listing"},
			{Name: "--charset", Description: "Use charset X for terminal/HTML and indentation line output"},
			{Name: "--filelimit", Description: "Do not descend dirs with more than # files in them"},
			{Name: "--timefmt", Description: "Print and format time according to the format <f>"},
			{Name: "-o", Description: "Output to file instead of stdout"},
			{Name: "-q", Description: "Print non-printable characters as '?'"},
			{Name: "-N", Description: "Print non-printable characters as is"},
			{Name: "-Q", Description: "Quote filenames with double quotes"},
			{Name: "-p", Description: "Print the protections for each file"},
			{Name: "-u", Description: "Displays file owner or UID number"},
			{Name: "-g", Description: "Displays file group owner or GID number"},
			{Name: "-s", Description: "Print the size in bytes of each file"},
			{Name: "-h", Description: "Print the size in a more human readable way"},
			{Name: "--si", Description: "Like -h but use SI units (powers of 1000) instead"},
			{Name: "--du", Description: "Appends '/', '=', '*', '@', '|' or '>' as per ls -F"},
			{Name: "--inodes", Description: "Print inode number of each file"},
			{Name: "--device", Description: "Print device ID number to which each file belongs"},
			{Name: "-v", Description: "Sort files alphanumerically by version"},
			{Name: "-t", Description: "Sort files by last modification time"},
			{Name: "-c", Description: "Sort files by last status change time"},
			{Name: "-U", Description: "Leave files unsorted"},
			{Name: "-r", Description: "Reverse the order of the sort"},
			{Name: "--dirsfirst", Description: "List directories before files (-U disables)"},
			{Name: "--sort", Description: "Select sort"},
			{Name: "-i", Description: "Don't print indentation lines"},
			{Name: "-A", Description: "Print ANSI lines graphic indentation lines"},
			{Name: "-S", Description: "Print with CP437 (console) graphics indentation lines"},
			{Name: "-n", Description: "Turn colorization off always (-C overrides)"},
			{Name: "-C", Description: "Turn colorization on always"},
		},
	})
}
