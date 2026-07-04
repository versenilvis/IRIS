package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "pnpx",
		Description: "Execute binaries from npm packages",
		Options: []core.Option{
			{Name: "--package", Description: "Package to be executed"},
			{Name: "--cache", Description: "Location of the npm cache"},
			{Name: "--always-spawn", Description: "Always spawn a child process to execute the command"},
			{Name: "--call", Description: "Execute string as if inside `npm run-script`"},
			{Name: "--shell", Description: "Shell to execute the command with, if any"},
			{Name: "--ignore-existing", Description: "Suppress output from pnpx itself. Subcommands will not be affected"},
			{Name: "--npm", Description: "Npm binary to use for internal operations"},
		},
	})
}
