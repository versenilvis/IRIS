package core

import "strings"

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

var registry = map[string]*Spec{}

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
		if !found {
			break
		}
		depth++
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

	// options
	if partial == "" || (len(partial) > 0 && partial[0] == '-') {
		for _, opt := range currentOpts {
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
	results := []Suggestion{}
	for name, spec := range registry {
		if query == "" || hasPrefix(name, query) {
			results = append(results, Suggestion{
				Cmd:  name,
				Desc: spec.Description,
				Icon: name,
			})
		}
	}
	return results
}

func tokenize(s string) []string {
	tokens := []string{}
	current := ""
	for _, c := range s {
		if c == ' ' {
			tokens = append(tokens, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	tokens = append(tokens, current) // always append the last token (could be empty if space was last character)
	return tokens
}

func hasPrefix(s, prefix string) bool {
	if len(prefix) > len(s) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		a, b := s[i], prefix[i]
		if a >= 'A' && a <= 'Z' { a += 32 }
		if b >= 'A' && b <= 'Z' { b += 32 }
		if a != b {
			return false
		}
	}
	return true
}
