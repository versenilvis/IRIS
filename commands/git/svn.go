package git

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "svn",
		Description: "Specify a username ARG",
		Subcommands: []core.Subcommand{
			{Name: "help", Description: "Show help for svn"},
			{Name: "subcommand", Description: "Help about specific subcommand"},
			{Name: "status", Description: "Show the working tree status"},
			{Name: "info", Description: "Show information about a local or remote item"},
			{Name: "checkout", Description: "Check out a working copy from a repository"},
			{Name: "repository", Description: "The repository you want to checkout"},
			{Name: "commit", Description: "Commit to a repository"},
		},
		Options: []core.Option{
			{Name: "--password", Description: "Specify a password ARG"},
			{Name: "--password-from-stdin", Description: "Read password from stdin"},
			{Name: "--no-auth-cache", Description: "Do not cache authentication tokens"},
			{Name: "--non-interactive", Description: "Do interactive prompting even if standard input is not a terminal device"},
			{Name: "--trust-server-cert", Description: "Specify a password ARG"},
			{Name: "--trust-server-cert-failures", Description: "Specify a password ARG"},
			{Name: "--config-dir", Description: "Read user configuration files from directory ARG"},
			{Name: "--config-option", Description: "Specify a password ARG"},
			{Name: "-m", Description: "Use the given message as the commit message"},
			{Name: "--version", Description: "Show help for svn"},
			{Name: "--verbose", Description: "Show help for svn"},
			{Name: "--quiet", Description: "Show help for svn"},
		},
	})
}
