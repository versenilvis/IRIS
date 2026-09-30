package ctxcheck

import (
	"strings"
	"unicode"
)

type Segment struct {
	Tokens      []string
	IsUnknown   bool
	IsCwdChange bool
}

type ParsedCommand struct {
	Segments  []Segment
	IsUnknown bool
}

var wrappers = map[string]bool{
	"env":     true,
	"sudo":    true,
	"time":    true,
	"command": true,
	"nohup":   true,
	"exec":    true,
}

func isEnvAssignment(token string) bool {
	idx := strings.IndexByte(token, '=')
	if idx <= 0 {
		return false
	}
	key := token[:idx]
	for i, r := range key {
		if i == 0 {
			if r != '_' && !unicode.IsLetter(r) {
				return false
			}
		} else {
			if r != '_' && !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				return false
			}
		}
	}
	return true
}

func Parse(cmd string, d Dialect) ParsedCommand {
	if d == UnknownDialect {
		return ParsedCommand{IsUnknown: true}
	}
	rawSegments, globalUnknown := splitIntoSegments(cmd, d)
	if globalUnknown {
		return ParsedCommand{IsUnknown: true}
	}

	var segments []Segment
	hasCwdChange := false

	for _, raw := range rawSegments {
		if raw.isUnknown {
			segments = append(segments, Segment{IsUnknown: true})
			continue
		}
		if len(raw.tokens) == 0 {
			continue
		}

		tokens := raw.tokens
		if d == Fish && len(tokens) > 0 && (tokens[0] == "and" || tokens[0] == "or") {
			tokens = tokens[1:]
		}
		for len(tokens) > 0 && isEnvAssignment(tokens[0]) {
			tokens = tokens[1:]
		}

		hasUnsupportedWrapper := false
		for len(tokens) > 0 && wrappers[tokens[0]] {
			tokens = tokens[1:]
			if len(tokens) > 0 && strings.HasPrefix(tokens[0], "-") {
				hasUnsupportedWrapper = true
				break
			}
		}

		if hasUnsupportedWrapper {
			segments = append(segments, Segment{IsUnknown: true})
			continue
		}

		if len(tokens) == 0 {
			continue
		}

		cmdWord := tokens[0]
		if cmdWord == "cd" || cmdWord == "pushd" || cmdWord == "popd" {
			hasCwdChange = true
		}

		segments = append(segments, Segment{
			Tokens:    tokens,
			IsUnknown: raw.isUnknown,
		})
	}

	if len(segments) > 1 && hasCwdChange {
		return ParsedCommand{IsUnknown: true}
	}

	// single pushd / popd is unknown because it manipulates directory stack
	if len(segments) == 1 && len(segments[0].Tokens) > 0 && (segments[0].Tokens[0] == "pushd" || segments[0].Tokens[0] == "popd") {
		return ParsedCommand{IsUnknown: true}
	}

	return ParsedCommand{Segments: segments}
}

type rawSegment struct {
	tokens    []string
	isUnknown bool
}

func splitIntoSegments(cmd string, d Dialect) ([]rawSegment, bool) {
	var segments []rawSegment
	var currentTokens []string
	var currentToken strings.Builder

	inSingle := false
	inDouble := false
	hasSubstitution := false
	hasRedirect := false

	n := len(cmd)
	i := 0

	flushToken := func() {
		if currentToken.Len() > 0 {
			currentTokens = append(currentTokens, currentToken.String())
			currentToken.Reset()
		}
	}

	flushSegment := func() {
		flushToken()
		if len(currentTokens) > 0 || hasSubstitution || hasRedirect {
			segUnknown := hasSubstitution || hasRedirect
			segments = append(segments, rawSegment{
				tokens:    currentTokens,
				isUnknown: segUnknown,
			})
			currentTokens = nil
			hasSubstitution = false
			hasRedirect = false
		}
	}

	for i < n {
		b := cmd[i]

		if inSingle {
			if d == Posix {
				if b == '\'' {
					inSingle = false
				} else {
					currentToken.WriteByte(b)
				}
				i++
				continue
			}
			// fish dialect: \' and \\ escape
			if b == '\\' && i+1 < n && (cmd[i+1] == '\'' || cmd[i+1] == '\\') {
				currentToken.WriteByte(cmd[i+1])
				i += 2
				continue
			}
			if b == '\'' {
				inSingle = false
			} else {
				currentToken.WriteByte(b)
			}
			i++
			continue
		}

		if inDouble {
			if b == '\\' && i+1 < n {
				// escape inside double quotes
				currentToken.WriteByte(cmd[i+1])
				i += 2
				continue
			}
			if b == '"' {
				inDouble = false
				i++
				continue
			}
			if b == '`' || (b == '$' && i+1 < n && cmd[i+1] == '(') {
				hasSubstitution = true
			}
			currentToken.WriteByte(b)
			i++
			continue
		}

		// outside quotes
		if b == '\'' {
			inSingle = true
			i++
			continue
		}
		if b == '"' {
			inDouble = true
			i++
			continue
		}

		// substitutions
		if b == '`' || (b == '$' && i+1 < n && cmd[i+1] == '(') {
			hasSubstitution = true
			currentToken.WriteByte(b)
			i++
			continue
		}
		if d == Fish && b == '(' {
			hasSubstitution = true
			currentToken.WriteByte(b)
			i++
			continue
		}

		// separators: ;, &&, ||, |, &, \n
		if b == ';' || b == '\n' {
			flushSegment()
			i++
			continue
		}
		if b == '&' {
			if i+1 < n && cmd[i+1] == '&' {
				flushSegment()
				i += 2
				continue
			}
			// single & can be background or redirect like 2>&1
			// check if part of 2>&1
			if currentToken.String() == "2>" && i+2 < n && cmd[i+1] == '&' && cmd[i+2] == '1' {
				currentToken.WriteString("&1")
				i += 2
				continue
			}
			flushSegment()
			i++
			continue
		}
		if b == '|' {
			if i+1 < n && cmd[i+1] == '|' {
				flushSegment()
				i += 2
				continue
			}
			flushSegment()
			i++
			continue
		}

		// redirects
		if b == '>' || b == '<' {
			// check 2>&1
			isAllowedRedirect := false
			if b == '>' && currentToken.String() == "2" && i+2 < n && cmd[i+1] == '&' && cmd[i+2] == '1' {
				isAllowedRedirect = true
				currentToken.WriteString(">&1")
				i += 3
				continue
			}
			if !isAllowedRedirect {
				hasRedirect = true
			}
			i++
			continue
		}

		// whitespace
		if unicode.IsSpace(rune(b)) {
			tokenStr := currentToken.String()
			// fish words 'and' / 'or' as separators or leading connectors
			if d == Fish && (tokenStr == "and" || tokenStr == "or") {
				currentToken.Reset()
				if len(currentTokens) > 0 {
					flushSegment()
				}
				i++
				continue
			}
			flushToken()
			i++
			continue
		}

		// backslash escape outside quotes
		if b == '\\' && i+1 < n {
			currentToken.WriteByte(cmd[i+1])
			i += 2
			continue
		}

		currentToken.WriteByte(b)
		i++
	}

	// trailing fish 'and'/'or' check
	if d == Fish {
		tokenStr := currentToken.String()
		if tokenStr == "and" || tokenStr == "or" {
			currentToken.Reset()
			flushSegment()
			return segments, false
		}
	}

	flushSegment()
	return segments, false
}
