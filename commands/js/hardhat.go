package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "hardhat",
		Description: "Ethereum development environment",
		Subcommands: []spec.Subcommand{
			{Name: "accounts", Description: "Prints the list of accounts"},
			{Name: "check", Description: "Check whatever you need"},
			{Name: "clean", Description: "Clears the cache and deletes all artifacts"},
			{Name: "compile", Description: "Compiles the entire project, building all artifacts"},
			{Name: "console", Description: "Opens a hardhat console"},
			{Name: "flatten", Description: "Flattens and prints contracts and their dependencies"},
			{Name: "help", Description: "Prints this message"},
			{Name: "node", Description: "Starts a JSON-RPC server on top of Hardhat Network"},
			{Name: "run", Description: "Runs a user-defined script after compiling the project"},
			{Name: "test", Description: "Runs mocha tests"},
		},
		Options: []spec.Option{
			{Name: "--emoji", Description: "Use emoji in messages"},
			{Name: "--max-memory", Description: "The maximum amount of memory that Hardhat can use"},
			{Name: "--help", Description: "Shows the help text or task's help if name is provided"},
			{Name: "--network", Description: "The network to connect to"},
			{Name: "--verbose", Description: "Enables Hardhat verbose logging"},
			{Name: "--version", Description: "Shows hardhat's version"},
			{Name: "--global", Description: "Clear the global cache"},
			{Name: "--force", Description: "Force compilation ignoring cache"},
			{Name: "--quiet", Description: "Makes the compilation process less verbose"},
			{Name: "--no-compile", Description: "Don't compile before running this task"},
			{Name: "--fork", Description: "The URL of the JSON-RPC server to fork from"},
			{Name: "--fork-block-number", Description: "The block number to fork from"},
			{Name: "--hostname", Description: "The host to which to bind to for new connections"},
			{Name: "--port", Description: "The port on which to listen for new connections (default: 8545)"},
		},
	})
}
