package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "elm-format",
		Description: "Format your code in the Elm idiomatic way",
		Options: []spec.Option{
			{Name: "--help", Description: "Show help for elm-format"},
			{Name: "--output", Description: "Write output to FILE instead of overwriting the given source file"},
			{Name: "--yes", Description: "Reply 'yes' to all automated prompts"},
			{Name: "--validate", Description: "Check if files are formatted without changing them"},
			{Name: "--stdin", Description: "Read from stdin, output to stdout"},
			{Name: "--elm-version", Description: "The Elm version of the source files being formatted"},
		},
	})
}
