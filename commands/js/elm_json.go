package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "elm-json",
		Description: "Deal with your elm.json",
		Subcommands: []core.Subcommand{
			{Name: "help", Description: "Prints help information or the help of the given subcommand(s)"},
			{Name: "install", Description: "Install a package"},
			{Name: "PACKAGE", Description: "Package to install, e.g. elm/core or elm/core@1.0.2 or elm/core@1"},
			{Name: "new", Description: "Create a new elm.json file"},
			{Name: "tree", Description: "List entire dependency graph as a tree"},
			{Name: "uninstall", Description: "Uninstall a package"},
			{Name: "upgrade", Description: "Bring your dependencies up to date"},
			{Name: "INPUT", Description: "The elm.json file to upgrade [default: elm.json]"},
		},
		Options: []core.Option{
			{Name: "--help", Description: "Prints help information"},
			{Name: "--test", Description: "Install as a test-dependency"},
			{Name: "--version", Description: "Prints version information"},
			{Name: "--yes", Description: "Package to install, e.g. elm/core or elm/core@1.0.2 or elm/core@1"},
			{Name: "--unsafe", Description: "Allow major versions bumps"},
			{Name: "--offline", Description: "Enable offline mode, which means no HTTP traffic will happen"},
			{Name: "--verbose", Description: "Sets the level of verbosity"},
		},
	})
}
