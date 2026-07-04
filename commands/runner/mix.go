package runner

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "mix",
		Description: "Build tool for Elixir",
		Subcommands: []core.Subcommand{
			{Name: "new", Description: "Creates a new Elixir project at the given path"},
			{Name: "run", Description: "Starts the current application and runs code"},
			{Name: "code", Description: "String containing code to execute"},
			{Name: "file|pattern", Description: "The file|pattern to execute"},
		},
		Options: []core.Option{
			{Name: "-h", Description: "Output usage information"},
			{Name: "--app", Description: "Name the OTP application for the project"},
			{Name: "--modules", Description: "Name the modules in the generated code skeleton"},
			{Name: "--sup", Description: "Generate an umbrella project"},
			{Name: "--config", Description: "Loads the given configuration files"},
			{Name: "-e", Description: "Evaluates the given code"},
			{Name: "-r", Description: "Executes the given pattern/file"},
			{Name: "-p", Description: "Executes the given pattern/file"},
			{Name: "--preload-modules", Description: "Preloads all modules defined in applications"},
			{Name: "--no-compile", Description: "Does not compile even if files require compilation"},
			{Name: "--no-deps-check", Description: "Does not check dependencies"},
			{Name: "--no-archives-check", Description: "Does not check archives"},
			{Name: "--no-halt", Description: "Does not halt the system after running the command"},
			{Name: "--no-mix-exs", Description: "Allows the command to run even if there is no mix.exs"},
			{Name: "--no-start", Description: "Does not start applications after compilation"},
			{Name: "--no-elixir-version-check", Description: "Does not check the Elixir version from mix.exs"},
			{Name: "--search", Description: "Prints all tasks and aliases that contain 'pattern' in the name"},
			{Name: "--names", Description: "Prints all task names and aliases"},
			{Name: "-v", Description: "Shows versioning information"},
		},
	})
}
