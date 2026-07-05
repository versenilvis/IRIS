package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "dcli",
		Description: "Display help for command",
		Subcommands: []spec.Subcommand{
			{Name: "sync", Description: "Manually synchronize the local vault with Dashlane"},
			{Name: "read", Description: "Retrieve a secret from the local vault via its path"},
			{Name: "path", Description: "Path to the secret (dl://<title>/<field> or dl://<id>/<field>)"},
			{Name: "password", Description: "Retrieve a password from the local vault and copy it to the clipboard"},
			{Name: "otp", Description: "Retrieve an OTP code from local vault and copy it to the clipboard"},
			{Name: "note", Description: "Retrieve a secure note from the local vault and open it"},
			{Name: "accounts", Description: "Manage your accounts connected to the CLI"},
			{Name: "whoami", Description: "Prints the login email of the account currently in use"},
			{Name: "devices", Description: "Operations on devices"},
			{Name: "list", Description: "Lists all registered devices that can access your account"},
			{Name: "device ids", Description: "Ids of the devices to remove"},
			{Name: "register", Description: "Registers a new device to be used in non-interactive mode"},
			{Name: "device name", Description: "Name of the device to register"},
			{Name: "team", Description: "Team related commands"},
			{Name: "credentials", Description: "Team credentials operations"},
			{Name: "generate", Description: "Generate new team credentials"},
			{Name: "revoke", Description: "Revoke credentials by access key"},
			{Name: "accessKey", Description: "Access key of the credentials to revoke"},
			{Name: "members", Description: "List team members"},
			{Name: "page", Description: "Page number"},
			{Name: "limit", Description: "Limit of members per page"},
			{Name: "logs", Description: "List audit logs"},
			{Name: "report", Description: "Get team report"},
			{Name: "days", Description: "Number of days in history"},
			{Name: "configure", Description: "Configure the CLI"},
			{Name: "disable-auto-sync", Description: "Disable automatic synchronization which is done once per hour (default: false)"},
			{Name: "backup", Description: "Backup your local vault (will use the current directory by default)"},
			{Name: "logout", Description: "Logout and clean your local database and OS keychain"},
		},
		Options: []spec.Option{
			{Name: "-i", Description: "Input file of a template to inject the credential into"},
			{Name: "-o", Description: "Output file to write the injected template to"},
			{Name: "--print", Description: "Prints just the OTP code, instead of copying it to the clipboard"},
			{Name: "--json", Description: "Output in JSON format"},
			{Name: "--all", Description: "Remove all devices including this one (dangerous)"},
			{Name: "--others", Description: "Remove all other devices"},
			{Name: "--csv", Description: "Output in CSV format"},
			{Name: "--human-readable", Description: "Output dates in human readable format"},
			{Name: "--start", Description: "Start timestamp in ms"},
			{Name: "--end", Description: "Log type"},
			{Name: "--category", Description: "Log category"},
			{Name: "--directory", Description: "Logout and clean your local database and OS keychain"},
			{Name: "-h", Description: "Display help for command"},
			{Name: "-V", Description: "Output the version number"},
			{Name: "--debug", Description: "Print debug messages"},
		},
	})
}
