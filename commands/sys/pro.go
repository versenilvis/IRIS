package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "pro",
		Description: "Manage Ubuntu Pro services from Canonical",
		Subcommands: []spec.Subcommand{
			{Name: "attach", Description: "Connect an Ubuntu Pro support contract to this machine"},
			{Name: "collect-logs", Description: "Create a tarball with all relevant logs and debug data"},
			{Name: "detach", Description: "Remove Ubuntu Pro from this machine"},
			{Name: "disable", Description: "Disable this machine's access to an Ubuntu Pro service"},
			{Name: "enable", Description: "Activate and configure this machine's access to an Ubuntu Pro service"},
			{Name: "fix", Description: "Fix a CVE or USN on the  system  by  upgrading  the  appropriate package(s)"},
			{Name: "refresh", Description: "Refresh contract and service details from Canonical"},
			{Name: "status", Description: "Report current status of Ubuntu Pro services on system"},
		},
		Options: []spec.Option{
			{Name: "--no-auto-enable", Description: "Disable  the  automatic enablement of recommended entitlements"},
			{Name: "--attach-config", Description: "Provide a file with configuration options"},
			{Name: "-o", Description: "Path for tarball. Uses ua_logs.tar.gz in current directory if not specified"},
			{Name: "--format", Description: "Output format"},
			{Name: "--simulate-with-token", Description: "Include beta and unavailable services"},
		},
	})
}
