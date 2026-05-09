package core

import (
	"strings"

	"github.com/versenilvis/iris/integration/shell"
)

var ShellAliases = map[string]string{}

func GetAlias(name string) (string, bool) {
	val, ok := ShellAliases[name]
	return val, ok
}

// Lookup finds matching suggestions for your input by looking at how many words you typed
// it changes aliases to real commands, finds subcommands inside others, and runs generators for more suggestions
// e.g. Lookup("git che") -> suggests "git checkout"
// e.g. Lookup("git checkout ") -> suggests branch names via generator
func Lookup(input string) []Suggestion {
	if shell.Current != nil {
		ShellAliases = shell.Current.ScanAliases()
	}

	if input == "" {
		return nil
	}

	tokens := Tokenize(input)

	if len(tokens) == 1 && tokens[0] == "" {
		return nil
	}
	// NOTE: remember that wrapper already check if we type nothing, but I just want to make sure
	// in the future, maybe we will write unit test or using Lookup in another module
	// it will be safer because we maybe forget to check it in that module

	pathCmds = make(map[string]bool)
	scanExternalCommands()

	// if you have an alias in your shell config like: alias gca="git commit -a"
	// if the first word match it, IRIS will suggest "git commit -a"
	if len(tokens) > 1 {
		if target, ok := ShellAliases[tokens[0]]; ok {
			aliasTokens := Tokenize(target)
			if len(aliasTokens) > 0 && aliasTokens[len(aliasTokens)-1] == "" {
				aliasTokens = aliasTokens[:len(aliasTokens)-1]
			}
			tokens = append(aliasTokens, tokens[1:]...)
		}
	}

	if len(tokens) == 1 {
		query := tokens[0]
		results := topLevelSuggestions(query)

		if spec, exists := Registry[query]; exists {
			hasTrailingSpace := query != "" && query[len(query)-1] == ' '

			if hasTrailingSpace {
				partial := ""
				prefix := query

				for _, sub := range spec.Subcommands {
					results = append(results, Suggestion{
						Cmd: strings.TrimSpace(query) + " " + sub.Name, Desc: sub.Description, Icon: query,
					})
				}
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
	spec, exists := Registry[rootCmdName]
	debugLog("[core] lookup tokens: %v, registry exists: %v", tokens, exists)
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
		break
	}

	results := []Suggestion{}

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

	argCount := 0
	for i := depth; i < len(tokens)-1; i++ {
		t := tokens[i]
		if t != "" && !strings.HasPrefix(t, "-") && !strings.Contains(t, "=") {
			argCount++
		}
	}

	partial := tokens[len(tokens)-1]
	allowMoreArgs := currentLimit <= 0 || argCount < currentLimit

	debugLog("[core] query tokens: %v (partial: '%s')", tokens, partial)
	debugLog("[core] depth: %d, argCount: %d, limit: %d, allowMore: %v", depth, argCount, currentLimit, allowMoreArgs)

	prefixBuilder := strings.Builder{}
	for i := 0; i < depth; i++ {
		if i > 0 {
			prefixBuilder.WriteByte(' ')
		}
		prefixBuilder.WriteString(tokens[i])
	}
	prefix := prefixBuilder.String()

	linePrefixBuilder := strings.Builder{}
	for i := 0; i < len(tokens)-1; i++ {
		if i > 0 {
			linePrefixBuilder.WriteByte(' ')
		}
		linePrefixBuilder.WriteString(tokens[i])
	}
	linePrefix := linePrefixBuilder.String()

	if currentGen != nil && allowMoreArgs {
		genResults := currentGen(tokens, prefix, partial)

		for _, g := range genResults {
			if partial != "" && !HasPrefix(g.Cmd, partial) && !strings.Contains(g.Cmd, partial) {
				continue
			}

			suggested := g.Cmd
			if strings.Contains(suggested, " ") && !strings.HasPrefix(suggested, "\"") {
				suggested = "\"" + suggested + "\""
			}

			// if the suggestion is a full path that includes 
			// words already in the command line (multi-word support), we replace 
			// the entire argument part by using prefix only
			finalCmd := ""
			if len(tokens) > depth+1 && strings.HasPrefix(g.Cmd, tokens[depth]) {
				finalCmd = prefix + " " + suggested
			} else {
				finalCmd = strings.TrimSpace(linePrefix) + " " + suggested
			}

			newTokens := Tokenize(finalCmd)
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
			if partial == "" || HasPrefix(sub.Name, partial) {
				results = append(results, Suggestion{
					Cmd: prefix + " " + sub.Name, Desc: sub.Description, Icon: rootCmdName,
				})
			}
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
			if !usedOpts[opt.Name] && (partial == "" || HasPrefix(opt.Name, partial)) {
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

	for name, target := range ShellAliases {
		if !seen[name] && (query == "" || HasPrefix(name, query)) {
			results = append(results, Suggestion{
				Cmd: target, Desc: "alias: " + name, Icon: "root",
			})
			seen[name] = true
		}
	}

	for name, spec := range Registry {
		if seen[name] {
			continue
		}
		match := false
		if query == "" || HasPrefix(name, query) {
			match = true
		} else {
			for _, a := range spec.Aliases {
				if HasPrefix(a, query) {
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

	for name := range pathCmds {
		if !seen[name] && (query == "" || HasPrefix(name, query)) {
			results = append(results, Suggestion{
				Cmd: name, Desc: "system command", Icon: "root",
			})
		}
	}

	return results
}
