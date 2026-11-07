package synthesize

import (
	"bufio"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/versenilvis/iris/spec"
)

var (
	reBashAssign = regexp.MustCompile(`(?i)(?:opts|options|boolean_options|options_with_args|flags)\s*=\s*(?:"([^"]*)"|'([^']*)')`)
	reCompgenW   = regexp.MustCompile(`compgen\s+-W\s+(?:"([^"]*)"|'([^']*)')`)
	reCaseBranch = regexp.MustCompile(`^\s*([^\s)]+(?:\s*\|\s*[^\s)]+)*)\s*\)`)
)

var bashBuiltinTestFlags = map[string]bool{
	"-eq": true, "-ne": true, "-lt": true, "-le": true, "-gt": true, "-ge": true,
	"-nt": true, "-ot": true, "-ef": true,
}

func parseBashReader(r io.Reader) ([]spec.UserOptionJSON, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var options []spec.UserOptionJSON
	seenFlags := make(map[string]int)

	recordOption := func(name string, aliases []string) {
		name = sanitizeText(name)
		if !reValidFlag.MatchString(name) || bashBuiltinTestFlags[name] {
			return
		}
		var cleanAliases []string
		for _, a := range aliases {
			a = sanitizeText(a)
			if reValidFlag.MatchString(a) && !bashBuiltinTestFlags[a] && a != name {
				cleanAliases = append(cleanAliases, a)
			}
		}

		if idx, exists := seenFlags[name]; exists {
			for _, a := range cleanAliases {
				if _, aliasExists := seenFlags[a]; !aliasExists {
					options[idx].Aliases = append(options[idx].Aliases, a)
					seenFlags[a] = idx
				}
			}
			return
		}

		idx := len(options)
		options = append(options, spec.UserOptionJSON{
			Name:    name,
			Aliases: cleanAliases,
		})
		seenFlags[name] = idx
		for _, a := range cleanAliases {
			seenFlags[a] = idx
		}
	}

	inCase := false
	var multiLineAssign strings.Builder
	var inMultiLineQuote byte

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			continue
		}

		// track case statement branches to discover flag aliases
		if strings.HasPrefix(trimmed, "case ") {
			inCase = true
		} else if strings.HasPrefix(trimmed, "esac") {
			inCase = false
		}

		if inCase {
			if sm := reCaseBranch.FindStringSubmatch(trimmed); sm != nil {
				pattern := sm[1]
				parts := strings.Split(pattern, "|")
				var flagsInBranch []string
				for _, p := range parts {
					p = strings.Trim(strings.TrimSpace(p), `'"`)
					if reFlagLong.MatchString(p) || reFlagShort.MatchString(p) {
						if !bashBuiltinTestFlags[p] {
							flagsInBranch = append(flagsInBranch, p)
						}
					}
				}
				if len(flagsInBranch) > 0 {
					var longFlags []string
					var shortFlags []string
					for _, f := range flagsInBranch {
						if strings.HasPrefix(f, "--") {
							longFlags = append(longFlags, f)
						} else {
							shortFlags = append(shortFlags, f)
						}
					}
					if len(longFlags) > 0 {
						primary := longFlags[0]
						var aliases []string
						aliases = append(aliases, longFlags[1:]...)
						aliases = append(aliases, shortFlags...)
						recordOption(primary, aliases)
					} else if len(shortFlags) > 0 {
						recordOption(shortFlags[0], shortFlags[1:])
					}
				}
			}
		}

		for _, match := range reCompgenW.FindAllStringSubmatch(line, -1) {
			text := match[1]
			if text == "" {
				text = match[2]
			}
			extractTokensFromText(text, recordOption)
		}

		if inMultiLineQuote != 0 {
			multiLineAssign.WriteString(" ")
			multiLineAssign.WriteString(trimmed)
			if strings.ContainsRune(trimmed, rune(inMultiLineQuote)) {
				extractTokensFromText(multiLineAssign.String(), recordOption)
				multiLineAssign.Reset()
				inMultiLineQuote = 0
			}
			continue
		}

		if sm := reBashAssign.FindStringSubmatch(trimmed); sm != nil {
			text := sm[1]
			if text == "" {
				text = sm[2]
			}
			extractTokensFromText(text, recordOption)
			continue
		}

		// check for opening of multiline quoted options assignment
		if strings.Contains(trimmed, "OPTS=") || strings.Contains(trimmed, "options=") || strings.Contains(trimmed, "boolean_options=") {
			for _, quote := range []byte{'"', '\''} {
				if idx := strings.IndexByte(trimmed, quote); idx != -1 {
					rest := trimmed[idx+1:]
					if !strings.ContainsRune(rest, rune(quote)) {
						inMultiLineQuote = quote
						multiLineAssign.WriteString(rest)
						break
					}
				}
			}
		}
	}

	return options, nil
}

func extractTokensFromText(text string, record func(name string, aliases []string)) {
	tokens := strings.FieldsSeq(text)
	for tok := range tokens {
		tok = sanitizeText(tok)
		tok = strings.Trim(tok, `'"$,;()`)
		if idx := strings.Index(tok, "="); idx != -1 {
			tok = tok[:idx]
		}
		if reValidFlag.MatchString(tok) {
			record(tok, nil)
		}
	}
}

func parseBashFile(dir, cmdName string, maxSize int64) ([]spec.UserOptionJSON, error) {
	cleanName := filepath.Base(cmdName)
	reader, err := safeOpenFile(dir, cleanName, maxSize)
	if err != nil {
		// some completions start with an underscore like _git
		if readerAlt, altErr := safeOpenFile(dir, "_"+cleanName, maxSize); altErr == nil {
			reader = readerAlt
		} else {
			return nil, err
		}
	}
	defer reader.Close()

	return parseBashReader(reader)
}
