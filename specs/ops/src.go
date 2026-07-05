package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "src",
		Description: "Interact with Sourcegraph from the command line",
		Subcommands: []core.Subcommand{
			{Name: "apply", Description: "Apply batch"},
			{Name: "exec", Description: "Execute batch"},
			{Name: "new", Description: "New batch"},
			{Name: "preview", Description: "Preview batch"},
			{Name: "repositories", Description: "Repositories to batch"},
			{Name: "validate", Description: "Validate batch"},
			{Name: "edit", Description: "Edit config"},
			{Name: "get", Description: "Get configs"},
			{Name: "list", Description: "List configs"},
			{Name: "upload", Description: "Upload LSIF dump"},
			{Name: "create", Description: "Create an organization"},
			{Name: "delete", Description: "Delete an organization"},
			{Name: "members", Description: "List organization members"},
			{Name: "tag", Description: "Tag user"},
		},
		Options: []core.Option{
			{Name: "-display", Description: "Log GraphQL requests and responses to stdout"},
			{Name: "-explain-json", Description: "Explain the JSON output schema and exit"},
			{Name: "-get-curl", Description: "Skip validation of TLD certificates against trusted chains"},
			{Name: "-json", Description: "Whether or not to output results as JSON"},
			{Name: "-less", Description: "Pipe output to `less -R` (only if stdout is terminal, and not json flag)"},
			{Name: "-stream", Description: "Log the trace ID for requests"},
			{Name: "-user-agent-telemetry", Description: "Sourcegraph API Access"},
			{Name: "-insecure-skip-verify", Description: "Skip validation of TLS certificates against trusted chains"},
			{Name: "-query", Description: "Log the trace ID for requests"},
			{Name: "-trace", Description: "Log the trace ID for requests"},
			{Name: "-addr", Description: "Address on which to server (end with : for unused port)"},
			{Name: "-list", Description: "List found repository names"},
			{Name: "-context", Description: "Comma-separated list of key=value pairs to add to the script execution context"},
			{Name: "-dump-requests", Description: "Log GraphQL requests and responses to stdout"},
			{Name: "-secrets", Description: "Log the trace ID for requests"},
		},
	})
}
