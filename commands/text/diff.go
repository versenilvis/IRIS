package text

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "diff",
		Description: "Similar, but format ${name} input groups with GFTM",
		Options: []core.Option{
			{Name: "-i", Description: "Ignore case differences in file contents"},
			{Name: "--ignore-file-name-case", Description: "Ignore case when comparing file names"},
			{Name: "--no-ignore-file-name-case", Description: "Consider case when comparing file names"},
			{Name: "-E", Description: "Ignore changes due to tab expansion"},
			{Name: "-b", Description: "Ignore changes in the amount of white space"},
			{Name: "-w", Description: "Ignore all white space"},
			{Name: "-B", Description: "Ignore changes whose lines are all blank"},
			{Name: "-I", Description: "Ignore changes whose lines all match RE"},
			{Name: "--strip-trailing-cr", Description: "Strip trailing carriage return on input"},
			{Name: "-a", Description: "Treat all files as text"},
			{Name: "-c", Description: "Output NUM lines of copied context"},
			{Name: "-u", Description: "Output NUM lines of unified context"},
			{Name: "--label", Description: "Use LABEL instead of file name"},
			{Name: "-p", Description: "Show which C function each change is in"},
			{Name: "-F", Description: "Show the most recent line matching RE"},
			{Name: "-q", Description: "Output only whether files differ"},
			{Name: "-e", Description: "Output an ed script"},
			{Name: "--normal", Description: "Output a normal diff"},
			{Name: "-n", Description: "Output an RCS format diff"},
			{Name: "-y", Description: "Output in two columns"},
			{Name: "-W", Description: "Output at most NUM (default 130) print columns"},
			{Name: "--left-column", Description: "Output only the left column of common lines"},
			{Name: "--suppress-common-lines", Description: "Do not output common lines"},
			{Name: "-D", Description: "Output merged file to show `#ifdef NAME' diffs"},
			{Name: "-l", Description: "Pass the output through `pr' to paginate it"},
			{Name: "-t", Description: "Expand tabs to spaces in output"},
			{Name: "-T", Description: "Make tabs line up by prepending a tab"},
			{Name: "-r", Description: "Recursively compare any subdirectories found"},
			{Name: "-N", Description: "Treat absent files as empty"},
			{Name: "--unidirectional-new-file", Description: "Treat absent first files as empty"},
			{Name: "-s", Description: "Report when two files are the same"},
			{Name: "-x", Description: "Exclude files that match PAT"},
			{Name: "-X", Description: "Exclude files that match any pattern in FILE"},
			{Name: "-S", Description: "Start with FILE when comparing directories"},
			{Name: "--from-file", Description: "Compare FILE1 to all operands. FILE1 can be a directory"},
			{Name: "--to-file", Description: "Compare all operands to FILE2. FILE2 can be a directory"},
			{Name: "--horizon-lines", Description: "Keep NUM lines of the common prefix and suffix"},
			{Name: "-d", Description: "Try hard to find a smaller set of changes"},
			{Name: "--speed-large-files", Description: "Assume large files and many scattered small changes"},
			{Name: "-v", Description: "Output version info"},
		},
	})
}
