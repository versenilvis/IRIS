package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "turbo",
		Description: "Print the version",
		Subcommands: []core.Subcommand{
			{Name: "bin", Description: "Get the path to the turbo binary"},
			{Name: "link", Description: "Link your local directory to a Vercel organization and enable remote caching"},
			{Name: "login", Description: "Login to your Vercel account"},
			{Name: "logout", Description: "Logout of your Vercel account"},
			{Name: "prune", Description: "Prepare a subset of your monorepo"},
			{Name: "run", Description: "Run tasks in your monorepo"},
			{Name: "new-only", Description: "Only new output with hashes for cached tasks"},
			{Name: "hash-only", Description: "Only turbo-computed task hashes"},
			{Name: "full", Description: "Show all output"},
			{Name: "none", Description: "Hide process output"},
		},
		Options: []core.Option{
			{Name: "--version", Description: "Print the version"},
			{Name: "--help", Description: "Print a help message"},
			{Name: "--no-gitignore", Description: "Do not create or modify .gitignore"},
			{Name: "--sso-team", Description: "Attempt to authenticate to the specified team using SSO"},
			{Name: "--scope", Description: "Specify package to act as entry point for pruned monorepo"},
			{Name: "--docker", Description: "Run tasks in your monorepo"},
			{Name: "--cache-dir", Description: "Specify local filesystem cache directory"},
			{Name: "--concurrency", Description: "Limit the concurrency of task execution (use `1` for serial)"},
			{Name: "--continue", Description: "Continue execution even if a task exits with an error"},
			{Name: "--force", Description: "Ignore the existing cache"},
			{Name: "--graph", Description: "Generate a Dot graph of the task execution"},
			{Name: "--global-deps", Description: "Limit/set scope to changed packages since a mergebase"},
			{Name: "--team", Description: "The slug of a turborepo.com team"},
			{Name: "--token", Description: "A turborepo.com access token"},
			{Name: "--ignore", Description: "Files to ignore when calculating changed files (supports globs)"},
			{Name: "--profile", Description: "File to write turbo's performance profile into"},
			{Name: "--parallel", Description: "Execute all tasks in parallel"},
			{Name: "--include-dependencies", Description: "Include the dependencies of tasks in execution"},
			{Name: "--no-deps", Description: "Exclude dependent task consumers from execution"},
			{Name: "--no-cache", Description: "Avoid saving task results to the cache (useful for development/watch tasks)"},
			{Name: "--output-logs", Description: "Only new output with hashes for cached tasks"},
			{Name: "--dry", Description: "List the packages in scope and the tasks that would be run"},
		},
	})
}
