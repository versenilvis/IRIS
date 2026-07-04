package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "tsh",
		Description: "Remote host login",
		Subcommands: []core.Subcommand{
			{Name: "version", Description: "Print the version"},
			{Name: "ssh", Description: "Run shell or execute a command on a remote SSH node"},
			{Name: "user@hostname", Description: "Address of remote machine to log into"},
			{Name: "aws", Description: "Access AWS API"},
			{Name: "apps", Description: "View and control proxied applications"},
			{Name: "ls", Description: "List available applications"},
			{Name: "logout", Description: "Remove app certificate"},
			{Name: "config", Description: "Print app connection information"},
			{Name: "login", Description: "Retrieve credentials for a database"},
			{Name: "connect", Description: "Connect   Connect to a database"},
			{Name: "join", Description: "Join the active SSH session"},
			{Name: "play", Description: "Replay the recorded SSH session"},
			{Name: "scp", Description: "Secure file copy"},
			{Name: "clusters", Description: "List available Teleport clusters"},
			{Name: "status", Description: "Display the list of proxy servers and retrieved certificates"},
			{Name: "env", Description: "Print commands to set Teleport session environment variables"},
			{Name: "request", Description: "Manage access requests"},
			{Name: "show", Description: "Show request details"},
			{Name: "new", Description: "Create a new access request"},
			{Name: "review", Description: "Review an access request"},
			{Name: "kube", Description: "Manage available kubernetes clusters"},
			{Name: "sessions", Description: "Get a list of active kubernetes sessions"},
			{Name: "exec", Description: "Execute a command in a kubernetes pod"},
			{Name: "mfa", Description: "Manage multi-factor authentication (MFA) devices"},
			{Name: "add", Description: "Add a new MFA device"},
			{Name: "rm", Description: "Remove a MFA device"},
		},
		Options: []core.Option{
			{Name: "-l", Description: "Remote host login"},
			{Name: "--proxy", Description: "SSH proxy address"},
			{Name: "--user", Description: "SSH proxy user"},
			{Name: "--ttl", Description: "Minutes to live for a SSH session"},
			{Name: "-i", Description: "Identity file"},
			{Name: "--cert-format", Description: "SSH certificate format"},
			{Name: "--insecure", Description: "Do not verify server's certificate and host name. Use only in test environments"},
			{Name: "--auth", Description: "Specify the name of authentication connector to use"},
			{Name: "--skip-version-check", Description: "Skip version checking between server and client"},
			{Name: "-d", Description: "Verbose logging to stdout"},
			{Name: "-k", Description: "Controls how keys are handled. Valid values are [auto no yes only]"},
			{Name: "--enable-escape-sequences", Description: "Override host:port used when opening a browser for cluster logins"},
			{Name: "-J", Description: "SSH jumphost"},
			{Name: "-f", Description: "Format output"},
			{Name: "-q", Description: "Quiet mode"},
		},
	})
}
