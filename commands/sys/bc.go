package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "bc",
		Description: "An arbitrary precision calculator language",
		Options: []core.Option{
			{Name: "--help", Description: "Print the usage and exit"},
			{Name: "--interactive", Description: "Force interactive mode"},
			{Name: "--mathlib", Description: "Define the standard math library"},
			{Name: "--warn", Description: "Give warnings for extensions to POSIX bc"},
			{Name: "--standard", Description: "Process exactly the POSIX bc language"},
			{Name: "--quiet", Description: "Do not print the normal GNU bc welcome"},
			{Name: "--version", Description: "Print the version number and copyright and quit"},
		},
	})
}
