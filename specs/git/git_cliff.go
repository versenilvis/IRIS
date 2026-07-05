package git

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "git-cliff",
		Description: "Increases the logging verbosity",
		Options: []core.Option{
			{Name: "--verbose", Description: "Increases the logging verbosity"},
			{Name: "--init", Description: "Writes the default configuration file to cliff.toml"},
			{Name: "--latest", Description: "Processes the commits starting from the latest tag"},
			{Name: "--current", Description: "Processes the commits that belong to the current tag"},
			{Name: "--unreleased", Description: "Processes the commits that do not belong to a tag"},
			{Name: "--date-order", Description: "Sorts the tags chronologically"},
			{Name: "--context", Description: "Prints changelog context as JSON"},
			{Name: "--help", Description: "Prints help information"},
			{Name: "--version", Description: "Prints version information"},
			{Name: "--config", Description: "Sets the configuration file"},
			{Name: "--workdir", Description: "Sets the working directory"},
			{Name: "--repository", Description: "Sets the git repository"},
			{Name: "--include-path", Description: "Sets the path to include related commits"},
			{Name: "--exclude-path", Description: "Sets the path to exclude related commits"},
			{Name: "--with-commit", Description: "Sets custom commit messages to include in the changelog"},
			{Name: "--prepend", Description: "Prepends entries to the given changelog file"},
			{Name: "--output", Description: "Writes output to the given file"},
			{Name: "--tag", Description: "Sets the tag for the latest version"},
			{Name: "--body", Description: "Sets the template for the changelog body"},
			{Name: "--strip", Description: "Strips the given parts from the changelog"},
			{Name: "--sort", Description: "Sets sorting of the commits inside sections"},
		},
	})
}
