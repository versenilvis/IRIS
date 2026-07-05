package text

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "jq",
		Description: "Output the jq version and exit with zero",
		Options: []spec.Option{
			{Name: "--version", Description: "Output the jq version and exit with zero"},
			{Name: "--seq", Description: "Parse the input in streaming fashion, outputting arrays of path and leaf values"},
			{Name: "--slurp", Description: "Use a tab for each indentation level instead of two spaces"},
			{Name: "--indent", Description: "Use the given number of spaces for indentation"},
			{Name: "--color-output", Description: "Disable color"},
			{Name: "--ascii-output", Description: "Flush the output after each JSON object is printed"},
			{Name: "--sort-keys", Description: "Output the fields of each object with the keys in sorted orde"},
			{Name: "--raw-output", Description: "Like -r but jq won't print a newline after each output"},
			{Name: "-f", Description: "Read filter from the file rather than from a command line"},
			{Name: "-L", Description: "Prepend directory to the search list for modules"},
			{Name: "-e", Description: "This option passes a value to the jq program as a predefined variable"},
			{Name: "--argjson", Description: "Command-line JSON processor"},
		},
	})
}
