package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "st2",
		Description: "Show this help and exit",
		Subcommands: []core.Subcommand{
			{Name: "list", Description: "Get the list of actions"},
			{Name: "get", Description: "Get individual action"},
			{Name: "create", Description: "Create a new action"},
			{Name: "update", Description: "Update an existing action"},
			{Name: "delete", Description: "Delete an existing action"},
			{Name: "enable", Description: "Enable an existing action"},
			{Name: "disable", Description: "Disable an existing action"},
			{Name: "execute", Description: "Invoke an action manually"},
		},
		Options: []core.Option{
			{Name: "-h", Description: "Show this help and exit"},
			{Name: "-t", Description: "Print output in JSON format"},
			{Name: "-y", Description: "Print output in YAML format"},
			{Name: "--attr", Description: "Display full detail of the execution in table format"},
			{Name: "-k", Description: "How long (in milliseconds) to delay the execution before scheduling"},
			{Name: "--tail", Description: "Automatically start tailing new execution"},
			{Name: "--auto-dict", Description: "A trace tag string to track execution later"},
			{Name: "--trace-id", Description: "Existing trace id for this execution"},
			{Name: "-a", Description: "Do not wait for action to finish"},
			{Name: "-e", Description: "User under which to run the action (admins only)"},
			{Name: "-w", Description: "Set the width of the columns in output"},
			{Name: "-p", Description: "Only return resources belonging to the provided pack"},
			{Name: "-r", Description: "Return policies for the resource ref"},
			{Name: "-pt", Description: "Return policies of the policy type"},
			{Name: "-l", Description: "On successful authentication, print only token to the console"},
			{Name: "-n", Description: "List N most recent; use -n -1 to fetch the full result set"},
			{Name: "-s", Description: "Show all attributes"},
			{Name: "--action", Description: "Action reference to filter the list"},
			{Name: "--file", Description: "Local file path to the workflow definition"},
			{Name: "--status", Description: "Trigger instance id to filter the list"},
			{Name: "--show-secrets", Description: "Full list of attributes"},
			{Name: "-m", Description: "Optional metadata to associate with the API Keys"},
			{Name: "-x", Description: "Don't retrieve and display the result field"},
			{Name: "--tasks", Description: "Name of the workflow tasks to re-run"},
			{Name: "--no-reset", Description: "Type of output to tail for. If not provided, defaults to all"},
			{Name: "--types", Description: "Types of content to register"},
			{Name: "--include-metadata", Description: "Include metadata (timestamp, output type) with the output"},
			{Name: "--prefix", Description: "Only return values with names starting with the provided prefix"},
			{Name: "-d", Description: "Decrypt secrets and displays plain text"},
			{Name: "--encrypted", Description: "Scope item is under. Example: 'user'"},
			{Name: "--skip-dependencies", Description: "Skip dependencies"},
			{Name: "--iftt", Description: "Show trigger and action in display list"},
			{Name: "--enabled", Description: "Show enabled"},
			{Name: "--disabled", Description: "Show disabled"},
			{Name: "-g", Description: "Trigger type reference to filter the list"},
			{Name: "-ty", Description: "Execution to filter the list"},
			{Name: "--show-executions", Description: "Only show executions"},
			{Name: "--show-rules", Description: "Only show rules"},
			{Name: "--show-trigger-instances", Description: "Only show trigger instances"},
			{Name: "--group-id", Description: "Group ID"},
		},
	})
}
