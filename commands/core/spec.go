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
	Description string
	Subcommands []Subcommand
	Options     []Option
	Generator   GeneratorFunc
}

// Subcommand contains nested subcommands
type Subcommand struct {
	Name        string
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
	registry = map[string]*Spec{}
	pathCmds = make(map[string]bool)
	pathOnce sync.Once
)

// scanPath populates pathCmds with all executable files found in $PATH
func scanPath() {
	pathVar := os.Getenv("PATH")
	dirs := filepath.SplitList(pathVar)

	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, f := range files {
			if f.IsDir() {
				continue
			}
			// check if it's executable
			info, err := f.Info()
			if err == nil && info.Mode()&0111 != 0 {
				pathCmds[f.Name()] = true
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
	tokens := tokenize(input)
	if len(tokens) == 0 {
		return topLevelSuggestions(input)
	}

	// if only 1 token and no trailing space, user is still typing the root command
	if len(tokens) == 1 {
		return topLevelSuggestions(tokens[0])
	}

	root := tokens[0]
	spec, exists := registry[root]
	if !exists {
		// no spec found for this command, return nothing to avoid clutter
		return nil
	}

	currentSubs := spec.Subcommands
	currentOpts := spec.Options
	currentGen := spec.Generator
	depth := 1

	for depth < len(tokens) {
		tok := tokens[depth]
		if tok == "" {
			break
		}

		// skip options
		if strings.HasPrefix(tok, "-") {
			depth++
			continue
		}

		// try to match subcommands
		found := false
		for _, sub := range currentSubs {
			if sub.Name == tok {
				currentSubs = sub.Subcommands
				currentOpts = sub.Options
				currentGen = sub.Generator
				found = true
				break
			}
		}
		if found {
			depth++
			continue
		}

		// if no subcommand matches but we have a generator,
		// and this is NOT the last token (meaning it's a finished argument),
		// we consume it and move depth forward.
		if currentGen != nil && depth < len(tokens)-1 {
			depth++
			continue
		}

		break
	}

	// build prefix from tokens consumed so far
	prefix := ""
	for i := 0; i < depth; i++ {
		if i > 0 {
			prefix += " "
		}
		prefix += tokens[i]
	}

	// partial is what user is currently typing (might be incomplete)
	partial := ""
	if depth < len(tokens) {
		partial = tokens[depth]
	}

	results := []Suggestion{}

	// file/dir
	if currentGen != nil {
		genResults := currentGen(tokens[:depth], prefix, partial)
		for _, g := range genResults {
			// extract simple name for prefix matching
			parts := strings.Split(g.Cmd, " ")
			name := parts[len(parts)-1]
			if partial == "" || hasPrefix(name, partial) {
				results = append(results, Suggestion{
					Cmd:  g.Cmd,
					Desc: g.Desc,
					Icon: root,
				})
			}
		}
	}

	// subcommands
	for _, sub := range currentSubs {
		if partial == "" || hasPrefix(sub.Name, partial) {
			results = append(results, Suggestion{
				Cmd:  prefix + " " + sub.Name,
				Desc: sub.Description,
				Icon: root,
			})
		}
	}

	// 3. Options (Flags)
	if partial == "" || (len(partial) > 0 && partial[0] == '-') {
		// Identify already used options to filter them out
		usedOpts := make(map[string]bool)
		for _, t := range tokens {
			if strings.HasPrefix(t, "-") {
				usedOpts[t] = true
			}
		}

		for _, opt := range currentOpts {
			// Skip if already used
			if usedOpts[opt.Name] {
				continue
			}

			if partial == "" || hasPrefix(opt.Name, partial) {
				results = append(results, Suggestion{
					Cmd:  prefix + " " + opt.Name,
					Desc: opt.Description,
					Icon: root,
				})
			}
		}
	}

	return results
}

func topLevelSuggestions(query string) []Suggestion {
	pathOnce.Do(scanPath)

	results := []Suggestion{}
	seen := make(map[string]bool)

	for name, spec := range registry {
		if query == "" || hasPrefix(name, query) {
			results = append(results, Suggestion{
				Cmd:  name,
				Desc: spec.Description,
				Icon: name,
			})
			seen[name] = true
		}
	}

	for name := range pathCmds {
		if seen[name] {
			continue
		}
		if query == "" || hasPrefix(name, query) {
			results = append(results, Suggestion{
				Cmd:  name,
				Desc: "system command",
				Icon: "root",
			})
		}
	}

	return results
}

func tokenize(s string) []string {
	tokens := []string{}
	current := ""
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
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			} else if len(tokens) > 0 && tokens[len(tokens)-1] != "" {
				// preserve trailing space as an empty token for suggestions
				// but only if it's the very last thing and not multiple spaces
			}
		default:
			current += string(c)
		}
	}

	// always append the last token (even if empty, representing the cursor position)
	tokens = append(tokens, current)
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
