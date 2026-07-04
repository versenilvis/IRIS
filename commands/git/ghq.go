package git

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "ghq",
		Description: "Clone/sync with a remote repository",
		Subcommands: []core.Subcommand{
			{Name: "get", Description: "Clone/sync with a remote repository"},
			{Name: "list", Description: "List local repositories"},
			{Name: "root", Description: "Show repositories' root"},
			{Name: "create", Description: "Create a new repository"},
		},
		Options: []core.Option{
			{Name: "-u", Description: "Update local repository if cloned already (default: false)"},
			{Name: "-p", Description: "Clone with SSH (default: false)"},
			{Name: "--shallow", Description: "Do a shallow clone (default: false)"},
			{Name: "-l", Description: "Look after get (default: false)"},
			{Name: "--vcs", Description: "Specify vcs backend for cloning"},
			{Name: "-s", Description: "Clone or update silently (default: false)"},
			{Name: "--no-recursive", Description: "Prevent recursive fetching (default: false)"},
			{Name: "-b", Description: "Specify branch name. This flag implies --single-branch on Git"},
			{Name: "-P", Description: "Import parallelly (default: false)"},
			{Name: "--bare", Description: "Do a bare clone (default: false)"},
			{Name: "-e", Description: "Perform an exact match (default: false)"},
			{Name: "--unique", Description: "Print unique subpaths (default: false)"},
			{Name: "--all", Description: "Show all roots (default: false)"},
			{Name: "-h", Description: "Show help"},
			{Name: "-v", Description: "Print the version"},
		},
	})
}
