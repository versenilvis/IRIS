package alias

import (
	"sync"
)

type AliasEntry struct {
	Name      string
	Expansion string
	Scope     string // "local" | "global"
}

type Provider interface {
	ToolName() string
	GetAliases(cwd string) []AliasEntry
}

var (
	providers = make(map[string]Provider)
	mu        sync.RWMutex
)

func Register(p Provider) {
	mu.Lock()
	defer mu.Unlock()
	providers[p.ToolName()] = p
}

func GetProvider(toolName string) Provider {
	mu.RLock()
	defer mu.RUnlock()
	return providers[toolName]
}
