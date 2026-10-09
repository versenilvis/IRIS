package spec

import (
	"bytes"
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/versenilvis/iris/internal/config"
	"github.com/versenilvis/iris/internal/logger"
)

var (
	userSpecsMu     sync.Mutex
	loadedUserSpecs = map[string]bool{}
	builtinSnapshot map[string]*Spec
	builtinOnce     sync.Once
)

type UserSpecJSON struct {
	Name        string               `json:"name"`
	Aliases     []string             `json:"aliases,omitempty"`
	Description string               `json:"description,omitempty"`
	Icon        string               `json:"icon,omitempty"`
	Generator   any                  `json:"generator,omitempty"`
	MaxArgs     int                  `json:"max_args,omitempty"`
	Options     []UserOptionJSON     `json:"options,omitempty"`
	Subcommands []UserSubcommandJSON `json:"subcommands,omitempty"`
}

type UserSubcommandJSON struct {
	Name        string               `json:"name"`
	Aliases     []string             `json:"aliases,omitempty"`
	Description string               `json:"description,omitempty"`
	Icon        string               `json:"icon,omitempty"`
	Generator   any                  `json:"generator,omitempty"`
	MaxArgs     int                  `json:"max_args,omitempty"`
	Priority    int                  `json:"priority,omitempty"`
	Options     []UserOptionJSON     `json:"options,omitempty"`
	Subcommands []UserSubcommandJSON `json:"subcommands,omitempty"`
}

type UserOptionJSON struct {
	Name        string   `json:"name"`
	Aliases     []string `json:"aliases,omitempty"`
	Description string   `json:"description,omitempty"`
	Priority    int      `json:"priority,omitempty"`
}

func CleanJSONC(data []byte) []byte {
	var noComments []byte
	inString := false
	escaped := false
	n := len(data)

	for i := 0; i < n; i++ {
		b := data[i]

		if inString {
			noComments = append(noComments, b)
			if escaped {
				escaped = false
			} else if b == '\\' {
				escaped = true
			} else if b == '"' {
				inString = false
			}
			continue
		}

		if b == '/' && i+1 < n {
			if data[i+1] == '/' {
				i += 2
				for i < n && data[i] != '\n' {
					i++
				}
				if i < n {
					noComments = append(noComments, '\n')
				}
				continue
			} else if data[i+1] == '*' {
				i += 2
				for i+1 < n && (data[i] != '*' || data[i+1] != '/') {
					if data[i] == '\n' {
						noComments = append(noComments, '\n')
					}
					i++
				}
				i++
				continue
			}
		}

		if b == '"' {
			inString = true
			noComments = append(noComments, b)
			continue
		}

		noComments = append(noComments, b)
	}

	var out []byte
	inString = false
	escaped = false
	m := len(noComments)

	for i := range m {
		b := noComments[i]

		if inString {
			out = append(out, b)
			if escaped {
				escaped = false
			} else if b == '\\' {
				escaped = true
			} else if b == '"' {
				inString = false
			}
			continue
		}

		if b == '"' {
			inString = true
			out = append(out, b)
			continue
		}

		if b == ',' {
			j := i + 1
			for j < m && (noComments[j] == ' ' || noComments[j] == '\t' || noComments[j] == '\r' || noComments[j] == '\n') {
				j++
			}
			if j < m && (noComments[j] == '}' || noComments[j] == ']') {
				continue
			}
		}

		out = append(out, b)
	}

	return out
}

func resolveGenerator(v any) GeneratorFunc {
	if v == nil {
		return nil
	}

	switch val := v.(type) {
	case string:
		lower := strings.ToLower(strings.TrimSpace(val))
		switch lower {
		case "file", "files":
			return FileGenerator()
		case "dir", "dirs", "directory":
			return FileGenerator("/")
		default:
			return nil
		}
	case []any:
		var filters []string
		for _, item := range val {
			if s, ok := item.(string); ok && s != "" {
				filters = append(filters, s)
			}
		}
		if len(filters) > 0 {
			return FileGenerator(filters...)
		}
	case []string:
		if len(val) > 0 {
			return FileGenerator(val...)
		}
	case map[string]any:
		genType, _ := val["type"].(string)
		lowerType := strings.ToLower(strings.TrimSpace(genType))
		var filters []string

		if rawFilter, ok := val["filter"]; ok {
			switch rf := rawFilter.(type) {
			case string:
				if rf != "" {
					filters = append(filters, rf)
				}
			case []any:
				for _, item := range rf {
					if s, ok := item.(string); ok && s != "" {
						filters = append(filters, s)
					}
				}
			}
		}

		switch lowerType {
		case "dir", "dirs", "directory":
			return FileGenerator("/")
		case "file", "files":
			if len(filters) > 0 {
				return FileGenerator(filters...)
			}
			return FileGenerator()
		}
	}

	return nil
}

func convertUserSubcommand(raw *UserSubcommandJSON) Subcommand {
	sub := Subcommand{
		Name:        raw.Name,
		Aliases:     raw.Aliases,
		Description: raw.Description,
		Icon:        raw.Icon,
		MaxArgs:     raw.MaxArgs,
		Priority:    raw.Priority,
		Generator:   resolveGenerator(raw.Generator),
	}
	for _, o := range raw.Options {
		sub.Options = append(sub.Options, Option{
			Name:        o.Name,
			Description: o.Description,
			Priority:    o.Priority,
		})
		for _, a := range o.Aliases {
			if a != "" {
				sub.Options = append(sub.Options, Option{
					Name:        a,
					Description: o.Description,
					Priority:    o.Priority,
				})
			}
		}
	}
	for _, s := range raw.Subcommands {
		sub.Subcommands = append(sub.Subcommands, convertUserSubcommand(&s))
	}
	return sub
}

func convertUserSpec(raw *UserSpecJSON) *Spec {
	s := &Spec{
		Name:        raw.Name,
		Aliases:     raw.Aliases,
		Description: raw.Description,
		Icon:        raw.Icon,
		MaxArgs:     raw.MaxArgs,
		Generator:   resolveGenerator(raw.Generator),
	}
	for _, o := range raw.Options {
		s.Options = append(s.Options, Option{
			Name:        o.Name,
			Description: o.Description,
			Priority:    o.Priority,
		})
		for _, a := range o.Aliases {
			if a != "" {
				s.Options = append(s.Options, Option{
					Name:        a,
					Description: o.Description,
					Priority:    o.Priority,
				})
			}
		}
	}
	for _, sub := range raw.Subcommands {
		s.Subcommands = append(s.Subcommands, convertUserSubcommand(&sub))
	}
	return s
}

func LoadUserSpecBytes(data []byte, defaultName string) ([]*Spec, error) {
	cleaned := CleanJSONC(data)
	trimmed := bytes.TrimSpace(cleaned)
	if len(trimmed) == 0 {
		return nil, nil
	}

	var rawList []UserSpecJSON
	if trimmed[0] == '[' {
		if err := json.Unmarshal(trimmed, &rawList); err != nil {
			return nil, err
		}
	} else {
		var wrapper struct {
			Specs []UserSpecJSON `json:"specs"`
		}
		if err := json.Unmarshal(trimmed, &wrapper); err == nil && len(wrapper.Specs) > 0 {
			rawList = wrapper.Specs
		} else {
			var single UserSpecJSON
			if err := json.Unmarshal(trimmed, &single); err != nil {
				return nil, err
			}
			if single.Name == "" && defaultName != "" {
				single.Name = defaultName
			}
			rawList = []UserSpecJSON{single}
		}
	}

	var specs []*Spec
	for i := range rawList {
		raw := &rawList[i]
		if raw.Name == "" && defaultName != "" && len(rawList) == 1 {
			raw.Name = defaultName
		}
		if raw.Name != "" {
			specs = append(specs, convertUserSpec(raw))
		}
	}

	return specs, nil
}

func GetUserSpecsDir() string {
	cfgDir, err := config.ConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(cfgDir, "specs")
}

func LoadUserSpecs(dir string) error {
	if dir == "" {
		return nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	userSpecsMu.Lock()
	defer userSpecsMu.Unlock()

	builtinOnce.Do(func() {
		builtinSnapshot = make(map[string]*Spec, len(Registry))
		maps.Copy(builtinSnapshot, Registry)
	})

	for name := range loadedUserSpecs {
		if orig, ok := builtinSnapshot[name]; ok {
			Registry[name] = orig
		} else {
			delete(Registry, name)
		}
	}
	clear(loadedUserSpecs)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".json" && ext != ".jsonc" {
			continue
		}

		filePath := filepath.Join(dir, name)
		data, readErr := os.ReadFile(filePath)
		if readErr != nil {
			logger.Errorf("failed to read user spec file %s: %v", filePath, readErr)
			continue
		}

		defaultName := strings.TrimSuffix(name, filepath.Ext(name))
		specs, parseErr := LoadUserSpecBytes(data, defaultName)
		if parseErr != nil {
			logger.Errorf("failed to parse user spec file %s: %v", filePath, parseErr)
			continue
		}

		for _, s := range specs {
			Registry[s.Name] = s
			loadedUserSpecs[s.Name] = true
			for _, alias := range s.Aliases {
				if _, exists := Registry[alias]; !exists {
					Registry[alias] = s
					loadedUserSpecs[alias] = true
				}
			}
			logger.Debugf("registered user spec: %s from %s", s.Name, filePath)
		}
	}

	return nil
}

func dirSignature(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	var b strings.Builder
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".json" && ext != ".jsonc" {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		b.WriteString(entry.Name())
		b.WriteString(strconv.FormatInt(info.ModTime().UnixNano(), 10))
		b.WriteString(strconv.FormatInt(info.Size(), 10))
	}
	return b.String()
}

func AutoDetectSpecsChange(dir string, onReload func()) {
	if dir == "" {
		return
	}

	go func() {
		lastSig := dirSignature(dir)
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			currentSig := dirSignature(dir)
			if currentSig != lastSig {
				lastSig = currentSig
				if err := LoadUserSpecs(dir); err == nil {
					if onReload != nil {
						onReload()
					}
				}
			}
		}
	}()
}
