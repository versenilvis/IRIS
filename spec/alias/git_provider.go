package alias

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type GitProvider struct {
	cacheKey  string
	cached    []AliasEntry
	lastCheck time.Time
	mu        sync.Mutex
}

func (p *GitProvider) ToolName() string {
	return "git"
}

func (p *GitProvider) GetAliases(cwd string) []AliasEntry {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	if p.cached != nil && now.Sub(p.lastCheck) < 2*time.Second {
		return p.cached
	}

	key := p.buildCacheKey(cwd)
	p.lastCheck = now

	if p.cacheKey != "" && p.cacheKey == key {
		return p.cached
	}

	p.cached = p.parse(cwd)
	p.cacheKey = key
	return p.cached
}

func (p *GitProvider) buildCacheKey(cwd string) string {
	var sb strings.Builder
	sb.WriteString(cwd)

	// Check local config via fast os.Stat traversal
	if configPath := resolveGitConfigPath(cwd); configPath != "" {
		if info, err := os.Stat(configPath); err == nil {
			fmt.Fprintf(&sb, "|local:%s", info.ModTime().String())
		}
	}

	// Check global config
	if home, err := os.UserHomeDir(); err == nil {
		globalConfig := filepath.Join(home, ".gitconfig")
		if info, err := os.Stat(globalConfig); err == nil {
			fmt.Fprintf(&sb, "|global:%s", info.ModTime().String())
		}
	}

	return sb.String()
}

func resolveGitConfigPath(cwd string) string {
	dir := cwd
	for dir != "" {
		gitPath := filepath.Join(dir, ".git")
		info, err := os.Stat(gitPath)
		if err == nil {
			if info.IsDir() {
				return filepath.Join(gitPath, "config")
			}
			content, errRead := os.ReadFile(gitPath)
			if errRead == nil {
				s := strings.TrimSpace(string(content))
				if after, ok := strings.CutPrefix(s, "gitdir: "); ok {
					gitDir := strings.TrimSpace(after)
					if !filepath.IsAbs(gitDir) {
						gitDir = filepath.Join(dir, gitDir)
					}
					return filepath.Join(gitDir, "config")
				}
			}
			return filepath.Join(dir, ".git", "config")
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func (p *GitProvider) parse(cwd string) []AliasEntry {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Try with --show-scope first
	cmd := exec.CommandContext(ctx, "git", "config", "--get-regexp", "--show-scope", "^alias\\.")
	cmd.Dir = cwd
	out, err := cmd.Output()
	hasScope := true
	if err != nil {
		ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel2()
		// Fallback for older git
		cmd = exec.CommandContext(ctx2, "git", "config", "--get-regexp", "^alias\\.")
		cmd.Dir = cwd
		out, err = cmd.Output()
		hasScope = false
		if err != nil {
			return nil
		}
	}
	return p.parseOutput(out, hasScope)
}

func (p *GitProvider) parseOutput(out []byte, hasScope bool) []AliasEntry {
	var entries []AliasEntry
	lines := strings.Split(string(bytes.TrimSpace(out)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		scope := "local"
		if hasScope {
			idx := strings.IndexAny(line, " \t")
			if idx == -1 {
				continue
			}
			scope = line[:idx]
			rest := strings.TrimSpace(line[idx:])
			parts := strings.SplitN(rest, " ", 2)
			if len(parts) == 2 {
				name := strings.TrimPrefix(parts[0], "alias.")
				entries = append(entries, AliasEntry{
					Name:      name,
					Expansion: parts[1],
					Scope:     scope,
				})
			}
		} else {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				name := strings.TrimPrefix(parts[0], "alias.")
				entries = append(entries, AliasEntry{
					Name:      name,
					Expansion: parts[1],
					Scope:     scope,
				})
			}
		}
	}
	return entries
}

func init() {
	Register(&GitProvider{})
}
