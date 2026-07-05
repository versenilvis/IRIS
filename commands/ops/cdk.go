package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "cdk",
		Description: "AWS CDK CLI",
		Subcommands: []spec.Subcommand{
			{Name: "init", Description: "Create a new, empty CDK project from a template"},
			{Name: "metadata", Description: "Returns all metadata associated with this stack"},
			{Name: "doctor", Description: "Check your set-up for potential problems"},
			{Name: "diff", Description: "Compares the specified stack with the deployed stack"},
			{Name: "destroy", Description: "Destroy the specified stack(s)"},
			{Name: "deploy", Description: "Deploy the specified stack(s) into your AWS account"},
			{Name: "bootstrap", Description: "Deploys the CDK toolkit stack into an AWS environment"},
			{Name: "synth", Description: "Synthesizes and prints the CloudFormation template for this stack"},
			{Name: "ls", Description: "List all stacks in the app"},
			{Name: "import", Description: "Import existing resource(s) into the given STACK"},
			{Name: "watch", Description: "Shortcut for 'deploy --watch'"},
			{Name: "ack", Description: "Acknowledge a notice so that it does not show up anymore"},
			{Name: "notices", Description: "Returns a list of relevant notices"},
			{Name: "context", Description: "Manage cached context values"},
			{Name: "doc", Description: "Opens the reference documentation in a browser"},
		},
		Options: []spec.Option{
			{Name: "--version", Description: "The current version"},
			{Name: "-h", Description: "Show help"},
		},
	})
}
