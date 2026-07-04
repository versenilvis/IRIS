package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "trex",
		Description: "trex script",
		Subcommands: []core.Subcommand{
			{Name: "i", Description: "Install a package"},
			{Name: "delete", Description: "Delete a package"},
			{Name: "upgrade", Description: "Upgrade trex"},
			{Name: "tree", Description: "View dependency tree"},
			{Name: "run", Description: "Run a script alias in a file run.json"},
			{Name: "purge", Description: "Remove a package or url from cache"},
			{Name: "ls", Description: "Shows the list of installed packages"},
			{Name: "exec", Description: "Execute a cli tool with out install then"},
			{Name: "check", Description: "Check deno.land [std/x] dependencies updates"},
		},
		Options: []core.Option{
			{Name: "-v", Description: "Print version"},
			{Name: "-m", Description: "Install package from deno.land"},
			{Name: "-n", Description: "Install package from nest.land"},
			{Name: "-p", Description: "Install package from some repository"},
			{Name: "-c", Description: "Install custom package"},
			{Name: "--canary", Description: "Install from dev branch"},
			{Name: "-w", Description: "Use reboot script alias protocol (rsap)"},
			{Name: "-wv", Description: "Verbose output in --watch mode (rsap)"},
			{Name: "--perms", Description: "Specify cli permissions"},
			{Name: "-f", Description: "Update outdated dependencies"},
			{Name: "-h", Description: "Print help info"},
		},
	})
}
