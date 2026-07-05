package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "nest",
		Description: "Report actions that would be taken without writing out results",
		Subcommands: []spec.Subcommand{
			{Name: "new", Description: "Creates a new nest project"},
			{Name: "project", Description: "The name of the project"},
			{Name: "generate", Description: "Generate and/or modifies files based on a schematic"},
			{Name: "application", Description: "Generate a new application workspace"},
			{Name: "class", Description: "Generate a new class"},
			{Name: "configuration", Description: "Generate a CLI configuration file"},
			{Name: "controller", Description: "Generate a controller declaration"},
			{Name: "decorator", Description: "Generate a custom decorator"},
			{Name: "filter", Description: "Generate a filter declaration"},
			{Name: "gateway", Description: "Generate a gateway declaration"},
			{Name: "guard", Description: "Generate a guard declaration"},
			{Name: "interceptor", Description: "Generate an interceptor declaration"},
			{Name: "interface", Description: "Generate an interface"},
			{Name: "middleware", Description: "Generate a middleware declaration"},
			{Name: "module", Description: "Generate a module declaration"},
			{Name: "pipe", Description: "Generate a pipe declaration"},
			{Name: "provider", Description: "Generate a provider declaration"},
			{Name: "resolver", Description: "Generate a GraphQL resolver declaration"},
			{Name: "service", Description: "Generate a service declaration"},
			{Name: "library", Description: "Generate a new library within a monorepo"},
			{Name: "sub-app", Description: "Generate a new application within a monorepo"},
			{Name: "resource", Description: "Generate a new CRUD resource"},
			{Name: "build", Description: "Builds Nest application"},
			{Name: "app", Description: "The name of the app"},
			{Name: "start", Description: "Run Nest application"},
			{Name: "info", Description: "Display Nest project details"},
			{Name: "update", Description: "Update Nest dependencies"},
		},
		Options: []spec.Option{
			{Name: "-d", Description: "Report actions that would be taken without writing out results"},
			{Name: "-p", Description: "Project in which to generate files"},
			{Name: "--flat", Description: "Enforce flat structure of generated element"},
			{Name: "--spec", Description: "Enforce spec files generation (default: true)"},
			{Name: "--no-spec", Description: "Disable spec files generation"},
			{Name: "-c", Description: "Schematics collection to use"},
			{Name: "--help", Description: "Show help for nest"},
		},
	})
}
