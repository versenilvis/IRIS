package synthesize

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/versenilvis/iris/spec"
)

const defaultSchemaURL = "https://raw.githubusercontent.com/versenilvis/iris/main/spec/schema.json"

type Synthesizer struct {
	bashDir     string
	manDirs     []string
	maxFileSize int64
}

type Option func(*Synthesizer)

func WithBashCompletionDir(dir string) Option {
	return func(s *Synthesizer) {
		s.bashDir = dir
	}
}

func WithManDir(dir string) Option {
	return func(s *Synthesizer) {
		s.manDirs = []string{dir}
	}
}

func WithManDirs(dirs []string) Option {
	return func(s *Synthesizer) {
		s.manDirs = dirs
	}
}

func WithMaxFileSize(size int64) Option {
	return func(s *Synthesizer) {
		s.maxFileSize = size
	}
}

func New(opts ...Option) *Synthesizer {
	s := &Synthesizer{
		bashDir:     "/usr/share/bash-completion/completions",
		manDirs:     []string{"/usr/share/man/man1", "/usr/share/man/man8"},
		maxFileSize: defaultMaxFileSize,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func Synthesize(cmdName string, opts ...Option) (*spec.UserSpecJSON, error) {
	return New(opts...).Synthesize(cmdName)
}

func (s *Synthesizer) Synthesize(cmdName string) (*spec.UserSpecJSON, error) {
	trimmedCmd := strings.TrimSpace(cmdName)
	if trimmedCmd == "" {
		return nil, errors.New("command name cannot be empty")
	}
	if filepath.Base(trimmedCmd) != trimmedCmd || strings.ContainsAny(trimmedCmd, "/\\ ") || strings.Contains(trimmedCmd, "..") {
		return nil, fmt.Errorf("invalid command name %q", cmdName)
	}

	var cmdDesc string
	var manOptions []spec.UserOptionJSON
	for _, manDir := range s.manDirs {
		if !filepath.IsAbs(manDir) {
			continue
		}
		desc, opts, err := parseManFile(manDir, trimmedCmd, s.maxFileSize)
		if err == nil && (len(opts) > 0 || desc != "") {
			cmdDesc = desc
			manOptions = opts
			break
		}
	}

	var bashOptions []spec.UserOptionJSON
	if filepath.IsAbs(s.bashDir) {
		opts, err := parseBashFile(s.bashDir, trimmedCmd, s.maxFileSize)
		if err == nil && len(opts) > 0 {
			bashOptions = opts
		}
	}

	if len(manOptions) == 0 && len(bashOptions) == 0 && cmdDesc == "" {
		return nil, fmt.Errorf("no static completion sources found for %q", trimmedCmd)
	}

	merged := mergeOptions(manOptions, bashOptions)

	cmdDesc = sanitizeText(cmdDesc)
	cmdDesc = truncateRunes(cmdDesc, 160)
	cmdDesc = strings.TrimRight(cmdDesc, ". \t\r\n")
	if cmdDesc == "" {
		cmdDesc = trimmedCmd + " command completion"
	}

	return &spec.UserSpecJSON{
		Schema:      defaultSchemaURL,
		Name:        trimmedCmd,
		Description: cmdDesc,
		Options:     merged,
	}, nil
}

func mergeOptions(manOpts, bashOpts []spec.UserOptionJSON) []spec.UserOptionJSON {
	var merged []spec.UserOptionJSON
	seenFlags := make(map[string]int)

	for _, opt := range manOpts {
		idx := len(merged)
		merged = append(merged, opt)
		seenFlags[opt.Name] = idx
		for _, a := range opt.Aliases {
			seenFlags[a] = idx
		}
	}

	for _, opt := range bashOpts {
		if idx, exists := seenFlags[opt.Name]; exists {
			for _, a := range opt.Aliases {
				if _, aliasExists := seenFlags[a]; !aliasExists {
					merged[idx].Aliases = append(merged[idx].Aliases, a)
					seenFlags[a] = idx
				}
			}
			continue
		}

		matchedIdx := -1
		for _, a := range opt.Aliases {
			if idx, exists := seenFlags[a]; exists {
				matchedIdx = idx
				break
			}
		}

		if matchedIdx >= 0 {
			if _, nameExists := seenFlags[opt.Name]; !nameExists {
				merged[matchedIdx].Aliases = append(merged[matchedIdx].Aliases, opt.Name)
				seenFlags[opt.Name] = matchedIdx
			}
			for _, a := range opt.Aliases {
				if _, aliasExists := seenFlags[a]; !aliasExists {
					merged[matchedIdx].Aliases = append(merged[matchedIdx].Aliases, a)
					seenFlags[a] = matchedIdx
				}
			}
		} else {
			idx := len(merged)
			merged = append(merged, opt)
			seenFlags[opt.Name] = idx
			for _, a := range opt.Aliases {
				seenFlags[a] = idx
			}
		}
	}

	var sanitized []spec.UserOptionJSON
	for _, opt := range merged {
		name := sanitizeText(opt.Name)
		if !reValidFlag.MatchString(name) {
			continue
		}
		var cleanAliases []string
		for _, a := range opt.Aliases {
			ca := sanitizeText(a)
			if reValidFlag.MatchString(ca) && ca != name {
				cleanAliases = append(cleanAliases, ca)
			}
		}
		desc := sanitizeText(opt.Description)
		desc = truncateRunes(desc, 160)
		desc = strings.TrimRight(desc, ". \t\r\n,;:")
		sanitized = append(sanitized, spec.UserOptionJSON{
			Name:        name,
			Aliases:     cleanAliases,
			Description: desc,
		})
	}
	merged = sanitized

	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Name < merged[j].Name
	})

	for i := range merged {
		sort.Strings(merged[i].Aliases)
	}

	return merged
}
