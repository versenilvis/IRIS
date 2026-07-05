package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "npx",
		Description: "Execute binaries from npm packages",
		Options: []spec.Option{
			{Name: "--package", Description: "Package to be installed"},
			{Name: "--cache", Description: "Location of the npm cache"},
			{Name: "--always-spawn", Description: "Always spawn a child process to execute the command"},
			{Name: "-y", Description: "Execute npx command without prompting for confirmation"},
			{Name: "--call", Description: "Execute string as if inside `npm run-script`"},
			{Name: "--shell", Description: "Shell to execute the command with, if any"},
			{Name: "--ignore-existing", Description: "Suppress output from npx itself. Subcommands will not be affected"},
			{Name: "--npm", Description: "Npm binary to use for internal operations"},
		},
	})
}
