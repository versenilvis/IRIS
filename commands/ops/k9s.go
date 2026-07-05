package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "k9s",
		Description: "Kubernetes namespace",
		Subcommands: []spec.Subcommand{
			{Name: "help", Description: "Help about any command"},
			{Name: "info", Description: "Print configuration info"},
			{Name: "version", Description: "Print version/build info"},
		},
		Options: []spec.Option{
			{Name: "-h", Description: "Help for info"},
			{Name: "-s", Description: "Prints K9s version info in short format"},
			{Name: "-A", Description: "Launch K9s in all namespaces"},
			{Name: "--as", Description: "Username to impersonate for the operation"},
			{Name: "--as-group", Description: "Group to impersonate for the operation"},
			{Name: "--certificate-authority", Description: "Path to a cert file for the certificate authority"},
			{Name: "--client-key", Description: "Path to a client key file for TLS"},
			{Name: "-c", Description: "Overrides the default resource to load when the application launches"},
			{Name: "--context", Description: "The name of the kubeconfig context to use"},
			{Name: "--crumbsless", Description: "Turn K9s crumbs off"},
			{Name: "--headless", Description: "Turn K9s header off"},
			{Name: "--insecure-skip-tls-verify", Description: "If true, the server's caCertFile will not be checked for validity"},
			{Name: "--kubeconfig", Description: "Path to the kubeconfig file to use for CLI requests"},
			{Name: "--logFile", Description: "Specify the log file"},
			{Name: "-l", Description: "Specify a log level (info, warn, debug, trace, error) (default 'info')"},
			{Name: "--logoless", Description: "Turn K9s logo off"},
			{Name: "-n", Description: "If present, the namespace scope for this CLI request"},
			{Name: "--readonly", Description: "Sets readOnly mode by overriding readOnly configuration setting"},
			{Name: "-r", Description: "Specify the default refresh rate as an integer (sec) (default 2)"},
			{Name: "--request-timeout", Description: "The length of time to wait before giving up on a single server request"},
			{Name: "--screen-dump-dir", Description: "Sets a path to a dir for a screen dumps"},
			{Name: "--token", Description: "Bearer token for authentication to the API server"},
			{Name: "--user", Description: "The name of the kubeconfig user to use"},
			{Name: "--write", Description: "Sets write mode by overriding the readOnly configuration setting"},
		},
	})
}
