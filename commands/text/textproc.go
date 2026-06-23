package text

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	// awk
	core.Register(&core.Spec{
		Name:        "awk",
		Description: "pattern-directed scanning",
		Generator:   core.FileGenerator(),
		Options: []core.Option{
			{Name: "-F", Description: "field separator"},
			{Name: "-v", Description: "assign variable"},
			{Name: "-f", Description: "read program from file"},
			{Name: "-i", Description: "in-place edit (gawk)"},
			{Name: "--posix", Description: "POSIX compat mode"},
			{Name: "-W", Description: "compatibility options"},
		},
	})

	// gawk is also common
	core.Register(&core.Spec{
		Name:        "gawk",
		Description: "GNU awk",
		Generator:   core.FileGenerator(),
		Options: []core.Option{
			{Name: "-F", Description: "field separator"},
			{Name: "-v", Description: "assign variable"},
			{Name: "-f", Description: "read program from file"},
			{Name: "-i", Description: "in-place edit"},
			{Name: "--sandbox", Description: "sandbox mode"},
			{Name: "--profile", Description: "profiling output"},
		},
	})

	// sed
	core.Register(&core.Spec{
		Name:        "sed",
		Description: "stream editor",
		Generator:   core.FileGenerator(),
		Options: []core.Option{
			{Name: "-e", Description: "add expression"},
			{Name: "-f", Description: "read script from file"},
			{Name: "-i", Description: "in-place edit"},
			{Name: "-n", Description: "suppress default output"},
			{Name: "-E", Description: "extended regex"},
			{Name: "-r", Description: "extended regex (GNU)"},
			{Name: "-z", Description: "null-delimited lines"},
			{Name: "--sandbox", Description: "sandbox mode"},
		},
	})

	// xargs
	core.Register(&core.Spec{
		Name:        "xargs",
		Description: "build and run commands from stdin",
		Options: []core.Option{
			{Name: "-I", Description: "replace string (e.g. -I{})"},
			{Name: "-n", Description: "max args per command"},
			{Name: "-P", Description: "parallel jobs"},
			{Name: "-0", Description: "null-delimited input"},
			{Name: "-d", Description: "custom delimiter"},
			{Name: "-t", Description: "print command before executing"},
			{Name: "-p", Description: "prompt before executing"},
			{Name: "-r", Description: "no run if empty input"},
			{Name: "--no-run-if-empty", Description: "don't run if no input"},
		},
	})

	// tr
	core.Register(&core.Spec{
		Name:        "tr",
		Description: "translate or delete characters",
		Options: []core.Option{
			{Name: "-d", Description: "delete characters"},
			{Name: "-s", Description: "squeeze repeated chars"},
			{Name: "-c", Description: "complement set"},
		},
	})

	// sort
	core.Register(&core.Spec{
		Name:        "sort",
		Description: "sort lines of text",
		Generator:   core.FileGenerator(),
		Options: []core.Option{
			{Name: "-n", Description: "numeric sort"},
			{Name: "-r", Description: "reverse sort"},
			{Name: "-k", Description: "sort by key field"},
			{Name: "-t", Description: "field separator"},
			{Name: "-u", Description: "unique (remove duplicates)"},
			{Name: "-f", Description: "ignore case"},
			{Name: "-h", Description: "human-readable sort"},
			{Name: "-V", Description: "version sort"},
			{Name: "-o", Description: "output file"},
			{Name: "--parallel", Description: "parallel sort"},
		},
	})

	// uniq
	core.Register(&core.Spec{
		Name:        "uniq",
		Description: "filter adjacent duplicate lines",
		Generator:   core.FileGenerator(),
		Options: []core.Option{
			{Name: "-c", Description: "prefix count"},
			{Name: "-d", Description: "only duplicates"},
			{Name: "-u", Description: "only unique"},
			{Name: "-i", Description: "ignore case"},
			{Name: "-f", Description: "skip n fields"},
			{Name: "-s", Description: "skip n chars"},
		},
	})

	// cut
	core.Register(&core.Spec{
		Name:        "cut",
		Description: "extract columns from lines",
		Generator:   core.FileGenerator(),
		Options: []core.Option{
			{Name: "-d", Description: "delimiter"},
			{Name: "-f", Description: "field numbers (e.g. 1,3 or 1-3)"},
			{Name: "-c", Description: "character positions"},
			{Name: "--output-delimiter", Description: "output delimiter"},
		},
	})

	// tee
	core.Register(&core.Spec{
		Name:        "tee",
		Description: "read stdin, write to stdout and files",
		Generator:   core.FileGenerator(),
		Options: []core.Option{
			{Name: "-a", Description: "append to files"},
			{Name: "-i", Description: "ignore interrupts"},
		},
	})
}
