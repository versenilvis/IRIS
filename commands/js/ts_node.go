package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "ts-node",
		Description: "Run the TypeScript interpreter for Node.JS",
		Options: []core.Option{
			{Name: "--help", Description: "Show help for ts-node"},
			{Name: "-v", Description: "Print version information of the ts-node module"},
			{Name: "-e", Description: "Evaluate script"},
			{Name: "-p", Description: "Evaluate script and print result"},
			{Name: "-r", Description: "Require module before executing"},
			{Name: "-i", Description: "Always open interactive REPL"},
			{Name: "--show-config", Description: "Print resolved Typescript config to the terminal"},
			{Name: "--cwd-mode", Description: "Resolve Typescript config based on the current working directory"},
			{Name: "-T", Description: "Use the Typescript transpile module mode"},
			{Name: "-H", Description: "Use the Typescript compiler host API"},
			{Name: "-I", Description: "Ignore patterns from Typescript compilation"},
			{Name: "-P", Description: "Specify TypeScript project location"},
			{Name: "-C", Description: "Use a custom compiler"},
			{Name: "--transpiler", Description: "Use a custom transpiler"},
			{Name: "-D", Description: "Specify Typescript diagnostic code to ignore"},
			{Name: "-O", Description: "JSON object that will be merged with the compiler options"},
			{Name: "--cwd", Description: "Specify working directory"},
			{Name: "--files", Description: "Load files, include and exclude from Typescript config on startup"},
			{Name: "--pretty", Description: "Use the pretty formatter for diagnostic errors"},
			{Name: "--skip-project", Description: "Skip reading Typescript config"},
			{Name: "--scope", Description: "Scope compilation to scope directory specified"},
			{Name: "--scope-dir", Description: "Directory for scope parameter"},
			{Name: "--skip-ignore", Description: "Skip --ignore checks"},
			{Name: "--prefer-ts-exts", Description: "Prefer Typescript files over JavaScript files when importing files"},
			{Name: "--log-error", Description: "Pipe Typescript errors to stderr instead of throwing exceptions"},
			{Name: "--no-experimental-repl-await", Description: "Disable the top-level await function in REPL"},
		},
	})
}
