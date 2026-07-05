package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "helmfile",
		Description: "Deploy helm charts",
		Subcommands: []spec.Subcommand{
			{Name: "apply", Description: "Apply all resources from state file only when there are changes"},
			{Name: "build", Description: "Build all resources from state file"},
			{Name: "cache", Description: "Cache management"},
			{Name: "completion", Description: "Generate the autocompletion script for the specified shell"},
			{Name: "destroy", Description: "Destroys and then purges releases"},
			{Name: "help", Description: "Help about any command"},
			{Name: "lint", Description: "Lint charts from state file (helm lint)"},
			{Name: "list", Description: "List releases defined in state file"},
			{Name: "repos", Description: "Repos releases defined in state file"},
			{Name: "secrets", Description: "Causes the helm-secrets plugin to be executed to decrypt the file"},
			{Name: "status", Description: "Retrieve status of releases in state file"},
			{Name: "sync", Description: "Sync releases defined in state file"},
			{Name: "template", Description: "Template releases defined in state file"},
			{Name: "test", Description: "Test charts from state file (helm test)"},
			{Name: "version", Description: "Print the CLI version"},
		},
		Options: []spec.Option{
			{Name: "--help", Description: "Do not exit with an error code if the provided selector has no matching releases"},
			{Name: "--allow-no-matching-release", Description: "Show help for helmfile"},
			{Name: "--chart", Description: "Output with color"},
			{Name: "--debug", Description: "Specify the environment name. defaults to default"},
			{Name: "--file", Description: "Request confirmation before attempting to modify clusters"},
			{Name: "--kube-context", Description: "Set kubectl context. Uses current context by default"},
			{Name: "--log-level", Description: "Set log level, default info (default info)"},
			{Name: "--namespace", Description: "Output without color"},
			{Name: "--quiet", Description: "Silence output. Equivalent to log-level warn"},
			{Name: "--selector", Description: "Specify state values in a YAML file"},
			{Name: "--state-values-set", Description: "Version for helmfile"},
		},
	})
}
