package python

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "jupyter",
		Description: "Set log level to logging.DEBUG (maximize logging output)",
		Subcommands: []core.Subcommand{
			{Name: "bundlerextension", Description: "Work with Jupyter bundler extensions"},
			{Name: "enable", Description: "Enable a bundler extension"},
			{Name: "disable", Description: "Disable a bundler extension"},
			{Name: "list", Description: "List bundler extensions"},
			{Name: "kernel", Description: "Run a kernel locally in a subprocess"},
			{Name: "kernelspec", Description: "Manage Jupyter kernel specifications"},
			{Name: "install", Description: "Install a kernel specification directory"},
			{Name: "uninstall", Description: "Alias for remove"},
			{Name: "remove", Description: "Remove one or more Jupyter kernelspecs by name"},
			{Name: "migrate", Description: "Migrate configuration and data from .ipython prior to 4.0 to Jupyter locations"},
			{Name: "notebook", Description: "Run the Jupyter notebook server"},
			{Name: "stop", Description: "Stop currently running notebook server"},
			{Name: "password", Description: "Set a password for the notebook server"},
			{Name: "run", Description: "Run a notebook"},
			{Name: "serverextension", Description: "Manage server extensions"},
			{Name: "troubleshoot", Description: "Log for troubleshooting"},
			{Name: "trust", Description: "Manage trust"},
		},
		Options: []core.Option{
			{Name: "--debug", Description: "Set log level to logging.DEBUG (maximize logging output)"},
			{Name: "-h", Description: "Show this message"},
			{Name: "--log-level", Description: "Set the log level by value or name"},
			{Name: "--config", Description: "Choose a config file"},
			{Name: "--user", Description: "Apply the operation only for the given user"},
			{Name: "--system", Description: "Apply the operation system-wide"},
			{Name: "--sys-prefix", Description: "Install from a Python package"},
			{Name: "--kernel", Description: "Manage Jupyter kernel specifications"},
			{Name: "--json", Description: "Output spec name and location as json"},
			{Name: "--generate-config", Description: "Generate default config file"},
			{Name: "-y", Description: "Answer yes to any questions instead of prompting"},
			{Name: "--help", Description: "Show help for jupyter"},
			{Name: "--version", Description: "Show the jupyter command's version and exit"},
			{Name: "--config-dir", Description: "Show Jupyter config dir"},
			{Name: "--data-dir", Description: "Show Jupyter data dir"},
			{Name: "--runtime-dir", Description: "Show Jupyter runtime dir"},
			{Name: "--paths", Description: "Show all Jupyter paths. Add --json for machine-readable format"},
		},
	})
}
