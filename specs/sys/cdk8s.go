package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "cdk8s",
		Description: "CDK for K8s",
		Subcommands: []core.Subcommand{
			{Name: "init", Description: "Create a new, empty CDK8S project"},
			{Name: "type", Description: "Select language you are using"},
			{Name: "import", Description: "Import a CRD schema to generate generate resources"},
			{Name: "spec", Description: "Path to the CRD schema"},
			{Name: "synth", Description: "Synthesizes Kubernetes manifests for all charts in your app"},
		},
		Options: []core.Option{
			{Name: "--language", Description: "Output programming language"},
			{Name: "--class-prefix", Description: "A prefix to add to all generated class names"},
			{Name: "--no-class-prefix", Description: "Does not add a prefix to generated class names"},
			{Name: "--exclude", Description: "Do not import types that match these regular expressions"},
			{Name: "--output", Description: "Output directory"},
			{Name: "--app", Description: "Command to use in order to execute cdk8s app"},
			{Name: "--stdout", Description: "Write synthesized manifests to STDOUT instead of the output directory"},
			{Name: "--plugins-dir", Description: "Directory to store cdk8s plugins"},
			{Name: "--validate", Description: "Apply validation plugins on the resulting manifests"},
			{Name: "--no-validate", Description: "Disable validation"},
			{Name: "--version", Description: "The current version"},
			{Name: "--help", Description: "Show help"},
			{Name: "--check-upgrade", Description: "Check for cdk8s-cli upgrade"},
		},
	})
}
