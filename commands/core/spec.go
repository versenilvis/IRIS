package core

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type GeneratorFunc func(tokens []string, prefix string, partial string) []Suggestion

// Spec defines a top-level CLI command
type Spec struct {
	Name        string
	Aliases     []string
	Description string
	Subcommands []Subcommand
	Options     []Option
	Generator   GeneratorFunc
}

// Subcommand contains nested subcommands
type Subcommand struct {
	Name        string
	Aliases     []string
	Description string
	Subcommands []Subcommand
	Options     []Option
	Generator   GeneratorFunc
}

// Option represents a flag or option
type Option struct {
	Name        string
	Description string
}

// Suggestion is what gets shown in the dropdown
type Suggestion struct {
	Cmd  string
	Desc string
	Icon string
}

var (
	registry     = map[string]*Spec{}
	pathCmds     = make(map[string]bool)
	shellAliases = make(map[string]string) // alias name -> target command
	pathOnce     sync.Once
)

// scanShellAliases parses shell config files for aliases
func scanShellAliases() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	for _, f := range []string{".zshrc", ".bashrc", ".bash_profile", ".bash_aliases"} {
		content, err := os.ReadFile(filepath.Join(home, f))
		if err != nil {
			continue
		}

		for _, line := range strings.Split(string(content), "\n") {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "alias") {
				continue
			}

			body := strings.TrimSpace(strings.TrimPrefix(line, "alias"))
			for _, pair := range splitAliasTokens(body) {
				if eqIdx := strings.IndexByte(pair, '='); eqIdx > 0 {
					k := strings.TrimSpace(pair[:eqIdx])
					v := strings.Trim(strings.TrimSpace(pair[eqIdx+1:]), "\"'")
					if k != "" && v != "" {
						shellAliases[k] = v
					}
				}
			}
		}
	}
}

func splitAliasTokens(s string) []string {
	var pairs []string
	var cur strings.Builder
	inQuote := false
	var quote rune
	for _, c := range s {
		switch {
		case !inQuote && (c == '"' || c == '\''):
			inQuote, quote = true, c
			cur.WriteRune(c)
		case inQuote && c == quote:
			inQuote = false
			cur.WriteRune(c)
		case c == ' ' && !inQuote:
			if cur.Len() > 0 {
				pairs = append(pairs, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteRune(c)
		}
	}
	if cur.Len() > 0 {
		pairs = append(pairs, cur.String())
	}
	return pairs
}

// scanExternalCommands populates pathCmds and shellAliases
func scanExternalCommands() {
	scanPath()
	scanShellAliases()
}

// scanPath populates pathCmds with all executable files found in $PATH
func scanPath() {
	dirs := filepath.SplitList(os.Getenv("PATH"))
	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, f := range files {
			if !f.IsDir() {
				info, err := f.Info()
				if err == nil && info.Mode()&0111 != 0 {
					pathCmds[f.Name()] = true
				}
			}
		}
	}
}

// Register adds a spec to the global registry
func Register(s *Spec) {
	registry[s.Name] = s
}

// Lookup finds suggestions based on user input
// We need to tokenize the input to detect the depth of the command,
// like when you type "gi" and it shows suggestion about "git",
// then if you press tab (not enter) it will continuely show suggestions about "git" subcommands,
// when you type like "git add", it will continuely show the next suggestion from "git add"
// if we don't tokenize it. iris cannot get which command to suggest.
// And tokenizer helps us to detect white space like we can't know if you are typing "git-lfs"
// or "git" to show the suggestions, without it, if you type "git " (with white space), it still shows
// suggestion for "git-lfs"
//
// priority: file/dir -> subcommands -> options
func Lookup(input string) []Suggestion {
	pathOnce.Do(scanExternalCommands)
	tokens := tokenize(input)

	// Token Injection: resolve shell alias (e.g. gca -> git commit -a)
	if len(tokens) > 1 {
		if target, ok := shellAliases[tokens[0]]; ok {
			aliasTokens := tokenize(target)
			if len(aliasTokens) > 0 && aliasTokens[len(aliasTokens)-1] == "" {
				aliasTokens = aliasTokens[:len(aliasTokens)-1]
			}
			tokens = append(aliasTokens, tokens[1:]...)
		}
	}

	if len(tokens) == 0 {
		return topLevelSuggestions(input)
	}

	if len(tokens) == 1 {
		return topLevelSuggestions(tokens[0])
	}

	rootCmdName := tokens[0]
	spec, exists := registry[rootCmdName]
	if !exists {
		return nil
	}

	currentSubs, currentOpts, currentGen := spec.Subcommands, spec.Options, spec.Generator
	depth := 1

	for depth < len(tokens) {
		tok := tokens[depth]
		if tok == "" || strings.HasPrefix(tok, "-") || strings.Contains(tok, "=") {
			depth++
			continue
		}

		found := false
		for _, sub := range currentSubs {
			match := sub.Name == tok
			if !match {
				for _, a := range sub.Aliases {
					if a == tok {
						match = true
						break
					}
				}
			}

			if match {
				currentSubs, currentOpts, currentGen = sub.Subcommands, sub.Options, sub.Generator
				found = true
				break
			}
		}
		if found {
			depth++
			continue
		}

		if currentGen != nil && depth < len(tokens)-1 {
			depth++
			continue
		}
		break
	}

	results := []Suggestion{}
	prefixBuilder := strings.Builder{}
	for i := 0; i < depth; i++ {
		if i > 0 {
			prefixBuilder.WriteByte(' ')
		}
		prefixBuilder.WriteString(tokens[i])
	}
	prefix := prefixBuilder.String()
	partial := tokens[len(tokens)-1]

	if currentGen != nil {
		genResults := currentGen(tokens[:depth], prefix, partial)
		for _, g := range genResults {
			parts := strings.Split(g.Cmd, " ")
			name := parts[len(parts)-1]
			if partial == "" || hasPrefix(name, partial) {
				results = append(results, Suggestion{
					Cmd: g.Cmd, Desc: g.Desc, Icon: rootCmdName,
				})
			}
		}
	}

	for _, sub := range currentSubs {
		if partial == "" || hasPrefix(sub.Name, partial) {
			results = append(results, Suggestion{
				Cmd: prefix + " " + sub.Name, Desc: sub.Description, Icon: rootCmdName,
			})
		}
	}

	if partial == "" || (len(partial) > 0 && partial[0] == '-') {
		usedOpts := make(map[string]bool)
		for _, t := range tokens {
			if strings.HasPrefix(t, "-") {
				usedOpts[t] = true
			}
		}
		for _, opt := range currentOpts {
			if !usedOpts[opt.Name] && (partial == "" || hasPrefix(opt.Name, partial)) {
				results = append(results, Suggestion{
					Cmd: prefix + " " + opt.Name, Desc: opt.Description, Icon: rootCmdName,
				})
			}
		}
	}

	return results
}

func topLevelSuggestions(query string) []Suggestion {
	pathOnce.Do(scanExternalCommands)
	results, seen := []Suggestion{}, make(map[string]bool)

	// 1. Manual specs (High Priority)
	for name, spec := range registry {
		match := false
		if query == "" || hasPrefix(name, query) {
			match = true
		} else {
			for _, a := range spec.Aliases {
				if hasPrefix(a, query) {
					match = true
					break
				}
			}
		}
		if match {
			results = append(results, Suggestion{Cmd: name, Desc: spec.Description, Icon: name})
			seen[name] = true
		}
	}

	// 2. shell aliases (User's choice precedence)
	for name, target := range shellAliases {
		if !seen[name] && (query == "" || hasPrefix(name, query)) {
			results = append(results, Suggestion{
				Cmd: target, Desc: "alias: " + name, Icon: "root",
			})
			seen[name] = true
		}
	}

	// 3. system commands from $PATH
	for name := range pathCmds {
		if !seen[name] && (query == "" || hasPrefix(name, query)) {
			results = append(results, Suggestion{
				Cmd: name, Desc: "system command", Icon: "root",
			})
		}
	}

	return results
}

func tokenize(s string) []string {
	tokens := []string{}
	var current strings.Builder
	inQuote := false
	var quoteChar rune

	for _, c := range s {
		switch {
		case !inQuote && (c == '"' || c == '\''):
			inQuote = true
			quoteChar = c
		case inQuote && c == quoteChar:
			inQuote = false
		case c == ' ' && !inQuote:
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			} else if len(tokens) > 0 && tokens[len(tokens)-1] != "" {
				// preserve trailing space
			}
		default:
			current.WriteRune(c)
		}
	}

	tokens = append(tokens, current.String())
	return tokens
}

func hasPrefix(s, prefix string) bool {
	if len(prefix) > len(s) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		a, b := s[i], prefix[i]
		if a >= 'A' && a <= 'Z' {
			a += 32
		}
		if b >= 'A' && b <= 'Z' {
			b += 32
		}
		if a != b {
			return false
		}
	}
	return true
}
