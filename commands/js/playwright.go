package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "playwright",
		Description: "Display help for command",
		Subcommands: []spec.Subcommand{
			{Name: "test", Description: "Run tests with Playwright Test"},
			{Name: "tests", Description: "Test files to run"},
			{Name: "install", Description: "Running without arguments will install default browsers"},
			{Name: "browsers", Description: "Browser to install"},
		},
		Options: []spec.Option{
			{Name: "--help", Description: "Display help for command"},
			{Name: "-g", Description: "Run the test with the title"},
			{Name: "--headed", Description: "Run tests in headed browsers"},
			{Name: "--with-deps", Description: "Install system dependencies for browsers"},
			{Name: "--version", Description: "Output the version number"},
		},
	})
}
