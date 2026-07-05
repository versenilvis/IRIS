package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "ampx",
		Description: "CLI for Amplify Gen 2",
		Subcommands: []spec.Subcommand{
			{Name: "sandbox", Description: "Deploy to your personal cloud sandbox"},
			{Name: "delete", Description: "Delete your personal cloud sandbox"},
			{Name: "secret", Description: "Manage backend secrets for your personal cloud sandbox"},
			{Name: "set", Description: "Set a secret"},
			{Name: "remove", Description: "Remove a secret"},
			{Name: "list", Description: "List all available secrets"},
			{Name: "get", Description: "View details of a secret"},
			{Name: "generate", Description: "Generate supplemental information or code"},
			{Name: "outputs", Description: "Generate backend outputs file"},
			{Name: "graphql-client-code", Description: "Generate GraphQL statements and types"},
			{Name: "forms", Description: "Generate React form components"},
			{Name: "info", Description: "Generate system information for troubleshooting"},
			{Name: "pipeline-deploy", Description: "Deploy Amplify project in a CI/CD pipeline"},
		},
		Options: []spec.Option{
			{Name: "--dir-to-watch", Description: "Directory to watch for file changes"},
			{Name: "--exclude", Description: "Paths or glob patterns to ignore"},
			{Name: "--identifier", Description: "Name to distinguish between sandbox environments"},
			{Name: "--outputs-out-dir", Description: "Directory where client config file is written"},
			{Name: "--outputs-format", Description: "Format of the client config file"},
			{Name: "--outputs-version", Description: "Version of the configuration"},
			{Name: "--profile", Description: "AWS profile name"},
			{Name: "--stream-function-logs", Description: "Stream function execution logs"},
			{Name: "--logs-filter", Description: "Regex pattern to filter logs"},
			{Name: "--logs-out-file", Description: "File to append streaming logs"},
			{Name: "--name", Description: "Name to distinguish between sandbox environments"},
			{Name: "-y", Description: "Do not ask for confirmation before deleting"},
			{Name: "--format", Description: "Format of the configuration"},
			{Name: "--out-dir", Description: "Directory where config is written"},
			{Name: "--branch", Description: "Git branch name"},
			{Name: "--app-id", Description: "Amplify App ID"},
			{Name: "--out", Description: "Directory where config is written"},
			{Name: "--model-target", Description: "Modelgen export target"},
			{Name: "--statement-target", Description: "Graphql codegen statement export target"},
			{Name: "--type-target", Description: "Graphql-codegen type export target"},
			{Name: "--all", Description: "Show hidden options"},
			{Name: "--debug", Description: "Print debug logs to the console"},
			{Name: "--statement-max-depth", Description: "Maximum depth of the generated GraphQL statements"},
			{Name: "--models", Description: "Model names to generate"},
			{Name: "--stack", Description: "CloudFormation stack name"},
			{Name: "--help", Description: "Display help information"},
			{Name: "--version", Description: "Display version information"},
		},
	})
}
