package view

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "bat",
		Description: "A cat(1) clone with syntax highlighting and Git integration",
		Options: []spec.Option{
			{Name: "-A", Description: "Show non-printable characters"},
			{Name: "-p", Description: "Show plain style, no decorations"},
			{Name: "-l", Description: "Explicitly set the language for syntax highlighting"},
			{Name: "-H", Description: "Highlight the specified line ranges"},
			{Name: "--file-name", Description: "File(s)"},
			{Name: "-d", Description: "Show lines that have been added/removed/modified with respect to the Git index"},
			{Name: "--diff-context", Description: "Lines of context"},
			{Name: "--tabs", Description: "Set the tab width to T spaces. Use a width of 0 to pass tabs through directly"},
			{Name: "--wrap", Description: "Specify when to use colored output"},
			{Name: "--terminal-width", Description: "Explicitly set the width of the terminal instead of determining it automatically"},
			{Name: "-n", Description: "Show line numbers, no other decorations"},
			{Name: "--color", Description: "Specify when to use colored output"},
			{Name: "--italic-text", Description: "Specify when to use ANSI sequences for italic text in the output"},
			{Name: "--decorations", Description: "Specify when to use the decorations that have been specified via '--style'"},
			{Name: "-f", Description: "Alias for '--decorations=always --color=always'"},
			{Name: "--paging", Description: "Specify when to use the pager"},
			{Name: "--pager", Description: "Determine which pager is used"},
			{Name: "-m", Description: "Map a glob pattern to an existing syntax name"},
			{Name: "--ignored-suffix", Description: "Ignore extension"},
			{Name: "--theme", Description: "Set the theme for syntax highlighting"},
			{Name: "--list-themes", Description: "Display a list of supported themes for syntax highlighting"},
			{Name: "--style", Description: "Display a list of supported themes for syntax highlighting"},
			{Name: "-r", Description: "Only print the specified range of lines for each file"},
			{Name: "-L", Description: "Display a list of supported languages for syntax highlighting"},
			{Name: "-u", Description: "Show diagnostic information for bug reports"},
			{Name: "--acknowledgements", Description: "Show acknowledgements"},
			{Name: "-h", Description: "Print help message"},
			{Name: "-V", Description: "Show version information"},
		},
	})
}
