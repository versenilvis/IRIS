package ai

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type ContextProvider interface {
	Name() string
	Matches(buf string) bool
	Gather(ctx context.Context) (string, error)
}

type cacheEntry struct {
	data       string
	expireTime time.Time
}

type ProviderCache struct {
	mu      sync.Mutex
	entries map[string]cacheEntry
	ttl     time.Duration
}

func NewProviderCache(ttl time.Duration) *ProviderCache {
	if ttl == 0 {
		ttl = 4 * time.Second
	}
	return &ProviderCache{
		entries: make(map[string]cacheEntry),
		ttl:     ttl,
	}
}

func (c *ProviderCache) GetOrGather(ctx context.Context, p ContextProvider) string {
	c.mu.Lock()
	entry, ok := c.entries[p.Name()]
	if ok && time.Now().Before(entry.expireTime) {
		c.mu.Unlock()
		return entry.data
	}
	c.mu.Unlock()

	data, err := p.Gather(ctx)
	if err != nil || ctx.Err() != nil {
		return ""
	}

	c.mu.Lock()
	c.entries[p.Name()] = cacheEntry{
		data:       data,
		expireTime: time.Now().Add(c.ttl),
	}
	c.mu.Unlock()
	return data
}

func (c *ProviderCache) Clear() {
	c.mu.Lock()
	c.entries = make(map[string]cacheEntry)
	c.mu.Unlock()
}

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
		return p.Label + ":\n" + s, nil
	}
	return "", nil
}

var DefaultProviders = []*CommandContextProvider{
	{NameStr: "docker_exec", Prefixes: []string{"docker exec", "docker logs", "docker stop", "docker restart", "docker rm"},
		GatherCmd: []string{"docker", "ps", "--format", "{{.Names}}\t{{.Image}}"}, Label: "Running containers"},
	{NameStr: "docker_compose", Prefixes: []string{"docker compose exec", "docker compose logs", "docker-compose exec", "docker-compose logs"},
		GatherCmd: []string{"docker", "compose", "ps", "--format", "{{.Name}}\t{{.Service}}"}, Label: "Compose services"},
	{NameStr: "kubectl_pods", Prefixes: []string{"kubectl exec", "kubectl logs", "kubectl describe pod", "kubectl delete pod"},
		GatherCmd: []string{"kubectl", "get", "pods", "--no-headers"}, Label: "Pods"},
	{NameStr: "git_branch", Prefixes: []string{"git checkout", "git switch", "git merge", "git rebase", "git branch -d", "git branch -D"},
		GatherCmd: []string{"git", "branch", "-a", "--format=%(refname:short)"}, Label: "Branches"},
	{NameStr: "kill_proc", Prefixes: []string{"kill ", "kill -9 "},
		GatherCmd: []string{"ps", "-eo", "pid,comm,%cpu,%mem", "--sort=-%cpu"}, Label: "Top processes"},
	{NameStr: "systemctl", Prefixes: []string{"systemctl restart", "systemctl stop", "systemctl status"},
		GatherCmd: []string{"systemctl", "list-units", "--type=service", "--no-legend"}, Label: "Services"},
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
			}
			names = append(names, name)
		}
		if len(names) > 0 {
			sb.WriteString(fmt.Sprintf("Files in Cwd: %s\n\n", strings.Join(names, ", ")))
		}
	}

	cmd := exec.CommandContext(ctxTimeout, "git", "-C", p.cwd, "rev-parse", "--is-inside-work-tree")
	if cmd.Run() == nil {
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
		if statusStr != "" {
			sb.WriteString(fmt.Sprintf("Status:\n%s\n\n", statusStr))
		}
		if diffStr != "" {
			sb.WriteString(fmt.Sprintf("Staged Diff:\n%s\n\n", diffStr))
		}
		if logStr != "" {
			sb.WriteString(fmt.Sprintf("User's recent commit messages (MUST follow this exact style, formatting, language, and casing conventions):\n%s\n", logStr))
		}
	}

	if fields := strings.Fields(p.buf); len(fields) > 0 {
		ctxHelp, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
		defer cancel()
		helpOut, err := exec.CommandContext(ctxHelp, fields[0], "--help").CombinedOutput()
		if err == nil {
			helpStr := strings.TrimSpace(string(helpOut))
			if len(helpStr) > 1200 {
				helpStr = helpStr[:1200]
			}
			if helpStr != "" {
				sb.WriteString(fmt.Sprintf("\nCommand help (%s --help):\n%s\n", fields[0], helpStr))
			}
		}
	}

	return strings.TrimSpace(sb.String()), nil
}
