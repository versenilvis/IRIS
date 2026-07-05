package rust

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "dprint",
		Description: "Prints the help of the given subcommand(s)",
		Subcommands: []spec.Subcommand{
			{Name: "init", Description: "Initializes a configuration file in the current directory"},
			{Name: "fmt", Description: "Formats the source files and writes the result to the file system"},
			{Name: "check", Description: "Checks for any files that haven't been formatted"},
			{Name: "config", Description: "Functionality related to the configuration file"},
			{Name: "add", Description: "Adds a plugin to the configuration file"},
			{Name: "update", Description: "Updates the plugins in the configuration file"},
			{Name: "output-format-times", Description: "Prints the amount of time it takes to format each file. Use this for debugging"},
			{Name: "clear-cache", Description: "Deletes the plugin cache directory"},
			{Name: "license", Description: "Outputs the software license"},
		},
		Options: []spec.Option{
			{Name: "--excludes", Description: "A pluggable and configurable code formatting platform written in Rust"},
			{Name: "--diff", Description: "Outputs a check-like diff of every formatted file"},
			{Name: "--stdin", Description: "Checks for any files that haven't been formatted"},
			{Name: "-y", Description: "Upgrade process plugins without prompting to confirm checksums"},
			{Name: "-c", Description: "Prints help information"},
			{Name: "--plugins", Description: "Prints additional diagnostic information"},
		},
	})
}
