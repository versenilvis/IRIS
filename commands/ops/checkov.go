package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "checkov",
		Description: "Branch",
		Options: []core.Option{
			{Name: "--help", Description: "Show help for checkov"},
			{Name: "--version", Description: "Show the version of checkov"},
			{Name: "--quiet", Description: "CLI output, display only failed checks"},
			{Name: "--compact", Description: "CLI output, do not display code blocks"},
			{Name: "--list", Description: "List checks"},
			{Name: "--no-guide", Description: "IaC root directory (can not be used together with --file)"},
			{Name: "--output", Description: "IaC frameworks to include checks for"},
			{Name: "--skip-framework", Description: "IaC frameworks to exclude checks for"},
			{Name: "--add-check", Description: "Generate a new check via CLI prompt"},
			{Name: "--file", Description: "IaC file(can not be used together with --directory)"},
			{Name: "--skip-path", Description: "Directory for custom checks to be loaded. Can be repeated"},
			{Name: "--bc-api-key", Description: "Bridgecrew API key. You may also use the environment variable: BC_API_KEY"},
			{Name: "--docker-image", Description: "Scan docker images by name or ID. Only works with --bc-api-key flag"},
			{Name: "--dockerfile-path", Description: "Path to the Dockerfile of the scanned docker image"},
			{Name: "--repo-id", Description: "Identity string of the repository, with form <repo_owner>/<repo_name>"},
			{Name: "--branch", Description: "Evaluate the values of variables and locals"},
			{Name: "--ca-certificate", Description: "Path to the Checkov configuration YAML file"},
			{Name: "--create-config", Description: "Runs checks but suppresses the error code"},
			{Name: "--soft-fail-on", Description: "Set minimum severity to return a non-zero exit code"},
			{Name: "--skip-cve-package", Description: "Use the Enforcement rules configured in the platform for hard / soft fail logic"},
			{Name: "--support", Description: "Enable debug logs and upload the logs to the server"},
			{Name: "--summary-position", Description: "Chose whether the summary will be appended on top or on bottom"},
			{Name: "--skip-download", Description: "Do not download any data from Prisma Cloud"},
			{Name: "--secrets-history-timeout", Description: "Maximum time to run the history scan"},
			{Name: "--scan-secrets-history", Description: "Will scan the history of commits for secrets"},
			{Name: "--prisma-api-url", Description: "The Prisma Cloud API URL"},
			{Name: "--policy-metadata-filter", Description: "Name of the output folder to save the chosen output formats"},
			{Name: "--output-baseline-as-skipped", Description: "Output checks that are skipped due to baseline file presence"},
			{Name: "--openai-api-key", Description: "Return exit code 0 instead of 2"},
			{Name: "--mask", Description: "Each entry in the list will be used for masking the desired attribute"},
			{Name: "--deep-analysis", Description: "Enable combine TF graph and TF Plan graph"},
			{Name: "--block-list-secret-scan", Description: "List of files to filter out from the secret scanner"},
		},
	})
}
