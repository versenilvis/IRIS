package commands

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "git",
		Description: "version control",
		Options: []core.Option{
			{Name: "--version", Description: "print version"},
			{Name: "--help", Description: "show help"},
		},
		Subcommands: []core.Subcommand{
			{
				Name: "init",
				Description: "create empty repo",
				Options: []core.Option{
					{Name: "--bare", Description: "create bare repo"},
					{Name: "-b", Description: "initial branch name"},
				},
			},
			{
				Name: "clone",
				Description: "clone a repository",
				Options: []core.Option{
					{Name: "--depth", Description: "shallow clone depth"},
					{Name: "--branch", Description: "specific branch"},
					{Name: "--bare", Description: "clone as bare"},
				},
			},
			{
				Name: "status",
				Description: "show working tree",
				Options: []core.Option{
					{Name: "-s", Description: "short format"},
					{Name: "--porcelain", Description: "machine format"},
				},
			},
			{
				Name: "add",
				Description: "stage changes",
				Options: []core.Option{
					{Name: "-A", Description: "add all files"},
					{Name: "-p", Description: "interactive patch"},
					{Name: ".", Description: "add current dir"},
				},
			},
			{
				Name: "commit",
				Description: "record changes",
				Options: []core.Option{
					{Name: "-m", Description: "commit message"},
					{Name: "-a", Description: "auto stage tracked"},
					{Name: "--amend", Description: "amend last commit"},
					{Name: "--no-edit", Description: "keep message"},
				},
			},
			{
				Name: "push",
				Description: "update remote refs",
				Options: []core.Option{
					{Name: "-u", Description: "set upstream"},
					{Name: "--force", Description: "force push"},
					{Name: "--tags", Description: "push tags"},
					{Name: "origin", Description: "default remote"},
				},
			},
			{
				Name: "pull",
				Description: "fetch and merge",
				Options: []core.Option{
					{Name: "--rebase", Description: "rebase on pull"},
					{Name: "origin", Description: "default remote"},
				},
			},
			{
				Name: "fetch",
				Description: "download objects",
				Options: []core.Option{
					{Name: "--all", Description: "fetch all remotes"},
					{Name: "--prune", Description: "remove stale refs"},
				},
			},
			{
				Name: "checkout",
				Description: "switch branches",
				Options: []core.Option{
					{Name: "-b", Description: "create new branch"},
				},
			},
			{
				Name: "switch",
				Description: "switch branches",
				Options: []core.Option{
					{Name: "-c", Description: "create and switch"},
				},
			},
			{
				Name: "branch",
				Description: "manage branches",
				Options: []core.Option{
					{Name: "-d", Description: "delete branch"},
					{Name: "-D", Description: "force delete"},
					{Name: "-a", Description: "list all"},
					{Name: "-m", Description: "rename branch"},
				},
			},
			{
				Name: "merge",
				Description: "join branches",
				Options: []core.Option{
					{Name: "--no-ff", Description: "no fast forward"},
					{Name: "--squash", Description: "squash commits"},
					{Name: "--abort", Description: "abort merge"},
				},
			},
			{
				Name: "rebase",
				Description: "reapply commits",
				Options: []core.Option{
					{Name: "-i", Description: "interactive"},
					{Name: "--onto", Description: "rebase onto"},
					{Name: "--abort", Description: "abort rebase"},
					{Name: "--continue", Description: "continue rebase"},
				},
			},
			{
				Name: "log",
				Description: "show commit log",
				Options: []core.Option{
					{Name: "--oneline", Description: "compact format"},
					{Name: "--graph", Description: "show graph"},
					{Name: "-n", Description: "limit count"},
				},
			},
			{
				Name: "diff",
				Description: "show changes",
				Options: []core.Option{
					{Name: "--staged", Description: "staged changes"},
					{Name: "--stat", Description: "diffstat only"},
				},
			},
			{
				Name: "stash",
				Description: "stash changes",
				Subcommands: []core.Subcommand{
					{Name: "pop", Description: "apply and drop"},
					{Name: "apply", Description: "apply stash"},
					{Name: "drop", Description: "remove stash"},
					{Name: "list", Description: "list stashes"},
					{Name: "show", Description: "show stash diff"},
				},
			},
			{
				Name: "reset",
				Description: "reset HEAD",
				Options: []core.Option{
					{Name: "--hard", Description: "discard changes"},
					{Name: "--soft", Description: "keep staged"},
					{Name: "--mixed", Description: "unstage changes"},
				},
			},
			{
				Name: "tag",
				Description: "manage tags",
				Options: []core.Option{
					{Name: "-a", Description: "annotated tag"},
					{Name: "-d", Description: "delete tag"},
					{Name: "-l", Description: "list tags"},
				},
			},
			{
				Name: "remote",
				Description: "manage remotes",
				Subcommands: []core.Subcommand{
					{Name: "add", Description: "add remote"},
					{Name: "remove", Description: "remove remote"},
					{Name: "rename", Description: "rename remote"},
					{Name: "-v", Description: "verbose list"},
				},
			},
			{
				Name: "cherry-pick",
				Description: "apply commit",
				Options: []core.Option{
					{Name: "--no-commit", Description: "no auto commit"},
					{Name: "--abort", Description: "abort pick"},
				},
			},
			{
				Name: "bisect",
				Description: "binary search bug",
				Subcommands: []core.Subcommand{
					{Name: "start", Description: "start bisect"},
					{Name: "good", Description: "mark good"},
					{Name: "bad", Description: "mark bad"},
					{Name: "reset", Description: "end bisect"},
				},
			},
		},
	})
}
