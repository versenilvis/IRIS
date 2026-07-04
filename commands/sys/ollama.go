package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "ollama",
		Description: "A command-line tool for managing and deploying machine learning models",
		Subcommands: []core.Subcommand{
			{Name: "serve", Description: "Start ollama"},
			{Name: "create", Description: "Create a model from a Modelfile"},
			{Name: "show", Description: "Show information for a model"},
			{Name: "run", Description: "Run a model"},
			{Name: "stop", Description: "Stop the ollama server"},
			{Name: "pull", Description: "Pull a model from a registry"},
			{Name: "push", Description: "Push a model to a registry"},
			{Name: "list", Description: "List models"},
			{Name: "ps", Description: "List running models"},
			{Name: "cp", Description: "Copy a model"},
			{Name: "rm", Description: "Remove a model"},
			{Name: "help", Description: "Help about any command"},
		},
		Options: []core.Option{
			{Name: "-f", Description: "Specify Modelfile"},
			{Name: "--verbose", Description: "Enable verbose output"},
			{Name: "--help", Description: "Show help for ollama"},
			{Name: "--version", Description: "Show version information"},
		},
	})
}
