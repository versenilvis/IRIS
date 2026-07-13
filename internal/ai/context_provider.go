package ai

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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

	cmd := exec.CommandContext(ctxTimeout, "git", "-C", p.cwd, "rev-parse", "--is-inside-work-tree")
	if cmd.Run() == nil {
		branchOut, _ := exec.CommandContext(ctxTimeout, "git", "-C", p.cwd, "rev-parse", "--abbrev-ref", "HEAD").Output()
		currentBranch := strings.TrimSpace(string(branchOut))

		prevBranchOut, _ := exec.CommandContext(ctxTimeout, "git", "-C", p.cwd, "rev-parse", "--abbrev-ref", "@{-1}").Output()
		prevBranch := strings.TrimSpace(string(prevBranchOut))

		recentBranchesOut, _ := exec.CommandContext(ctxTimeout, "git", "-C", p.cwd, "branch", "--sort=-committerdate", "--format=%(refname:short)").Output()
		recentBranchesList := strings.Split(strings.TrimSpace(string(recentBranchesOut)), "\n")
		var recentBranches []string
		for i, b := range recentBranchesList {
			b = strings.TrimSpace(b)
			if b != "" && i < 10 {
				recentBranches = append(recentBranches, b)
			}
		}

		statusOut, _ := exec.CommandContext(ctxTimeout, "git", "-C", p.cwd, "status", "-s").Output()
		statusStr := strings.TrimSpace(string(statusOut))
		if len(statusStr) > 1000 {
			statusStr = statusStr[:1000] + "\n... (truncated)"
		}

		diffOut, _ := exec.CommandContext(ctxTimeout, "git", "-C", p.cwd, "diff", "--staged").Output()
		diffStr := strings.TrimSpace(string(diffOut))
		if len(diffStr) > 1500 {
			diffStr = diffStr[:1500] + "\n... (truncated)"
		}

		logOut, _ := exec.CommandContext(ctxTimeout, "git", "-C", p.cwd, "log", "-n", "5", "--no-decorate", "--pretty=format:%s").Output()
		logStr := strings.TrimSpace(string(logOut))

		sb.WriteString("Git Repository State:\n")
		if currentBranch != "" {
			fmt.Fprintf(&sb, "Current Branch: %s\n", currentBranch)
		}
		if prevBranch != "" && prevBranch != "HEAD" && prevBranch != currentBranch {
			fmt.Fprintf(&sb, "Previous Checkout/Switch Branch (@{-1}): %s\n", prevBranch)
		}
		if len(recentBranches) > 0 {
			fmt.Fprintf(&sb, "Recent Local Branches (by recent activity): %s\n\n", strings.Join(recentBranches, ", "))
		} else if currentBranch != "" {
			sb.WriteString("\n")
		}
		if statusStr != "" {
			fmt.Fprintf(&sb, "Status:\n%s\n\n", statusStr)
		}
		if diffStr != "" {
			fmt.Fprintf(&sb, "Staged Diff:\n%s\n\n", diffStr)
		}
		if logStr != "" {
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
