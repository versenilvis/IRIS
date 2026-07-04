package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "man",
		Description: "${section} ${description}",
		Options: []core.Option{
			{Name: "-C", Description: "Specify the configuration file to use"},
			{Name: "-M", Description: "Specify the list of directories to search (colon separated)"},
			{Name: "-P", Description: "Specify the pager program"},
			{Name: "-B", Description: "Specify which browser to use for HTML files"},
			{Name: "-H", Description: "Specify a command that renders HTML files as text"},
			{Name: "-S", Description: "Specify a colon-separated list of manual sections to search"},
			{Name: "-a", Description: "Open every matching page instead of just the first"},
			{Name: "-c", Description: "Reformat the source page, even when an up-to-date cat-page exists"},
			{Name: "-d", Description: "Don't actually display the pages (dry run)"},
			{Name: "-D", Description: "Both display and print debugging info"},
			{Name: "-f", Description: "Equivalent to `whatis`"},
			{Name: "-F", Description: "Format only, do not display"},
			{Name: "-h", Description: "Print a help message and exit"},
			{Name: "-k", Description: "Equivalent to apropos"},
			{Name: "-K", Description: "Search for a given string in all pages"},
			{Name: "-m", Description: "Specify an alternate set of pages to search based on the system name given"},
			{Name: "-p", Description: "Specify the sequence of preprocessors to run before nroff or troff"},
			{Name: "-t", Description: "Use `/usr/bin/groff -Tps -mandoc -c` to format the page"},
			{Name: "-w", Description: "Print the location of files that would be displayed"},
			{Name: "-W", Description: "Print file locations, one per line"},
		},
	})
}
