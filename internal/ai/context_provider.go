package ai

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/versenilvis/iris/internal/workspace"
)

type CommandContextProvider struct {
	NameStr   string
	Prefixes  []string
	GatherCmd []string
	Label     string
}

func (p *CommandContextProvider) Name() string { return p.NameStr }

func (p *CommandContextProvider) Matches(buf string) bool {
	trimmed := strings.ToLower(strings.TrimSpace(buf))
	for _, prefix := range p.Prefixes {
		if strings.HasPrefix(trimmed, prefix) {
			return true
		}
	}
	return false
}

func (p *CommandContextProvider) Gather(ctx context.Context) (string, error) {
	if len(p.GatherCmd) == 0 {
		return "", nil
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	out, err := exec.CommandContext(ctxTimeout, p.GatherCmd[0], p.GatherCmd[1:]...).Output()
	if err != nil {
		return "", err
	}
	if s := strings.TrimSpace(string(out)); s != "" {
		// Cap gathered command output to 1000 characters to keep prompt concise and avoid blowing up token budget
		if len(s) > 1000 {
			s = s[:1000] + "\n... (truncated)"
		}
		return p.Label + ":\n" + s, nil
	}
	return "", nil
}

var allowedHelpCommands = map[string]bool{
	"git": true, "docker": true, "kubectl": true, "npm": true, "yarn": true,
	"pnpm": true, "cargo": true, "go": true, "systemctl": true, "helm": true,
	"terraform": true, "aws": true, "gcloud": true, "az": true, "make": true,
	"bun": true, "pip": true, "python": true, "python3": true, "node": true,
	"deno": true, "tar": true, "curl": true, "wget": true, "ssh": true,
	"podman": true, "tofu": true, "ansible": true, "gh": true, "nix": true,
}

func isAllowedForHelp(cmdName string) bool {
	if strings.ContainsAny(cmdName, "/\\") {
		return false
	}
	return allowedHelpCommands[cmdName]
}

type universalProvider struct {
	cwd string
	buf string
}

func (p *universalProvider) Name() string {
	firstWord := ""
	if fields := strings.Fields(p.buf); len(fields) > 0 {
		firstWord = fields[0]
	}
	return "universal:" + p.cwd + ":" + firstWord
}

func (p *universalProvider) Matches(buf string) bool {
	return true
}

func (p *universalProvider) Gather(ctx context.Context) (string, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 1200*time.Millisecond)
	defer cancel()

	var sb strings.Builder

	ws := workspace.DetectCached(p.cwd)
	var ecosystems []string
	if ws.HasGit {
		ecosystems = append(ecosystems, "Git")
	}
	if ws.HasNodeProject {
		ecosystems = append(ecosystems, "Node/Bun")
	}
	if ws.HasGoProject {
		ecosystems = append(ecosystems, "Go")
	}
	if ws.HasRustProject {
		ecosystems = append(ecosystems, "Rust")
	}
	if ws.HasPythonProject {
		ecosystems = append(ecosystems, "Python")
	}
	if ws.HasJustfile {
		ecosystems = append(ecosystems, "Just")
	}
	if ws.HasMakefile {
		ecosystems = append(ecosystems, "Makefile/C++")
	}
	if ws.HasDockerfile {
		ecosystems = append(ecosystems, "Docker")
	}
	if ws.HasK8s {
		ecosystems = append(ecosystems, "K8s")
	}
	if len(ecosystems) > 0 {
		fmt.Fprintf(&sb, "Detected Workspace Ecosystems: %s\n\n", strings.Join(ecosystems, ", "))
	}

	ExtractScriptsAndTargets(&sb, p.cwd, "")

	if entries, err := os.ReadDir(p.cwd); err == nil {
		var names []string
		for i, e := range entries {
			if i >= 30 {
				names = append(names, "...")
				break
			}
			name := e.Name()
			if e.IsDir() {
				name += "/"
				if !strings.HasPrefix(e.Name(), ".") && e.Name() != "node_modules" && i < 15 {
					ExtractScriptsAndTargets(&sb, filepath.Join(p.cwd, e.Name()), e.Name())
				}
			}
			names = append(names, name)
		}
		if len(names) > 0 {
			fmt.Fprintf(&sb, "Files in Cwd: %s\n\n", strings.Join(names, ", "))
		}
	}

	if ws.HasGit {
		var wg sync.WaitGroup
		var branchErr, prevErr, recentErr, statusErr, diffErr, logErr error
		var branchOut, prevOut, recentOut, statusOut, diffOut, logOut []byte

		wg.Add(6)
		go func() {
			defer wg.Done()
			ctxProbe, cancel := context.WithTimeout(ctx, 400*time.Millisecond)
			defer cancel()
			branchOut, branchErr = exec.CommandContext(ctxProbe, "git", "-C", p.cwd, "rev-parse", "--abbrev-ref", "HEAD").Output()
		}()
		go func() {
			defer wg.Done()
			ctxProbe, cancel := context.WithTimeout(ctx, 400*time.Millisecond)
			defer cancel()
			prevOut, prevErr = exec.CommandContext(ctxProbe, "git", "-C", p.cwd, "rev-parse", "--abbrev-ref", "@{-1}").Output()
		}()
		go func() {
			defer wg.Done()
			ctxProbe, cancel := context.WithTimeout(ctx, 400*time.Millisecond)
			defer cancel()
			recentOut, recentErr = exec.CommandContext(ctxProbe, "git", "-C", p.cwd, "for-each-ref", "--sort=-committerdate", "--format=%(refname:short)", "--count=10", "refs/heads/").Output()
		}()
		go func() {
			defer wg.Done()
			statusOut, statusErr = exec.CommandContext(ctxTimeout, "git", "-C", p.cwd, "status", "-s").Output()
		}()
		go func() {
			defer wg.Done()
			diffOut, diffErr = exec.CommandContext(ctxTimeout, "git", "-C", p.cwd, "diff", "--staged").Output()
		}()
		go func() {
			defer wg.Done()
			logOut, logErr = exec.CommandContext(ctxTimeout, "git", "-C", p.cwd, "log", "-n", "5", "--no-decorate", "--pretty=format:%s").Output()
		}()
		wg.Wait()

		formatGitErr := func(name string, err error) string {
			if err == context.DeadlineExceeded || strings.Contains(err.Error(), "signal: killed") {
				return fmt.Sprintf("[Git probe timed out: %s]\n", name)
			}
			return fmt.Sprintf("[Git probe failed (%s): %v]\n", name, err)
		}

		sb.WriteString("Git Repository State:\n")
		if branchErr != nil {
			sb.WriteString(formatGitErr("current branch", branchErr))
		} else if currentBranch := strings.TrimSpace(string(branchOut)); currentBranch != "" {
			fmt.Fprintf(&sb, "Current Branch: %s\n", currentBranch)
		}

		if prevErr != nil {
			if !strings.Contains(prevErr.Error(), "exit status 128") {
				sb.WriteString(formatGitErr("previous branch", prevErr))
			}
		} else if prevBranch := strings.TrimSpace(string(prevOut)); prevBranch != "" && prevBranch != "HEAD" && prevBranch != strings.TrimSpace(string(branchOut)) {
			fmt.Fprintf(&sb, "Previous Checkout/Switch Branch (@{-1}): %s\n", prevBranch)
		}

		if recentErr != nil {
			sb.WriteString(formatGitErr("recent branches", recentErr))
		} else {
			recentBranchesList := strings.Split(strings.TrimSpace(string(recentOut)), "\n")
			var recentBranches []string
			for _, b := range recentBranchesList {
				b = strings.TrimSpace(b)
				if b != "" {
					recentBranches = append(recentBranches, b)
				}
			}
			if len(recentBranches) > 0 {
				fmt.Fprintf(&sb, "Recent Local Branches (by recent activity): %s\n\n", strings.Join(recentBranches, ", "))
			} else if strings.TrimSpace(string(branchOut)) != "" && branchErr == nil {
				sb.WriteString("\n")
			}
		}

		if statusErr != nil {
			sb.WriteString(formatGitErr("status", statusErr))
		} else if statusStr := strings.TrimSpace(string(statusOut)); statusStr != "" {
			if len(statusStr) > 1000 {
				statusStr = statusStr[:1000] + "\n... (truncated)"
			}
			fmt.Fprintf(&sb, "Status:\n%s\n\n", statusStr)
		}

		if diffErr != nil {
			sb.WriteString(formatGitErr("staged diff", diffErr))
		} else if diffStr := strings.TrimSpace(string(diffOut)); diffStr != "" {
			if len(diffStr) > 1500 {
				diffStr = diffStr[:1500] + "\n... (truncated)"
			}
			fmt.Fprintf(&sb, "Staged Diff:\n%s\n\n", diffStr)
		}

		if logErr != nil {
			sb.WriteString(formatGitErr("commit messages", logErr))
		} else if logStr := strings.TrimSpace(string(logOut)); logStr != "" {
			fmt.Fprintf(&sb, "User's recent commit messages (MUST follow this exact style, formatting, language, and casing conventions):\n%s\n", logStr)
		}
	}

	if fields := strings.Fields(p.buf); len(fields) > 0 {
		cmdName := fields[0]
		if isAllowedForHelp(cmdName) {
			ctxHelp, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
			defer cancel()
			helpOut, err := exec.CommandContext(ctxHelp, cmdName, "--help").CombinedOutput()
			if err == nil {
				helpStr := strings.TrimSpace(string(helpOut))
				if len(helpStr) > 600 {
					helpStr = helpStr[:600]
				}
				if helpStr != "" {
					fmt.Fprintf(&sb, "\nCommand help (%s --help):\n%s\n", cmdName, helpStr)
				}
			}
		}
	}

	return strings.TrimSpace(sb.String()), nil
}
