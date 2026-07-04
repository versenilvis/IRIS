package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "oh-my-posh",
		Description: "The config file to use",
		Subcommands: []core.Subcommand{
			{Name: "init", Description: "Initialize oh-my-posh for your shell"},
			{Name: "get", Description: "Get oh-my-posh values"},
			{Name: "shell", Description: "Get the current shell, example usage: 'oh-my-posh get shell'"},
			{Name: "debug", Description: "Debug oh-my-posh, example usage: 'oh-my-posh debug'"},
			{Name: "config", Description: "Interact with the oh-my-posh configuration"},
			{Name: "edit", Description: "Edit the config file, example usage: 'oh-my-posh config edit'"},
			{Name: "OUTPUT", Description: "The file to write to"},
			{Name: "print", Description: "Print a prompt"},
			{Name: "WIDTH", Description: "The terminal width"},
			{Name: "COMMAND", Description: "The tooltip command"},
			{Name: "ERROR CODE", Description: "The last error code"},
			{Name: "EXECUTION TIME", Description: "The last command's execution time"},
			{Name: "POWERSHELL WORKING DIRECTORY", Description: "The working directory according to PowerShell"},
			{Name: "WORKING DIRECTORY", Description: "The working directory"},
			{Name: "SHELL", Description: "The shell used"},
			{Name: "NUM", Description: "The number of stacks"},
			{Name: "version", Description: "Display the oh-my-posh version, example usage: 'oh-my-posh version'"},
		},
		Options: []core.Option{
			{Name: "--config", Description: "The config file to use"},
			{Name: "--format", Description: "The file format to use"},
			{Name: "--print", Description: "Print the init script"},
			{Name: "--write", Description: "The file to write to"},
			{Name: "--terminal-width", Description: "The terminal width"},
			{Name: "--command", Description: "The tooltip command"},
			{Name: "--error", Description: "The last exit code, example usage: 'oh-my-posh print primary --error 127'"},
			{Name: "--eval", Description: "Use eval to render the prompt, example usage: 'oh-my-posh print primary --eval'"},
			{Name: "--execution-time", Description: "The last command's execution time"},
			{Name: "--plain", Description: "The working directory according to PowerShell"},
			{Name: "--pwd", Description: "The working directory"},
			{Name: "--shell", Description: "The current shell, example usage: 'oh-my-posh print primary --shell fish'"},
			{Name: "--stack-count", Description: "The number of stacks"},
			{Name: "--help", Description: "Show help for oh-my-posh"},
		},
	})
}
