package runner

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "laravel",
		Description: "The output format (txt, xml, json, or md)",
		Subcommands: []core.Subcommand{
			{Name: "new", Description: "Create a new Laravel application"},
		},
		Options: []core.Option{
			{Name: "--format", Description: "The output format (txt, xml, json, or md)"},
			{Name: "--raw", Description: "To output raw command list"},
			{Name: "--dev", Description: "Initialize a Git repository"},
			{Name: "--branch", Description: "The branch that should be created for a new repository"},
			{Name: "--github", Description: "Create a new repository on GitHub"},
			{Name: "--jet", Description: "Installs the Laravel Jetstream scaffolding"},
			{Name: "--stack", Description: "The Jetstream stack that should be installed"},
			{Name: "--teams", Description: "Indicates whether Jetstream should be scaffolded with team support"},
			{Name: "--prompt-jetstream", Description: "Issues a prompt to determine if Jetstream should be installed"},
			{Name: "-f", Description: "Forces install even if the directory already exists"},
			{Name: "-h", Description: "Display the help message"},
			{Name: "-q", Description: "Do not output any message"},
			{Name: "-V", Description: "Display this application version"},
			{Name: "--ansi", Description: "Force ANSI output"},
			{Name: "--no-ansi", Description: "Disable ANSI output"},
			{Name: "-n", Description: "Do not ask any interactive question"},
		},
	})
}
