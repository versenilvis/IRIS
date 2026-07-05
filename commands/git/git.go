package git

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/versenilvis/iris/spec"
)

func GitRemoteGenerator(tokens []string, _ string, _ string) []spec.Suggestion {
	return getGitResults(tokens, "remote")
}

func GitStashGenerator(tokens []string, _ string, _ string) []spec.Suggestion {
	return getGitResults(tokens, "stash", "list", "--format=%gd: %gs")
}

func GitTagGenerator(tokens []string, _ string, _ string) []spec.Suggestion {
	return getGitResults(tokens, "tag", "-l")
}

func GitCommitGenerator(tokens []string, _ string, _ string) []spec.Suggestion {
	cwd := spec.GetCWD()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "log", "--format=%h [%cr] %s", "-30")
	cmd.Dir = cwd
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	var results []spec.Suggestion
	for line := range strings.SplitSeq(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		hash := parts[0]
		desc := ""
		if len(parts) == 2 {
			desc = parts[1]
		}
		results = append(results, spec.Suggestion{Cmd: hash, Desc: desc})
	}
	return results
}

func getGitResults(tokens []string, args ...string) []spec.Suggestion {
	return getGitResultsFiltered(tokens, false, args...)
}

func getGitResultsFiltered(tokens []string, localOnly bool, args ...string) []spec.Suggestion {
	cwd := spec.GetCWD()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = cwd
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	activeBranch := ""
	if args[0] == "branch" {
		activeCmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
		activeCmd.Dir = cwd
		if activeOut, err := activeCmd.Output(); err == nil {
			activeBranch = strings.TrimSpace(string(activeOut))
		}
	}

	seen := make(map[string]bool)
	lines := strings.Split(string(out), "\n")
	var results []spec.Suggestion

	// find the subcommand by skipping global flags that take arguments
	subcommand := ""
	for i := 1; i < len(tokens); i++ {
		t := tokens[i]
		if strings.HasPrefix(t, "-") {
			if t == "-c" || t == "-C" || t == "--git-dir" || t == "--work-tree" || t == "--namespace" || t == "--super-prefix" || t == "--config-env" || t == "--exec-path" {
				i++
			}
			continue
		}
		subcommand = t
		break
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		isRemote := strings.HasPrefix(line, "remotes/")

		// skip remote tracking branches if localOnly mode
		if localOnly && isRemote {
			continue
		}

		// strip remotes/ prefix to get the usable form: origin/main
		if isRemote {
			line = strings.TrimPrefix(line, "remotes/")
		}

		// only skip active branch for checkout/switch commands
		if subcommand == "checkout" || subcommand == "switch" {
			if line == activeBranch {
				continue
			}
			// also skip any remote branch that tracks the active branch (e.g. origin/main)
			if idx := strings.Index(line, "/"); isRemote && idx != -1 {
				if line[idx+1:] == activeBranch {
					continue
				}
			}
		}

		// dedup: origin/dev and dev would both appear with -a; skip if already seen the short name
		shortName := line
		if idx := strings.Index(line, "/"); isRemote && idx != -1 {
			shortName = line[idx+1:] // "origin/dev" -> "dev"
		}
		if seen[shortName] {
			continue
		}
		seen[shortName] = true

		suggestionCmd := line
		suggestionDesc := args[0]

		if args[0] == "stash" {
			parts := strings.SplitN(line, ": ", 2)
			if len(parts) == 2 {
				suggestionCmd = parts[0]
				suggestionDesc = parts[1]
			}
		}

		results = append(results, spec.Suggestion{
			Cmd:  suggestionCmd,
			Desc: suggestionDesc,
		})
	}

	if activeBranch != "" {
		for i, r := range results {
			if r.Cmd == activeBranch {
				copy(results[1:i+1], results[0:i])
				results[0] = r
				break
			}
		}
	}

	return results
}

// GitBranchGenerator suggests git branches (local + remote, deduped)
func GitBranchGenerator(tokens []string, _ string, _ string) []spec.Suggestion {
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

	return getGitResults(tokens, "branch", "-a", "--format=%(refname:short)")
}

// gitLocalBranchGenerator is like GitBranchGenerator but only local branches
// used for push/pull where remote tracking branches cause duplicates
func gitLocalBranchGenerator(tokens []string, _ string, _ string) []spec.Suggestion {
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
	return getGitResultsFiltered(tokens, true, "branch", "-a", "--format=%(refname:short)")
}

func GitPushPullGenerator(tokens []string, prefix string, partial string) []spec.Suggestion {
	// count completed positional args (not flags, not the partial being typed)
	// tokens[0] = "git", tokens[1] = "push"/"pull", so start at 2
	// exclude tokens[len-1] because that's the partial being typed
	pArgs := []string{}
	for i := 2; i < len(tokens)-1; i++ {
		t := tokens[i]
		if t == "" || strings.HasPrefix(t, "-") {
			continue
		}
		pArgs = append(pArgs, t)
	}

	// no remote confirmed yet, suggest remotes
	if len(pArgs) == 0 {
		return GitRemoteGenerator(tokens, prefix, partial)
	}

	// remote is set, suggest local branches only (no duplicates with origin/xxx)
	if len(pArgs) >= 2 {
		return nil
	}
	return gitLocalBranchGenerator(tokens, prefix, partial)
}

func init() {
	spec.Register(&spec.Spec{
		Name:        "git",
		Description: "version control",
		Options: []spec.Option{
			{Name: "--version", Description: "print version"},
			{Name: "--help", Description: "show help"},
		},
		Subcommands: []spec.Subcommand{
			{
				Name:        "init",
				Description: "create empty repo",
				Options: []spec.Option{
					{Name: "--bare", Description: "create bare repo"},
					{Name: "-b", Description: "initial branch name"},
				},
			},
			{
				Name:        "clone",
				Description: "clone a repository",
				Options: []spec.Option{
					{Name: "--depth", Description: "shallow clone depth"},
					{Name: "--branch", Description: "specific branch"},
					{Name: "--bare", Description: "clone as bare"},
				},
			},
			{
				Name:        "status",
				Description: "show working tree",
				Options: []spec.Option{
					{Name: "-s", Description: "short format"},
					{Name: "--porcelain", Description: "machine format"},
				},
			},
			{
				Name:        "add",
				Description: "stage changes",
				Generator:   spec.FileGenerator(),
				Options: []spec.Option{
					{Name: "-A", Description: "add all files"},
					{Name: "-p", Description: "interactive patch"},
					{Name: ".", Description: "add current dir"},
				},
			},
			{
				Name:        "commit",
				Description: "record changes",
				Options: []spec.Option{
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
				Options: []spec.Option{
					{Name: "-u", Description: "set upstream"},
					{Name: "--force", Description: "force push"},
					{Name: "--tags", Description: "push tags"},
				},
			},
			{
				Name:        "pull",
				Description: "fetch and merge",
				Generator:   GitPushPullGenerator,
				Options: []spec.Option{
					{Name: "--rebase", Description: "rebase on pull"},
				},
			},
			{
				Name:        "fetch",
				Description: "download objects",
				Generator:   GitRemoteGenerator,
				Options: []spec.Option{
					{Name: "--all", Description: "fetch all remotes"},
					{Name: "--prune", Description: "remove stale refs"},
				},
			},
			{
				Name:        "checkout",
				Description: "switch branches",
				Generator: func(tokens []string, prefix string, partial string) []spec.Suggestion {
					for _, t := range tokens {
						if t == "-b" || t == "-B" {
							return nil
						}
					}
					branches := GitBranchGenerator(tokens, prefix, partial)
					files := spec.FileGenerator()(tokens, prefix, partial)
					return append(branches, files...)
				},
				Options: []spec.Option{
					{Name: "-b", Description: "create new branch"},
				},
			},
			{
				Name:        "switch",
				Description: "switch branches",
				Generator:   GitBranchGenerator,
				Options: []spec.Option{
					{Name: "-c", Description: "create and switch"},
				},
			},
			{
				Name:        "branch",
				Description: "manage branches",
				Generator:   GitBranchGenerator,
				Options: []spec.Option{
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
				Options: []spec.Option{
					{Name: "--no-ff", Description: "no fast forward"},
					{Name: "--squash", Description: "squash commits"},
					{Name: "--abort", Description: "abort merge"},
				},
			},
			{
				Name:        "rebase",
				Description: "reapply commits",
				Generator:   GitBranchGenerator,
				Options: []spec.Option{
					{Name: "-i", Description: "interactive"},
					{Name: "--onto", Description: "rebase onto"},
					{Name: "--abort", Description: "abort rebase"},
					{Name: "--continue", Description: "continue rebase"},
				},
			},
			{
				Name:        "log",
				Description: "show commit log",
				Generator:   spec.FileGenerator(),
				Options: []spec.Option{
					{Name: "--oneline", Description: "compact format"},
					{Name: "--graph", Description: "show graph"},
					{Name: "-n", Description: "limit count"},
				},
			},
			{
				Name:        "diff",
				Description: "show changes",
				Generator:   spec.FileGenerator(),
				Options: []spec.Option{
					{Name: "--staged", Description: "staged changes"},
					{Name: "--stat", Description: "diffstat only"},
					{Name: "--", Description: "separate paths"},
				},
			},
			{
				Name:        "tag",
				Description: "manage tags",
				Generator:   GitTagGenerator,
				Options: []spec.Option{
					{Name: "-a", Description: "annotated tag"},
					{Name: "-d", Description: "delete tag"},
					{Name: "-l", Description: "list tags"},
					{Name: "--delete", Description: "delete tag"},
					{Name: "-m", Description: "tag message"},
				},
			},
			{
				Name:        "show",
				Description: "show object",
				Generator: func(tokens []string, prefix string, partial string) []spec.Suggestion {
					tags := GitTagGenerator(tokens, prefix, partial)
					commits := GitCommitGenerator(tokens, prefix, partial)
					return append(tags, commits...)
				},
				Options: []spec.Option{
					{Name: "--stat", Description: "diffstat only"},
					{Name: "--name-only", Description: "filenames only"},
				},
			},
			{
				Name:        "revert",
				Description: "revert a commit",
				Generator:   GitCommitGenerator,
				Options: []spec.Option{
					{Name: "--no-commit", Description: "no auto commit"},
					{Name: "--abort", Description: "abort revert"},
					{Name: "--continue", Description: "continue revert"},
				},
			},
			{
				Name:        "reset",
				Description: "reset HEAD",
				Generator: func(tokens []string, prefix string, partial string) []spec.Suggestion {
					branches := GitBranchGenerator(tokens, prefix, partial)
					files := spec.FileGenerator()(tokens, prefix, partial)
					return append(branches, files...)
				},
				Options: []spec.Option{
					{Name: "--hard", Description: "discard changes"},
					{Name: "--soft", Description: "keep staged"},
					{Name: "--mixed", Description: "unstage changes"},
				},
			},
			{
				Name:        "restore",
				Description: "restore working tree files",
				Generator:   spec.FileGenerator(),
				Options: []spec.Option{
					{Name: "-s", Description: "source tree"},
					{Name: "-W", Description: "working tree"},
				},
			},
			{
				Name:        "rm",
				Description: "remove files",
				Generator:   spec.FileGenerator(),
				Options: []spec.Option{
					{Name: "-f", Description: "force"},
					{Name: "-r", Description: "recursive"},
					{Name: "--cached", Description: "unstage only"},
				},
			},
			{
				Name:        "stash",
				Description: "stash changes",
				Subcommands: []spec.Subcommand{
					{Name: "pop", Description: "apply and drop", Generator: GitStashGenerator, Options: []spec.Option{{Name: "--index", Description: "try to reinstate index"}}},
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
				Subcommands: []spec.Subcommand{
					{Name: "add", Description: "add remote", Options: []spec.Option{{Name: "-f", Description: "fetch immediately"}}},
					{Name: "remove", Description: "remove remote", Generator: GitRemoteGenerator},
					{Name: "rename", Description: "rename remote", Generator: GitRemoteGenerator},
					{Name: "set-url", Description: "change remote url", Generator: GitRemoteGenerator},
					{Name: "-v", Description: "verbose list"},
				},
			},
			{
				Name:        "cherry-pick",
				Description: "apply commit",
				Generator:   GitCommitGenerator,
				Options: []spec.Option{
					{Name: "--no-commit", Description: "no auto commit"},
					{Name: "--abort", Description: "abort pick"},
					{Name: "--continue", Description: "continue pick"},
				},
			},
			{
				Name:        "worktree",
				Description: "manage worktrees",
				Subcommands: []spec.Subcommand{
					{Name: "add", Description: "add new worktree", Generator: spec.FileGenerator("/")},
					{Name: "list", Description: "list worktrees"},
					{Name: "remove", Description: "remove worktree"},
					{Name: "prune", Description: "prune stale worktrees"},
				},
			},
			{
				Name:        "submodule",
				Description: "manage submodules",
				Subcommands: []spec.Subcommand{
					{Name: "add", Description: "add submodule"},
					{Name: "init", Description: "init submodule config"},
					{Name: "update", Description: "update submodules", Options: []spec.Option{{Name: "--init", Description: "init if needed"}, {Name: "--recursive", Description: "recursive update"}}},
					{Name: "status", Description: "show submodule status"},
					{Name: "foreach", Description: "run command in each submodule"},
				},
			},
			{
				Name:        "bisect",
				Description: "binary search bug",
				Subcommands: []spec.Subcommand{
					{Name: "start", Description: "start bisect"},
					{Name: "good", Description: "mark good"},
					{Name: "bad", Description: "mark bad"},
					{Name: "reset", Description: "end bisect"},
				},
			},
		},
	})
}
