package runner

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "ruby",
		Description: "Interpreted object-oriented scripting language",
		Options: []core.Option{
			{Name: "--copyright", Description: "Prints the copyright notice"},
			{Name: "--version", Description: "Prints the version of Ruby interpreter"},
			{Name: "-0", Description: "Specifies the input record separator ($/) as an octal number"},
			{Name: "-C", Description: "Causes Ruby to switch to the directory"},
			{Name: "-F", Description: "Specifies input field separator ($;)"},
			{Name: "-I", Description: "Specifies KANJI (Japanese) encoding"},
			{Name: "-S", Description: "Turns on taint checks at the specified level (default 1)"},
			{Name: "-a", Description: "Turns on auto-split mode when used with -n or -p"},
			{Name: "-c", Description: "Turns on debug mode. $DEBUG will be set to true"},
			{Name: "-e", Description: "Prints a summary of the options"},
			{Name: "-i", Description: "Causes Ruby to assume the following loop around your script"},
			{Name: "-p", Description: "Causes Ruby to load the library using require"},
			{Name: "-s", Description: "Enables verbose mode"},
		},
	})
}
