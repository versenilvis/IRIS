package dev

import (
	"context"
	"os/exec"
	"strings"

	"github.com/versenilvis/iris/commands/core"
)

// GitRemoteGenerator suggests git remotes
func GitRemoteGenerator(tokens []string, prefix string, partial string) []core.Suggestion {
	return getGitResults(prefix, "remote")
}

// GitStashGenerator suggests git stashes
func GitStashGenerator(tokens []string, prefix string, partial string) []core.Suggestion {
	return getGitResults(prefix, "stash", "list", "--format=%gd: %gs")
}

func getGitResults(prefix string, args ...string) []core.Suggestion {
	cwd := core.GetCWD()
	cmd := exec.CommandContext(context.Background(), "git", args...)
	cmd.Dir = cwd
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	activeBranch := ""
	switch args[0] {
	case "branch":
		// Try to find the current active branch to filter it out later
		activeCmd := exec.CommandContext(context.Background(), "git", "rev-parse", "--abbrev-ref", "HEAD")
		activeCmd.Dir = cwd
		if activeOut, err := activeCmd.Output(); err == nil {
			activeBranch = strings.TrimSpace(string(activeOut))
		}
	}

	lines := strings.Split(string(out), "\n")
	var results []core.Suggestion
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "*") { // skip active branch marker if any
			line = strings.TrimSpace(strings.TrimPrefix(line, "*"))
		}
		if line == "" || line == activeBranch {
			continue
		}
		
		// handle remote branches that look like "remotes/origin/main"
		line = strings.TrimPrefix(line, "remotes/")

		suggestionCmd := line
		suggestionDesc := args[0]

		// for stash list, the format is "stash@{0}: message"
		if args[0] == "stash" {
			parts := strings.SplitN(line, ": ", 2)
			if len(parts) == 2 {
				suggestionCmd = parts[0]
				suggestionDesc = parts[1]
			}
		}

		results = append(results, core.Suggestion{
			Cmd:  prefix + " " + suggestionCmd,
			Desc: suggestionDesc,
		})
	}
	return results
}

// GitBranchGenerator suggests git branches
func GitBranchGenerator(tokens []string, prefix string, partial string) []core.Suggestion {
	// check if we are in "create" mode (-b or -B or -c)
	isCreateMode := false
	for _, t := range tokens {
		if t == "-b" || t == "-B" || t == "-c" || t == "-C" {
			isCreateMode = true
			break
		}
	}

	if isCreateMode {
		return nil
	}

	return getGitResults(prefix, "branch", "-a", "--format=%(refname:short)")
}

// GitPushPullGenerator suggests remotes for the first arg, and branches for the second
func GitPushPullGenerator(tokens []string, prefix string, partial string) []core.Suggestion {
	// Filter out flags to find positional arguments
	args := []string{}
	for i := 1; i < len(tokens); i++ {
		t := tokens[i]
		if t != "" && !strings.HasPrefix(t, "-") {
			args = append(args, t)
		}
	}

	// args[0] is subcommand (push/pull)
	// args[1] should be remote
	// args[2] should be branch
	
	// If we only have subcommand, suggest remotes
	if len(args) == 1 {
		return GitRemoteGenerator(tokens, prefix, partial)
	}
	
	// If we have subcommand + remote, suggest branches
	if len(args) == 2 {
		return GitBranchGenerator(tokens, prefix, partial)
	}
	
	return nil
}

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
				Generator: core.FileGenerator(),
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
				Name:        "push",
				Description: "update remote refs",
				Generator:   GitPushPullGenerator,
				Options: []core.Option{
					{Name: "-u", Description: "set upstream"},
					{Name: "--force", Description: "force push"},
					{Name: "--tags", Description: "push tags"},
				},
			},
			{
				Name:        "pull",
				Description: "fetch and merge",
				Generator:   GitPushPullGenerator,
				Options: []core.Option{
					{Name: "--rebase", Description: "rebase on pull"},
				},
			},
			{
				Name:        "fetch",
				Description: "download objects",
				Generator:   GitRemoteGenerator,
				Options: []core.Option{
					{Name: "--all", Description: "fetch all remotes"},
					{Name: "--prune", Description: "remove stale refs"},
				},
			},
			{
				Name:        "checkout",
				Description: "switch branches",
				Generator: func(tokens []string, prefix string, partial string) []core.Suggestion {
					for _, t := range tokens {
						if t == "-b" || t == "-B" {
							return nil
						}
					}
					branches := GitBranchGenerator(tokens, prefix, partial)
					files := core.FileGenerator()(tokens, prefix, partial)
					return append(branches, files...)
				},
				Options: []core.Option{
					{Name: "-b", Description: "create new branch"},
				},
			},
			{
				Name:        "switch",
				Description: "switch branches",
				Generator:   GitBranchGenerator,
				Options: []core.Option{
					{Name: "-c", Description: "create and switch"},
				},
			},
			{
				Name:        "branch",
				Description: "manage branches",
				Generator:   GitBranchGenerator,
				Options: []core.Option{
					{Name: "-d", Description: "delete branch"},
					{Name: "-D", Description: "force delete"},
					{Name: "-a", Description: "list all"},
					{Name: "-m", Description: "rename branch"},
				},
			},
			{
				Name:        "merge",
				Description: "join branches",
				Generator:   GitBranchGenerator,
				Options: []core.Option{
					{Name: "--no-ff", Description: "no fast forward"},
					{Name: "--squash", Description: "squash commits"},
					{Name: "--abort", Description: "abort merge"},
				},
			},
			{
				Name:        "rebase",
				Description: "reapply commits",
				Generator:   GitBranchGenerator,
				Options: []core.Option{
					{Name: "-i", Description: "interactive"},
					{Name: "--onto", Description: "rebase onto"},
					{Name: "--abort", Description: "abort rebase"},
					{Name: "--continue", Description: "continue rebase"},
				},
			},
			{
				Name:        "log",
				Description: "show commit log",
				Generator:   core.FileGenerator(),
				Options: []core.Option{
					{Name: "--oneline", Description: "compact format"},
					{Name: "--graph", Description: "show graph"},
					{Name: "-n", Description: "limit count"},
				},
			},
			{
				Name:        "diff",
				Description: "show changes",
				Generator:   core.FileGenerator(),
				Options: []core.Option{
					{Name: "--staged", Description: "staged changes"},
					{Name: "--stat", Description: "diffstat only"},
					{Name: "--", Description: "separate paths"},
				},
			},
			{
				Name:        "tag",
				Description: "manage tags",
				Generator:   func(tokens []string, prefix string, partial string) []core.Suggestion { return getGitResults(prefix, "tag", "-l") },
				Options: []core.Option{
					{Name: "-a", Description: "annotated tag"},
					{Name: "-d", Description: "delete tag"},
					{Name: "-l", Description: "list tags"},
					{Name: "--delete", Description: "delete tag"},
					{Name: "-m", Description: "tag message"},
				},
			},
			{
				Name:        "reset",
				Description: "reset HEAD",
				Generator: func(tokens []string, prefix string, partial string) []core.Suggestion {
					branches := GitBranchGenerator(tokens, prefix, partial)
					files := core.FileGenerator()(tokens, prefix, partial)
					return append(branches, files...)
				},
				Options: []core.Option{
					{Name: "--hard", Description: "discard changes"},
					{Name: "--soft", Description: "keep staged"},
					{Name: "--mixed", Description: "unstage changes"},
				},
			},
			{
				Name:        "restore",
				Description: "restore working tree files",
				Generator:   core.FileGenerator(),
				Options: []core.Option{
					{Name: "-s", Description: "source tree"},
					{Name: "-W", Description: "working tree"},
				},
			},
			{
				Name:        "rm",
				Description: "remove files",
				Generator:   core.FileGenerator(),
				Options: []core.Option{
					{Name: "-f", Description: "force"},
					{Name: "-r", Description: "recursive"},
					{Name: "--cached", Description: "unstage only"},
				},
			},
			{
				Name:        "stash",
				Description: "stash changes",
				Subcommands: []core.Subcommand{
					{Name: "pop", Description: "apply and drop", Generator: GitStashGenerator, Options: []core.Option{{Name: "--index", Description: "try to reinstate index"}}},
					{Name: "apply", Description: "apply stash", Generator: GitStashGenerator},
					{Name: "drop", Description: "remove stash", Generator: GitStashGenerator},
					{Name: "list", Description: "list stashes"},
					{Name: "show", Description: "show stash diff", Generator: GitStashGenerator},
					{Name: "push", Description: "push to stash"},
					{Name: "branch", Description: "create branch from stash", Generator: GitBranchGenerator},
				},
			},
			{
				Name:        "remote",
				Description: "manage remotes",
				Subcommands: []core.Subcommand{
					{Name: "add", Description: "add remote", Options: []core.Option{{Name: "-f", Description: "fetch immediately"}}},
					{Name: "remove", Description: "remove remote", Generator: GitRemoteGenerator},
					{Name: "rename", Description: "rename remote", Generator: GitRemoteGenerator},
					{Name: "set-url", Description: "change remote url", Generator: GitRemoteGenerator},
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
