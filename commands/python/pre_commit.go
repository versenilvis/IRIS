package python

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "pre-commit",
		Description: "Show help message and exit",
		Subcommands: []core.Subcommand{
			{Name: "autoupdate", Description: "Auto-update pre-commit config to the latest repos' versions"},
			{Name: "clean", Description: "Clean out pre-commit files"},
			{Name: "gc", Description: "Clean unused cached repos"},
			{Name: "install", Description: "Install the pre-commit script"},
			{Name: "install-hooks", Description: "Whether to install hook environments for all environments in the config file"},
			{Name: "migrate-config", Description: "Migrate list configuration to new map configuration"},
			{Name: "run", Description: "Run hooks"},
			{Name: "hook", Description: "A single hook-id to run"},
			{Name: "sample-config", Description: "Produce a sample .pre-commit-config.yaml file"},
			{Name: "try-repo", Description: "Try the hooks in a repository, useful for developing new hooks"},
			{Name: "repo", Description: "Repository to source hooks from"},
			{Name: "uninstall", Description: "Uninstall the pre-commit script"},
			{Name: "help", Description: "Show help for a specific command"},
		},
		Options: []core.Option{
			{Name: "-h", Description: "Show help message and exit"},
			{Name: "--color", Description: "Whether to use color in output. Defaults to `auto`"},
			{Name: "--config", Description: "Path to alternate config file"},
			{Name: "-t", Description: "Type of hook to install"},
			{Name: "--verbose", Description: "Run all files in the repo"},
			{Name: "--files", Description: "Specific filenames to run hooks on"},
			{Name: "--show-diff-on-failure", Description: "When hooks fail, run `git diff` directly afterward"},
			{Name: "--hook-stage", Description: "The stage during which the hook is fired"},
			{Name: "--remote-branch", Description: "Remote branch ref used by `git push`"},
			{Name: "--local-branch", Description: "Local branch ref used by `git push`"},
			{Name: "--from-ref", Description: "Filename to check when running during `commit-msg`"},
			{Name: "--remote-name", Description: "Remote name used by `git push`"},
			{Name: "--remote-url", Description: "Remote URL used by `git push`"},
			{Name: "--checkout-type", Description: "During a post-merge hook, indicates whether the merge was a squash merge"},
			{Name: "--rewrite-command", Description: "During a post-rewrite hook, specifies the command that invoked the rewrite"},
			{Name: "-V", Description: "Show program's version number and exit"},
			{Name: "--bleeding-edge", Description: "Store 'frozen' hashes in `rev` instead of tag names"},
			{Name: "--repo", Description: "Only update this repository -- may be specified multiple times"},
			{Name: "--no-allow-missing-config", Description: "Assume cloned repos should have a `pre-commit` config"},
			{Name: "-f", Description: "Overwrite existing hooks / remove migration mode"},
			{Name: "--allow-missing-config", Description: "Migrate list configuration to new map configuration"},
			{Name: "--ref", Description: "Manually select a rev to run against, otherwise the `HEAD` revision will be used"},
		},
	})
}
