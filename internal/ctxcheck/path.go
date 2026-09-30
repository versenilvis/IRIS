package ctxcheck

import (
	"os"
	"path/filepath"
	"strings"
)

var interpreterNames = map[string]bool{
	"python":  true,
	"python3": true,
	"node":    true,
	"sh":      true,
	"bash":    true,
	"zsh":     true,
	"fish":    true,
	"ruby":    true,
	"perl":    true,
	"php":     true,
}

func expandHome(path string) string {
	if path == "~" {
		home, err := os.UserHomeDir()
		if err == nil {
			return home
		}
		return path
	}
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

func isExplicitPath(tok string) bool {
	if strings.ContainsAny(tok, "*?[") || strings.Contains(tok, "...") || strings.Contains(tok, "$") {
		return false
	}
	return strings.HasPrefix(tok, "./") ||
		strings.HasPrefix(tok, "../") ||
		strings.HasPrefix(tok, "~/") ||
		strings.HasPrefix(tok, "/")
}

func resolvePath(tok, cwd string) string {
	expanded := expandHome(tok)
	if filepath.IsAbs(expanded) {
		return filepath.Clean(expanded)
	}
	return filepath.Clean(filepath.Join(cwd, expanded))
}

func ValidateCd(tokens []string, cwd string) Verdict {
	if len(tokens) == 0 || tokens[0] != "cd" {
		return Free
	}
	if len(tokens) == 1 {
		return Free
	}

	target := tokens[1]
	if target == "~" || target == "-" {
		return Free
	}
	if strings.HasPrefix(target, "-") && target != "-" {
		if len(tokens) > 2 {
			target = tokens[2]
		} else {
			return Free
		}
	}
	if target == "~" || target == "-" {
		return Free
	}

	if strings.ContainsAny(target, "*?[") || strings.Contains(target, "$") {
		return Free
	}

	resolved := resolvePath(target, cwd)
	fi, err := os.Stat(resolved)
	if err == nil {
		if fi.IsDir() {
			return Valid
		}
		return Invalid
	}

	if !filepath.IsAbs(target) && !strings.HasPrefix(target, "~") && os.Getenv("CDPATH") != "" {
		return Unknown
	}

	return Invalid
}

func ValidatePathTokens(tokens []string, cwd string) Verdict {
	if len(tokens) == 0 {
		return Free
	}

	cmdWord := tokens[0]
	if cmdWord == "cd" {
		return ValidateCd(tokens, cwd)
	}

	// check if command in executable position is an explicit path
	if isExplicitPath(cmdWord) {
		resolved := resolvePath(cmdWord, cwd)
		if _, err := os.Stat(resolved); err != nil {
			return Invalid
		}
		return Valid
	}

	// check if command is an interpreter: check first non-flag script argument
	if interpreterNames[cmdWord] {
		for _, arg := range tokens[1:] {
			if strings.HasPrefix(arg, "-") {
				continue
			}
			if strings.ContainsAny(arg, "*?[") || strings.Contains(arg, "$") || strings.Contains(arg, "...") {
				return Free
			}
			resolved := resolvePath(arg, cwd)
			if _, err := os.Stat(resolved); err != nil {
				return Invalid
			}
			return Valid
		}
		return Free
	}

	// check other explicit path tokens in arguments
	hasExplicit := false
	for _, tok := range tokens[1:] {
		if isExplicitPath(tok) {
			hasExplicit = true
			resolved := resolvePath(tok, cwd)
			if _, err := os.Stat(resolved); err != nil {
				return Invalid
			}
		}
	}

	if hasExplicit {
		return Valid
	}
	return Free
}
