package ai

import (
	"context"
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

type DockerExecProvider struct{}

func (p DockerExecProvider) Name() string {
	return "docker-exec"
}

func (p DockerExecProvider) Matches(buf string) bool {
	return strings.HasPrefix(buf, "docker exec")
}

func (p DockerExecProvider) Gather(ctx context.Context) (string, error) {
	out, err := exec.CommandContext(ctx, "docker", "ps", "--format", "{{.Names}}\t{{.Image}}").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

type GitCommitProvider struct{}

func (p GitCommitProvider) Name() string {
	return "git-commit"
}

func (p GitCommitProvider) Matches(buf string) bool {
	return strings.HasPrefix(buf, "git commit")
}

func (p GitCommitProvider) Gather(ctx context.Context) (string, error) {
	out, err := exec.CommandContext(ctx, "git", "status", "-s").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
