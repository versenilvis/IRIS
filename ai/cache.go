package ai

import (
	"context"
	"sync"
	"time"

	"github.com/versenilvis/iris/spec"
)

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

type ContextCache struct {
	mu               sync.Mutex
	lastSnapshotHash string
	lastSuggestion   *spec.Suggestion
	lastFetchedAt    time.Time
}

func NewContextCache() *ContextCache {
	return &ContextCache{}
}

func (c *ContextCache) ShouldCallAI(snap EnvSnapshot, minInterval time.Duration) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	hash := snap.Hash()
	if hash == c.lastSnapshotHash {
		return false
	}
	if time.Since(c.lastFetchedAt) < minInterval {
		return false
	}
	return true
}

func (c *ContextCache) GetCachedSuggestion(snap EnvSnapshot) *spec.Suggestion {
	c.mu.Lock()
	defer c.mu.Unlock()
	if snap.Hash() == c.lastSnapshotHash {
		return c.lastSuggestion
	}
	return nil
}

func (c *ContextCache) Update(snap EnvSnapshot, sugg *spec.Suggestion) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastSnapshotHash = snap.Hash()
	c.lastSuggestion = sugg
	c.lastFetchedAt = time.Now()
}

func (c *ContextCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastSnapshotHash = ""
	c.lastSuggestion = nil
}
