package runner

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "rvm",
		Description: "Show version of rvm",
		Subcommands: []core.Subcommand{
			{Name: "use", Description: "Setup current shell to use a specific ruby version"},
			{Name: "reset", Description: "Remove default and current settings, exit the shell"},
			{Name: "info", Description: "Show information about the current ruby environment"},
			{Name: "list", Description: "List currently installed version"},
			{Name: "reload", Description: "Reload RVM source itself (useful after changing RVM source)"},
			{Name: "implode", Description: "Remove all ruby installations it manages"},
			{Name: "get", Description: "Upgrades RVM to the stable or git head branches"},
			{Name: "do", Description: "Runs a named ruby file against specified and/or all rubies"},
			{Name: "install", Description: "Install one or many ruby versions"},
			{Name: "upgrade", Description: "Install new ruby, copy gemsets, make gems pristine, remove old rubies"},
			{Name: "reinstall", Description: "Remove ruby, install it, make gems pristine"},
			{Name: "uninstall", Description: "Uninstall one or many ruby versions, leaves their sources"},
			{Name: "remove", Description: "Remove one or many ruby versions, including their sources"},
		},
		Options: []core.Option{
			{Name: "--version", Description: "Show version of rvm"},
			{Name: "--default", Description: "When used with ruby selection, sets a default ruby for new shells"},
			{Name: "--debug", Description: "Enable debug mode"},
			{Name: "--force", Description: "Force install, removes old install & source before install"},
			{Name: "--all", Description: "Used with 'rvm list' to display 'most' available versions"},
			{Name: "--summary", Description: "Used with 'do' to print out a summary of the commands run"},
			{Name: "-C", Description: "Rubinius"},
			{Name: "--help", Description: "Show help for rvm"},
		},
	})
}
