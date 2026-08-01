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
	cacheKey string
	cached   []AliasEntry
	mu       sync.Mutex
}

func (p *GitProvider) ToolName() string {
	return "git"
}

func (p *GitProvider) GetAliases(cwd string) []AliasEntry {
	key := p.buildCacheKey(cwd)
	p.mu.Lock()
	defer p.mu.Unlock()

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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Check local config
	cmd := exec.CommandContext(ctx, "git", "-C", cwd, "rev-parse", "--show-toplevel")
	if gitDir, err := cmd.Output(); err == nil {
		localConfig := filepath.Join(strings.TrimSpace(string(gitDir)), ".git", "config")
		if info, err := os.Stat(localConfig); err == nil {
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
