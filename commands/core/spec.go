package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/versenilvis/iris/integration/shell"
)

var DebugWriter io.Writer

func debugLog(format string, a ...interface{}) {
	if DebugWriter != nil {
		fmt.Fprintf(DebugWriter, format+"\n", a...)
	}
}

type GeneratorFunc func(tokens []string, prefix string, partial string) []Suggestion

// Spec defines a top-level CLI command
type Spec struct {
	Name        string
	Aliases     []string
	Description string
	Subcommands []Subcommand
	Options     []Option
	Generator   GeneratorFunc
	MaxArgs     int // 0 means unlimited
}

// Subcommand contains nested subcommands
type Subcommand struct {
	Name        string
	Aliases     []string
	Description string
	Subcommands []Subcommand
	Options     []Option
	Generator   GeneratorFunc
	MaxArgs     int // 0 means unlimited
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
)

// scanExternalCommands populates pathCmds
func scanExternalCommands() {
	scanPath()
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
	if shell.Current != nil {
		shellAliases = shell.Current.ScanAliases()
	} else {
		shellAliases = make(map[string]string)
	}
	pathCmds = make(map[string]bool)
	scanExternalCommands()

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
		query := tokens[0]
		results := topLevelSuggestions(query)

		// If query is a full match for a root command, also include its sub-items
		if spec, exists := registry[query]; exists {
			// ONLY suggest parameters if there's a trailing space in the query
			// or if we are filtering. But for 'just' (len=1), if no space, wait.
			hasTrailingSpace := query != "" && query[len(query)-1] == ' '

			if hasTrailingSpace {
				partial := ""
				prefix := query

				// Subcommands
				for _, sub := range spec.Subcommands {
					results = append(results, Suggestion{
						Cmd: strings.TrimSpace(query) + " " + sub.Name, Desc: sub.Description, Icon: query,
					})
				}
				// Generator (recipes, branches, etc.)
				if spec.Generator != nil {
					genResults := spec.Generator(tokens, prefix, partial)
					for _, g := range genResults {
						results = append(results, Suggestion{
							Cmd: strings.TrimSpace(query) + " " + g.Cmd, Desc: g.Desc, Icon: query,
						})
					}
				}
			}
		}
		return results
	}

	rootCmdName := tokens[0]
	spec, exists := registry[rootCmdName]
	if !exists {
		return nil
	}

	currentSubs, currentOpts, currentGen := spec.Subcommands, spec.Options, spec.Generator
	depth := 1

	for depth < len(tokens)-1 {
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
		// If not a subcommand, it's an argument. Stop traversal here.
		break
	}

	results := []Suggestion{}

	// Track current max args
	// Re-track the limit from the actual deepest subcommand reached
	currentLimit := spec.MaxArgs
	tempSubs := spec.Subcommands
	for i := 1; i < depth; i++ {
		tok := tokens[i]
		for _, sub := range tempSubs {
			if sub.Name == tok {
				currentLimit = sub.MaxArgs
				tempSubs = sub.Subcommands
				break
			}
		}
	}

	// Count arguments ALREADY TYPED for this subcommand/spec
	argCount := 0
	for i := depth; i < len(tokens)-1; i++ {
		t := tokens[i]
		// Skip empty tokens or flags
		if t != "" && !strings.HasPrefix(t, "-") && !strings.Contains(t, "=") {
			argCount++
		}
	}

	partial := tokens[len(tokens)-1]
	// If we hit the limit, don't allow ANY more argument suggestions
	// (only flags/options which start with '-' are allowed later)
	allowMoreArgs := currentLimit <= 0 || argCount < currentLimit

	debugLog("[Core] query tokens: %v (partial: '%s')", tokens, partial)
	debugLog("[Core] depth: %d, argCount: %d, limit: %d, allowMore: %v", depth, argCount, currentLimit, allowMoreArgs)

	rootCmdName = tokens[0]

	// 'prefix' is the breadcrumb for the command (e.g. "git checkout")
	prefixBuilder := strings.Builder{}
	for i := 0; i < depth; i++ {
		if i > 0 {
			prefixBuilder.WriteByte(' ')
		}
		prefixBuilder.WriteString(tokens[i])
	}
	prefix := prefixBuilder.String()

	// 'linePrefix' is the full context except the last partial word (e.g. "just run build")
	linePrefixBuilder := strings.Builder{}
	for i := 0; i < len(tokens)-1; i++ {
		if i > 0 {
			linePrefixBuilder.WriteByte(' ')
		}
		linePrefixBuilder.WriteString(tokens[i])
	}
	linePrefix := linePrefixBuilder.String()

	// already have partial from above
	rootCmdName = tokens[0]

	if currentGen != nil && allowMoreArgs {
		genResults := currentGen(tokens, prefix, partial)

		for _, g := range genResults {
			// Basic filtering based on partial input
			if partial != "" && !hasPrefix(g.Cmd, partial) && !strings.Contains(g.Cmd, partial) {
				continue
			}

			finalCmd := g.Cmd
			// Smart Absolute Building: join fragments to the line context if needed.
			// Only add linePrefix if the command doesn't already seem to have it.
			cleanLinePrefix := strings.TrimSpace(linePrefix)
			if cleanLinePrefix != "" && !strings.HasPrefix(finalCmd, cleanLinePrefix) && !strings.HasPrefix(finalCmd, rootCmdName) {
				finalCmd = cleanLinePrefix + " " + g.Cmd
			}

			// UNIVERSAL DEDUPLICATION:
			// Check if the NEW word being suggested is already present in the typed tokens.
			newTokens := tokenize(finalCmd)
			if len(newTokens) > 0 {
				lastToken := newTokens[len(newTokens)-1]
				isDuplicate := false
				for i := 0; i < len(tokens)-1; i++ {
					if tokens[i] == lastToken {
						isDuplicate = true
						break
					}
				}
				if isDuplicate {
					continue
				}
			}

			results = append(results, Suggestion{
				Cmd:  finalCmd,
				Desc: g.Desc,
				Icon: rootCmdName,
			})
		}
	}

	if allowMoreArgs {
		for _, sub := range currentSubs {
			if partial == "" || hasPrefix(sub.Name, partial) {
				results = append(results, Suggestion{
					Cmd: prefix + " " + sub.Name, Desc: sub.Description, Icon: rootCmdName,
				})
			}
		}
	}

	// Only suggest options if partial starts with '-'
	if len(partial) > 0 && partial[0] == '-' {
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
	results, seen := []Suggestion{}, make(map[string]bool)

	// 1. shell aliases (User's choice precedence)
	for name, target := range shellAliases {
		if !seen[name] && (query == "" || hasPrefix(name, query)) {
			results = append(results, Suggestion{
				Cmd: target, Desc: "alias: " + name, Icon: "root",
				// Cmd: target, Desc: "alias: " + name, Icon: "DEBUG", //--> this is just for debug on reloading, dont mind it
			})
			seen[name] = true
		}
	}

	// 2. Manual specs (High Priority)
	for name, spec := range registry {
		if seen[name] {
			continue
		}
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
