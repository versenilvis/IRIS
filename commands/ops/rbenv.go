package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "rbenv",
		Description: "List all available rbenv commands",
		Subcommands: []spec.Subcommand{
			{Name: "commands", Description: "List all available rbenv commands"},
			{Name: "global", Description: "Set or show the global Ruby version"},
			{Name: "install", Description: "Install a Ruby version using ruby-build"},
			{Name: "local", Description: "Set or show the local application-specific Ruby version"},
			{Name: "rehash", Description: "Rehash rbenv shims (run this after installing executables)"},
			{Name: "shell", Description: "Set or show the shell-specific Ruby version"},
			{Name: "uninstall", Description: "Uninstall a specific Ruby version"},
			{Name: "versions", Description: "List installed Ruby versions"},
			{Name: "whence", Description: "List all Ruby versions that contain the given executable"},
			{Name: "which", Description: "Display the full path to an executable"},
		},
		Options: []spec.Option{
			{Name: "--unset", Description: "List all available rbenv commands"},
			{Name: "--sh", Description: "Set or show the global Ruby version"},
			{Name: "--version", Description: "Show version of ruby-build"},
			{Name: "-f", Description: "If the version does not exist, do not display an error message"},
		},
	})
}
