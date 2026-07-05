package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "pgcli",
		Description: "Host address of the postgres database",
		Options: []core.Option{
			{Name: "-h", Description: "Host address of the postgres database"},
			{Name: "-p", Description: "Port number at which the postgres instance is listening"},
			{Name: "-U", Description: "Username to connect to the postgres database"},
			{Name: "-u", Description: "Username to connect to the postgres database"},
			{Name: "-W", Description: "Force password prompt"},
			{Name: "-w", Description: "Never prompt for password"},
			{Name: "--single-connection", Description: "Do not use a separate connection for completions"},
			{Name: "-v", Description: "Version of pgcli"},
			{Name: "-d", Description: "Database name to connect to"},
			{Name: "--pgclirc", Description: "Location of pgclirc file"},
			{Name: "-D", Description: "Use DSN configured into the [alias_dsn] section of pgclirc file"},
			{Name: "--list-dsn", Description: "List of DSN configured into the [alias_dsn] section of pgclirc file"},
			{Name: "--row-limit", Description: "Set threshold for row limit prompt. Use 0 to disable prompt"},
			{Name: "--less-chatty", Description: "Skip intro on startup and goodbye on exit"},
			{Name: "--prompt", Description: "List available databases, then exit"},
			{Name: "--auto-vertical-output", Description: "Warn before running a destructive query"},
			{Name: "--ssh-tunnel", Description: "Open an SSH tunnel to the given address and connect to the database from it"},
			{Name: "--help", Description: "Show this message and exit"},
		},
	})
}
