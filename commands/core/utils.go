package core

import (
	"fmt"
	"io"
	"strings"
)

var DebugWriter io.Writer

func debugLog(format string, a ...interface{}) {
	if DebugWriter != nil {
		fmt.Fprintf(DebugWriter, format+"\n", a...)
	}
}

// SplitAliasTokens parses the input string into shell-like tokens handling quotes
// example: SplitAliasTokens("git commit -m \"hello world\"")
func Tokenize(s string) []string {
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
			}
		default:
			current.WriteRune(c)
		}
	}

	tokens = append(tokens, current.String())
	return tokens
}

// HasPrefix checks if s starts with prefix using case-insensitive matching
func HasPrefix(s, prefix string) bool {
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

// CI = case insensitive
func HasPrefixCI(s, prefix string) bool {
	return HasPrefix(s, prefix)
}
