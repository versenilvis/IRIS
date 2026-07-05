package rust

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "wasm-pack",
		Description: "Build an npm package",
		Subcommands: []core.Subcommand{
			{Name: "help", Description: "Prints this message or the help of the given subcommand(s)"},
			{Name: "test", Description: "Run tests for WebAssembly module"},
		},
		Options: []core.Option{
			{Name: "--help", Description: "Show help for wasm-pack or for the given subcommand(s)"},
			{Name: "--quiet", Description: "Suppress output from stdout"},
			{Name: "--version", Description: "Show version for wasm-pack"},
			{Name: "--verbose", Description: "Log verbosity is based off the number of v used"},
			{Name: "--log-level", Description: "The maximum level of messages that should be logged by wasm-pack"},
		},
	})
}
