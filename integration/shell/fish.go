package shell

import (
	"maps"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type fishDefs struct {
	aliases map[string]string
	abbrs   map[string]string
}

// Aliases and abbreviations share one entry: scanning them separately would
// make each call evict the other's result, re-reading the tree every keystroke.
var fishScanCache struct {
	sync.Mutex
	dir    string
	stamps map[string]fileStamp
	defs   fishDefs
}

func scanFishDefs(dir string) fishDefs {
	fishScanCache.Lock()
	defer fishScanCache.Unlock()

	if fishScanCache.defs.aliases != nil && fishScanCache.dir == dir && stampsUnchanged(fishScanCache.stamps) {
		return cloneFishDefs(fishScanCache.defs)
	}

	defs := fishDefs{aliases: map[string]string{}, abbrs: map[string]string{}}
	stamps := make(map[string]fileStamp)
	visited := make(map[string]bool)
	walk := func(data string, onSource func(target string)) {
		walkFishConfig(dir, data, defs, onSource)
	}

	files, dirs := fishConfigFiles(dir)
	// A file added to conf.d/ leaves every scanned file untouched, so only the
	// directory's own stamp can invalidate the cache.
	for _, d := range dirs {
		stamps[d], _ = statStamp(d)
	}
	for _, path := range files {
		scanConfigFile(path, walk, visited, stamps, 0)
	}

	fishScanCache.dir = dir
	fishScanCache.stamps = stamps
	fishScanCache.defs = defs

	return cloneFishDefs(defs)
}

func cloneFishDefs(defs fishDefs) fishDefs {
	return fishDefs{aliases: maps.Clone(defs.aliases), abbrs: maps.Clone(defs.abbrs)}
}

func fishConfigFiles(dir string) ([]string, []string) {
	functions := filepath.Join(dir, "functions")
	confd := filepath.Join(dir, "conf.d")

	// An autoloaded function only loads while its name is still undefined, so
	// an `alias --save` file loses to whatever conf.d and config.fish define.
	files := fishFiles(functions)
	files = append(files, fishFiles(confd)...)
	files = append(files, filepath.Join(dir, "config.fish"))

	return files, []string{functions, confd}
}

func fishFiles(dir string) []string {
	matches, err := filepath.Glob(filepath.Join(dir, "*.fish"))
	if err != nil {
		return nil
	}
	return matches
}

type fishBlock struct {
	branching  bool
	arms       int
	suppressed bool
}

func walkFishConfig(dir, data string, defs fishDefs, onSource func(target string)) {
	// fish closes every block with `end`, so `for`/`function`/`begin` need
	// tracking too or their `end` pops the enclosing `if` and the fallback arm
	// starts overwriting the preferred one.
	var stack []fishBlock
	altArm := 0

	vars := map[string][]string{}
	if dir != "" {
		vars["__fish_config_dir"] = []string{dir}
	}

	for _, line := range fishLogicalLines(data) {
		for _, segment := range splitShellSegments(stripFishComment(line)) {
			segment = trimFishModifiers(strings.TrimSpace(segment))
			if segment == "" {
				continue
			}

			switch firstWord(segment) {
			case "else", "case":
				if n := len(stack) - 1; n >= 0 {
					stack[n].arms++
					if stack[n].branching && stack[n].arms > 1 && !stack[n].suppressed {
						stack[n].suppressed = true
						altArm++
					}
				}
			case "end":
				if n := len(stack) - 1; n >= 0 {
					if stack[n].suppressed {
						altArm--
					}
					stack = stack[:n]
				}
			case "if":
				// the `if` line is already the first arm; `switch` has none
				// until its first `case`
				stack = append(stack, fishBlock{branching: true, arms: 1})
			case "switch":
				stack = append(stack, fishBlock{branching: true})
			case "for", "while", "begin", "function":
				stack = append(stack, fishBlock{})
			}

			parseFishSet(segment, vars)
			parseFishFor(segment, vars)

			if onSource != nil {
				for _, target := range fishSourceTargets(segment, vars) {
					onSource(target)
				}
			}

			if altArm > 0 {
				continue
			}
			parseFishAlias(segment, defs.aliases)
			parseFishAliasFunction(segment, defs.aliases)
			parseFishAbbr(segment, defs.abbrs)
		}
	}
}

func ParseFishConfig(dir, data string) (map[string]string, map[string]string) {
	defs := fishDefs{aliases: map[string]string{}, abbrs: map[string]string{}}
	walkFishConfig(dir, data, defs, nil)
	return defs.aliases, defs.abbrs
}

// A `\` at end of line continues the statement, and real configs use it to
// spread a `set` list over many lines.
func fishLogicalLines(data string) []string {
	var lines []string
	var joined strings.Builder

	for line := range strings.SplitSeq(data, "\n") {
		trimmed := strings.TrimRight(line, " \t")
		if strings.HasSuffix(trimmed, `\`) && !strings.HasSuffix(trimmed, `\\`) {
			joined.WriteString(strings.TrimSuffix(trimmed, `\`))
			joined.WriteString(" ")
			continue
		}
		joined.WriteString(line)
		lines = append(lines, joined.String())
		joined.Reset()
	}
	if joined.Len() > 0 {
		lines = append(lines, joined.String())
	}
	return lines
}

func firstWord(segment string) string {
	if i := strings.IndexAny(segment, " \t"); i >= 0 {
		return segment[:i]
	}
	return segment
}

func stripFishComment(line string) string {
	inQuote := false
	var quoteChar rune
	// fish only starts a comment at a word boundary, so `echo foo#bar` keeps
	// its hash.
	atWordStart := true

	for i, c := range line {
		switch {
		case !inQuote && (c == '"' || c == '\''):
			inQuote = true
			quoteChar = c
		case inQuote && c == quoteChar:
			inQuote = false
		case !inQuote && c == '#' && atWordStart:
			return line[:i]
		}
		atWordStart = c == ' ' || c == '\t'
	}
	return line
}

func trimFishModifiers(segment string) string {
	modifiers := map[string]bool{"and": true, "or": true, "not": true, "command": true, "builtin": true}
	for modifiers[firstWord(segment)] {
		_, rest, _ := strings.Cut(segment, " ")
		segment = strings.TrimSpace(rest)
	}
	return segment
}

func cutFishCommand(segment, name string) (string, bool) {
	if firstWord(segment) != name {
		return "", false
	}
	rest := strings.TrimSpace(segment[len(name):])
	return rest, rest != ""
}

func fishSourceTargets(segment string, vars map[string][]string) []string {
	rest, ok := cutFishCommand(segment, "source")
	if !ok {
		if rest, ok = cutFishCommand(segment, "."); !ok {
			return nil
		}
	}

	fields := SplitAliasTokens(rest)
	// fish spells command substitution as a bare (cmd), which expandConfigPath
	// has no way to recognize.
	if len(fields) == 0 || strings.Contains(fields[0], "(") {
		return nil
	}

	candidates, ok := expandFishToken(fields[0], vars)
	if !ok {
		return nil
	}

	targets := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		if path, ok := expandConfigPath(candidate); ok {
			targets = append(targets, path)
		}
	}
	return targets
}

// Config files routinely reach their real content through `set`/`for` rather
// than a literal path, so those two have to be resolved or `source` finds
// nothing to follow.
func parseFishSet(segment string, vars map[string][]string) {
	rest, ok := cutFishCommand(segment, "set")
	if !ok {
		return
	}

	reject := map[string]bool{
		"-e": true, "--erase": true, "-q": true, "--query": true,
		"-S": true, "--show": true, "-h": true, "--help": true,
	}
	positional, ok := fishPositionals(SplitAliasTokens(rest), nil, reject)
	if !ok || len(positional) == 0 {
		return
	}

	name := unquoteAliasValue(positional[0])
	if name == "" || strings.ContainsAny(name, "$[") {
		return
	}

	values := make([]string, 0, len(positional)-1)
	for _, token := range positional[1:] {
		expanded, ok := expandFishToken(token, vars)
		if !ok {
			delete(vars, name)
			return
		}
		values = append(values, expanded...)
	}
	vars[name] = values
}

func parseFishFor(segment string, vars map[string][]string) {
	rest, ok := cutFishCommand(segment, "for")
	if !ok {
		return
	}

	tokens := SplitAliasTokens(rest)
	if len(tokens) < 3 || tokens[1] != "in" {
		return
	}

	name := unquoteAliasValue(tokens[0])
	if name == "" {
		return
	}

	// every iteration is followed at once, since which one the shell takes
	// cannot be known from the text
	values := make([]string, 0, len(tokens)-2)
	for _, token := range tokens[2:] {
		expanded, ok := expandFishToken(token, vars)
		if !ok {
			delete(vars, name)
			return
		}
		values = append(values, expanded...)
	}
	vars[name] = values
}

func expandFishToken(token string, vars map[string][]string) ([]string, bool) {
	// fish leaves single-quoted text alone
	if strings.HasPrefix(token, "'") {
		return []string{unquoteAliasValue(token)}, true
	}
	return expandFishWord(unquoteAliasValue(token), vars)
}

const maxFishExpansion = 64

func expandFishWord(word string, vars map[string][]string) ([]string, bool) {
	results := []string{""}

	for {
		before, after, ok := strings.Cut(word, "$")
		if !ok {
			for i := range results {
				results[i] += word
			}
			return results, true
		}

		name, rest := cutFishVarName(after)
		values, ok := fishVarValues(name, vars)
		if !ok {
			return nil, false
		}

		head := before
		expanded := make([]string, 0, len(results)*len(values))
		for _, prefix := range results {
			for _, value := range values {
				expanded = append(expanded, prefix+head+value)
			}
		}
		if len(expanded) > maxFishExpansion {
			return nil, false
		}
		results, word = expanded, rest
	}
}

func fishVarValues(name string, vars map[string][]string) ([]string, bool) {
	if name == "" {
		return nil, false
	}
	if values, ok := vars[name]; ok {
		return values, len(values) > 0
	}
	if value, ok := os.LookupEnv(name); ok && value != "" {
		return []string{value}, true
	}
	return nil, false
}

func cutFishVarName(s string) (string, string) {
	s = strings.TrimPrefix(s, "{")
	end := 0
	for end < len(s) && (s[end] == '_' ||
		(s[end] >= 'a' && s[end] <= 'z') ||
		(s[end] >= 'A' && s[end] <= 'Z') ||
		(s[end] >= '0' && s[end] <= '9')) {
		end++
	}
	return s[:end], strings.TrimPrefix(s[end:], "}")
}

func fishPositionals(tokens []string, valueFlags, reject map[string]bool) ([]string, bool) {
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		if token == "--" {
			return tokens[i+1:], true
		}
		if len(token) < 2 || !strings.HasPrefix(token, "-") {
			return tokens[i:], true
		}

		flag, _, hasValue := strings.Cut(token, "=")
		if reject[flag] {
			return nil, false
		}
		if !hasValue && valueFlags[flag] {
			i++
		}
	}
	return nil, true
}

func parseFishAlias(segment string, aliases map[string]string) {
	rest, ok := cutFishCommand(segment, "alias")
	if !ok {
		return
	}

	tokens, ok := fishPositionals(SplitAliasTokens(rest), nil, nil)
	if !ok || len(tokens) == 0 || !isLiteralFishToken(tokens[0]) {
		return
	}

	name, target := fishAliasNameAndTarget(tokens)
	if name == "" || target == "" {
		return
	}
	aliases[name] = target
}

// A name the shell computes -- `(string split ...)`, `$tool` -- is not a name
// iris can suggest. Single quotes make it literal, so `'!$'` still counts.
func isLiteralFishToken(token string) bool {
	if strings.HasPrefix(token, "'") {
		return true
	}
	return !strings.ContainsAny(token, "$()")
}

func fishAliasNameAndTarget(tokens []string) (string, string) {
	if len(tokens) == 0 {
		return "", ""
	}

	parts := make([]string, 0, len(tokens))
	for _, token := range tokens {
		parts = append(parts, unquoteAliasValue(token))
	}

	// `alias name=value` still spills the rest of the value into later tokens
	if name, head, ok := strings.Cut(parts[0], "="); ok {
		parts[0] = unquoteAliasValue(head)
		return strings.TrimSpace(name), strings.TrimSpace(strings.Join(parts, " "))
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(strings.Join(parts[1:], " "))
}

func parseFishAbbr(segment string, abbrs map[string]string) {
	rest, ok := cutFishCommand(segment, "abbr")
	if !ok {
		return
	}

	// --regex leaves no literal name and --function builds the expansion at
	// runtime, so neither can be recovered from static text.
	reject := map[string]bool{
		"-r": true, "--regex": true, "-f": true, "--function": true,
		"-e": true, "--erase": true, "--rename": true,
		"-q": true, "--query": true, "-s": true, "--show": true,
		"-l": true, "--list": true, "-h": true, "--help": true,
	}
	valueFlags := map[string]bool{
		"--position": true, "-c": true, "--command": true, "--color": true,
	}

	tokens, ok := fishPositionals(SplitAliasTokens(rest), valueFlags, reject)
	if !ok || len(tokens) < 2 || !isLiteralFishToken(tokens[0]) {
		return
	}

	name := strings.TrimSpace(unquoteAliasValue(tokens[0]))
	parts := make([]string, 0, len(tokens)-1)
	for _, token := range tokens[1:] {
		parts = append(parts, unquoteAliasValue(token))
	}

	expansion := strings.TrimSpace(strings.Join(parts, " "))
	if name == "" || expansion == "" {
		return
	}
	abbrs[name] = expansion
}

// `alias --save` is the only writer of functions/, and it records the original
// definition verbatim in the function's description.
func parseFishAliasFunction(segment string, aliases map[string]string) {
	rest, ok := cutFishCommand(segment, "function")
	if !ok {
		return
	}

	tokens := SplitAliasTokens(rest)
	for i, token := range tokens {
		flag, value, hasValue := strings.Cut(token, "=")
		if flag != "--description" && flag != "-d" {
			continue
		}
		if !hasValue {
			if i+1 >= len(tokens) {
				return
			}
			value = tokens[i+1]
		}
		parseFishAlias(unquoteAliasValue(value), aliases)
		return
	}
}
