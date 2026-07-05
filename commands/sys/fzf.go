package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "fzf",
		Description: "A general-purpose command-line fuzzy finder",
		Options: []spec.Option{
			{Name: "-x", Description: "Enables extended-search mode"},
			{Name: "-e", Description: "Enables Exact-match"},
			{Name: "--algo", Description: "Fuzzy matching algorithm"},
			{Name: "-i", Description: "Case-insensitive match (default: smart-case match)"},
			{Name: "--literal", Description: "Do not normalize latin script letters before matching"},
			{Name: "-n", Description: "Comma-separated list of field index expressions for limiting search scope"},
			{Name: "--with-nth", Description: "Transform the presentation of each line using field index expressions"},
			{Name: "-d", Description: "Field delimiter regex (default: AWK-style)"},
			{Name: "--tac", Description: "Reverse the order of the input"},
			{Name: "--disabled", Description: "Do not perform search"},
			{Name: "--tiebreak", Description: "Comma-separated list of sort criteria to apply when the scores are tied"},
			{Name: "-m", Description: "Enable multi-select with tab/shift-tab"},
			{Name: "--no-mouse", Description: "Disable mouse"},
			{Name: "--bind", Description: "Custom key bindings. Refer to the man page"},
			{Name: "--cycle", Description: "Enable cyclic scroll"},
			{Name: "--keep-right", Description: "Keep the right end of the line visible on overflow"},
			{Name: "--no-hscroll", Description: "Disable horizontal scroll"},
			{Name: "--hscroll-off", Description: "Number of screen columns to keep to the right of the highlighted substring"},
			{Name: "--filepath-word", Description: "Make word-wise movements respect path separators"},
			{Name: "--jump-labels", Description: "Label characters for jump and jump-accept"},
			{Name: "--height", Description: "Height[%]"},
			{Name: "--min-height", Description: "Minimum height when --height is given in percent"},
			{Name: "--layout", Description: "Choose layout"},
			{Name: "--border", Description: "Draw border around the finder"},
			{Name: "--margin", Description: "Screen margin (TRBL | TB,RL | T,RL,B | T,R,B,L)"},
			{Name: "--padding", Description: "Padding inside border (TRBL | TB,RL | T,RL,B | T,R,B,L)"},
			{Name: "--info", Description: "Finder info style"},
			{Name: "--prompt", Description: "Input prompt"},
			{Name: "--pointer", Description: "Pointer to the current line"},
			{Name: "--marker", Description: "Multi-select marker"},
			{Name: "--header", Description: "String to print as header"},
			{Name: "--header-lines", Description: "The first N lines of the input are treated as header"},
			{Name: "--ansi", Description: "Enable processing of ANSI color codes"},
			{Name: "--tabstop", Description: "Number of spaces for a tab character"},
			{Name: "--color", Description: "Base scheme"},
			{Name: "--no-bold", Description: "Do not use bold text"},
			{Name: "--history", Description: "History file"},
			{Name: "--history-size", Description: "Maximum number of history entries"},
			{Name: "--preview", Description: "Command to preview highlighted line ({})"},
			{Name: "--preview-window", Description: "Preview window layout"},
		},
	})
}
