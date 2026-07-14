package spec

import (
	"slices"
	"strings"
)

// TryExtractSkeleton walks the registered subcommand tree for the given command buffer.
// It stops when a token doesn't match any registered subcommand name (treating that token as an argument/value)
// or when encountering a flag ('-' prefix).
// Returns the normalized subcommand skeleton (e.g. "git checkout feature-x" -> "git checkout") and true if a spec exists.
func TryExtractSkeleton(buf string) (string, bool) {
	buf = strings.TrimSpace(buf)
	if buf == "" {
		return "", false
	}

	aliases := GetAliasesCopy()
	tokens := Tokenize(buf)
	// filter out empty tokens (e.g. from trailing space)
	var cleanTokens []string
	for _, t := range tokens {
		if t != "" {
			cleanTokens = append(cleanTokens, t)
		}
	}
	if len(cleanTokens) == 0 {
		return "", false
	}

	// expand shell alias for the root command if present
	if target, ok := aliases[cleanTokens[0]]; ok {
		aliasTokens := Tokenize(target)
		var cleanAliasTokens []string
		for _, t := range aliasTokens {
			if t != "" {
				cleanAliasTokens = append(cleanAliasTokens, t)
			}
		}
		cleanTokens = append(cleanAliasTokens, cleanTokens[1:]...)
	}

	if len(cleanTokens) == 0 {
		return "", false
	}

	rootName := cleanTokens[0]
	spec, exists := Registry[rootName]
	if !exists || spec == nil {
		return "", false
	}

	skeletonTokens := []string{rootName}
	currentSubs := spec.Subcommands

	for i := 1; i < len(cleanTokens); i++ {
		tok := cleanTokens[i]
		if strings.HasPrefix(tok, "-") || strings.Contains(tok, "=") {
			continue
		}

		found := false
		for _, sub := range currentSubs {
			if sub.Name == tok || slices.Contains(sub.Aliases, tok) {
				skeletonTokens = append(skeletonTokens, sub.Name)
				currentSubs = sub.Subcommands
				found = true
				break
			}
		}

		if found {
			continue
		}

		// If the previous token was a flag (started with '-'), this token might be the flag's argument/value (e.g. -C /tmp).
		// In that case, we skip this token and continue looking for subcommands in subsequent tokens.
		if i > 0 && strings.HasPrefix(cleanTokens[i-1], "-") {
			continue
		}

		// token did not match any registered subcommand and is not a flag argument -> must be a positional argument (branch, file, etc.)
		break
	}

	return strings.Join(skeletonTokens, " "), true
}
