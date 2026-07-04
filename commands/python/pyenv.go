package python

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "pyenv",
		Description: "Pyenv",
		Subcommands: []core.Subcommand{
			{Name: "commands", Description: "Lists all available pyenv commands"},
			{Name: "local", Description: "Sets a local application-specific Python version"},
			{Name: "global", Description: "Sets the global version of Python to be used in all shells"},
			{Name: "shell", Description: "Sets a shell-specific Python version"},
			{Name: "install", Description: "Install a Python version"},
			{Name: "uninstall", Description: "Performs a deployment (default)"},
			{Name: "rehash", Description: "Performs a deployment (default)"},
			{Name: "whence", Description: "Lists all Python versions with the given command installed"},
		},
		Options: []core.Option{
			{Name: "-h", Description: "Output usage information"},
			{Name: "--unset", Description: "Sets the global version of Python to be used in all shells"},
			{Name: "-l", Description: "List all available versions"},
			{Name: "-f", Description: "Install even if the version appears to be installed already"},
			{Name: "-s", Description: "Skip the installation if the version appears to be installed already"},
			{Name: "-k", Description: "Keep source tree in $PYENV_BUILD_ROOT after installation"},
			{Name: "-v", Description: "Verbose mode: print compilation status to stdout"},
			{Name: "-p", Description: "Apply a patch from stdin before building"},
			{Name: "-g", Description: "Build a debug version"},
			{Name: "--bare", Description: "Print only the version names, one per line"},
			{Name: "--skip-aliases", Description: "Skip printing aliases"},
		},
	})
}
