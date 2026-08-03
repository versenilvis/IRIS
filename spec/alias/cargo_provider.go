package alias

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
)

type CargoProvider struct {
	cacheKey string
	cached   []AliasEntry
	mu       sync.Mutex
}

func (p *CargoProvider) ToolName() string {
	return "cargo"
}

func (p *CargoProvider) GetAliases(cwd string) []AliasEntry {
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

func (p *CargoProvider) buildCacheKey(cwd string) string {
	var sb strings.Builder
	sb.WriteString(cwd)

	dir := cwd
	for {
		localConfig := filepath.Join(dir, ".cargo", "config.toml")
		if info, err := os.Stat(localConfig); err == nil {
			fmt.Fprintf(&sb, "|local:%s", info.ModTime().String())
		}
		localConfig2 := filepath.Join(dir, ".cargo", "config")
		if info, err := os.Stat(localConfig2); err == nil {
			fmt.Fprintf(&sb, "|local:%s", info.ModTime().String())
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	if cargoHome := os.Getenv("CARGO_HOME"); cargoHome != "" {
		globalConfig := filepath.Join(cargoHome, "config.toml")
		if info, err := os.Stat(globalConfig); err == nil {
			fmt.Fprintf(&sb, "|global:%s", info.ModTime().String())
		}
	} else if home, err := os.UserHomeDir(); err == nil {
		globalConfig := filepath.Join(home, ".cargo", "config.toml")
		if info, err := os.Stat(globalConfig); err == nil {
			fmt.Fprintf(&sb, "|global:%s", info.ModTime().String())
		}
	}

	return sb.String()
}

func (p *CargoProvider) parse(cwd string) []AliasEntry {
	var rawEntries []AliasEntry

	dir := cwd
	for {
		rawEntries = append(rawEntries, p.parseFile(filepath.Join(dir, ".cargo", "config.toml"), "local")...)
		rawEntries = append(rawEntries, p.parseFile(filepath.Join(dir, ".cargo", "config"), "local")...)

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	if cargoHome := os.Getenv("CARGO_HOME"); cargoHome != "" {
		rawEntries = append(rawEntries, p.parseFile(filepath.Join(cargoHome, "config.toml"), "global")...)
		rawEntries = append(rawEntries, p.parseFile(filepath.Join(cargoHome, "config"), "global")...)
	} else if home, err := os.UserHomeDir(); err == nil {
		rawEntries = append(rawEntries, p.parseFile(filepath.Join(home, ".cargo", "config.toml"), "global")...)
		rawEntries = append(rawEntries, p.parseFile(filepath.Join(home, ".cargo", "config"), "global")...)
	}

	seen := make(map[string]bool)
	var entries []AliasEntry
	for _, e := range rawEntries {
		if !seen[e.Name] {
			seen[e.Name] = true
			entries = append(entries, e)
		}
	}

	return entries
}

func (p *CargoProvider) parseFile(path, scope string) []AliasEntry {
	var config struct {
		Alias map[string]interface{} `toml:"alias"`
	}
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil
	}

	var entries []AliasEntry
	for k, v := range config.Alias {
		switch val := v.(type) {
		case string:
			entries = append(entries, AliasEntry{Name: k, Expansion: val, Scope: scope})
		case []interface{}:
			var parts []string
			for _, item := range val {
				if s, ok := item.(string); ok {
					parts = append(parts, s)
				}
			}
			entries = append(entries, AliasEntry{Name: k, Expansion: strings.Join(parts, " "), Scope: scope})
		}
	}
	return entries
}

func init() {
	Register(&CargoProvider{})
}
