package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "datree",
		Description: "Help for",
		Subcommands: []core.Subcommand{
			{Name: "completion", Description: "Generate completion script for bash,zsh,fish,powershell"},
			{Name: "config", Description: "Internal configuration management for datree config file"},
			{Name: "get", Description: "Get configuration value"},
			{Name: "set", Description: "Set configuration value"},
			{Name: "help", Description: "Help about any command"},
			{Name: "kustomize", Description: "Generate kustomization files from manifests"},
			{Name: "test", Description: "Test kustomization files"},
			{Name: "version", Description: "Print the version number"},
		},
		Options: []core.Option{
			{Name: "--help", Description: "Help for 'test'"},
			{Name: "--ignore-missing-schemas", Description: "Ignore missing schemas when executing schema validation step"},
			{Name: "--no-record", Description: "Do not send policy checks metadata to the backend"},
			{Name: "--only-k8s-files", Description: "Define output format (simple, yaml, json, xml, JUnit)"},
			{Name: "-p", Description: "Policy name to run against"},
			{Name: "--policy-config", Description: "Path for local policies configuration file"},
			{Name: "--schema-location", Description: "Override schemas location search path (can be specified multiple times)"},
			{Name: "-s", Description: "Set kubernetes version to validate against. Defaults to 1.19.0"},
			{Name: "--verbose", Description: "Display 'How to Fix' link"},
			{Name: "-h", Description: "Set configuration value"},
		},
	})
}
