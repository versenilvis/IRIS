package rust

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "trunk",
		Description: "Run on all files instead of only changed files",
		Subcommands: []core.Subcommand{
			{Name: "init", Description: "Setup trunk in this repo"},
			{Name: "check", Description: "Universal code checker"},
			{Name: "upgrade", Description: "Upgrade all linters to latest versions"},
			{Name: "linters", Description: "Linter(s) to upgrade (upgrades all if none specified)"},
			{Name: "download", Description: "Download all files needed for trunk to work offline"},
			{Name: "tools", Description: "Tool(s) to download (if omitted, downloads all configured tools)"},
			{Name: "enable", Description: "Enable linters"},
			{Name: "disable", Description: "Disable linters"},
			{Name: "fmt", Description: "Universal code formatter"},
			{Name: "git-hooks", Description: "Git hooks"},
			{Name: "install", Description: "Install trunk git hooks"},
			{Name: "cache", Description: "Cache management"},
			{Name: "clean", Description: "Clean the cache"},
			{Name: "print-config", Description: "Print the resolved trunk config"},
			{Name: "daemon", Description: "Daemon management"},
			{Name: "launch", Description: "Start the trunk daemon if its not already running"},
			{Name: "shutdown", Description: "Shutdown the trunk daemon if it is running"},
			{Name: "status", Description: "Report daemon status"},
		},
		Options: []core.Option{
			{Name: "-a", Description: "Run on all files instead of only changed files"},
			{Name: "-n", Description: "Don't automatically apply fixes"},
			{Name: "--include-existing-autofixes", Description: "Show autofixes for existing issues"},
			{Name: "--force", Description: "Run on all files, even if ignored"},
			{Name: "--diff", Description: "Diff printing mode"},
			{Name: "--filter", Description: "Shorthand for an inverse --filter"},
			{Name: "-j", Description: "Number of concurrent jobs (does not affect background linting)"},
			{Name: "--sample", Description: "Run each linter on N files (implies --no-fix and --all if no paths are given)"},
			{Name: "--upstream", Description: "Upstream branch used to compute changed files (autodetected by default)"},
			{Name: "--lock", Description: "Add sha256s to trunk.yaml for additional verification"},
			{Name: "--check-sample", Description: "Run `trunk check sample` without prompting"},
			{Name: "--nocheck-sample", Description: "Do not run `trunk check sample` post-init"},
			{Name: "-y", Description: "Automatically apply all fixes without prompting"},
			{Name: "--all", Description: "Delete all files (including results cache)"},
			{Name: "-h", Description: "Usage information"},
			{Name: "--version", Description: "The version"},
			{Name: "-m", Description: "Enable the trunk daemon to monitor file changes in your repo"},
			{Name: "--ci", Description: "Run in continuous integration mode"},
			{Name: "-o", Description: "Output format"},
			{Name: "--no-progress", Description: "Don't show progress updates"},
			{Name: "--ci-progress", Description: "Output details about what's happening under the hood"},
		},
	})
}
