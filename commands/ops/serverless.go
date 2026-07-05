package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "serverless",
		Description: "AWS profile to use with the command",
		Subcommands: []spec.Subcommand{
			{Name: "deploy", Description: "Deploy a Serverless service"},
			{Name: "info", Description: "Display information about the service"},
			{Name: "invoke", Description: "Invoke a deployed function"},
			{Name: "logs", Description: "Output the logs of a deployed function"},
			{Name: "metrics", Description: "Show metrics for a specific function"},
			{Name: "remove", Description: "Remove Serverless service and all resources"},
			{Name: "rollback", Description: "Rollback the Serverless service to a specific deployment"},
			{Name: "studio", Description: "Develop a Serverless application in the cloud"},
			{Name: "test", Description: "Run HTTP tests"},
			{Name: "package", Description: "Package a Serverless service"},
			{Name: "plugin", Description: "Handle plugins"},
			{Name: "print", Description: "Print your compiled and resolved config file"},
			{Name: "create", Description: "Create new Serverless service"},
			{Name: "dashboard", Description: "Open the Serverless dashboard"},
			{Name: "generate-event", Description: "Generate event"},
			{Name: "login", Description: "Login or sign up for Serverless"},
			{Name: "logout", Description: "Logout from Serverless"},
			{Name: "output", Description: "Get/list value of dashboard deployment profile parameter"},
			{Name: "param", Description: "Get/list value of dashboard service output"},
			{Name: "slstats", Description: "Enable or disable stats"},
		},
		Options: []spec.Option{
			{Name: "--aws-profile", Description: "AWS profile to use with the command"},
			{Name: "--function", Description: "Name of the function"},
			{Name: "--region", Description: "Region of the service"},
			{Name: "--use-local-credentials", Description: "Path to serverless config file"},
			{Name: "--app", Description: "Dashboard app"},
			{Name: "--org", Description: "Dashboard org"},
			{Name: "--help", Description: "Show help"},
			{Name: "--version", Description: "Show version info"},
			{Name: "--stage", Description: "Stage of the service"},
			{Name: "--verbose", Description: "Show all stack events during deployment"},
			{Name: "--package", Description: "Path of the deployment package"},
			{Name: "--conceal", Description: "Hide secrets from the output (e.g. API Gateway key values)"},
			{Name: "--qualifier", Description: "Version number or alias to invoke"},
			{Name: "--path", Description: "Path to JSON or YAML file holding input data"},
			{Name: "--type", Description: "Type of invocation"},
			{Name: "--log", Description: "Trigger logging data output"},
			{Name: "--data", Description: "Input data"},
			{Name: "--raw", Description: "Flag to pass input data as a raw string"},
			{Name: "--context", Description: "Context of the service"},
			{Name: "--contextPath", Description: "Path to JSON or YAML file holding context data"},
			{Name: "--env", Description: "Override environment variables. e.g. --env VAR1=val1 --env VAR2=val2"},
			{Name: "--force", Description: "Forces a deployment to take place"},
			{Name: "--update-config", Description: "Forces a deployment to take place"},
			{Name: "--aws-s3-accelerate", Description: "Enables S3 Transfer Acceleration making uploading artifacts much faster"},
		},
	})
}
