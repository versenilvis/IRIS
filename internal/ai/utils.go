package ai

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/versenilvis/iris/spec"
)

func CleanSuggestion(raw string) string {
	s := strings.TrimSpace(raw)
	if strings.HasPrefix(s, "```") {
		lines := strings.Split(s, "\n")
		if len(lines) > 1 {
			endIdx := len(lines)
			if strings.HasPrefix(strings.TrimSpace(lines[len(lines)-1]), "```") {
				endIdx = len(lines) - 1
			}
			s = strings.TrimSpace(strings.Join(lines[1:endIdx], "\n"))
		}
	}
	if len(s) >= 2 && strings.HasPrefix(s, "`") && strings.HasSuffix(s, "`") && !strings.HasPrefix(s, "``") {
		s = s[1 : len(s)-1]
	}
	if len(s) >= 2 && ((strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")) || (strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'"))) {
		inner := s[1 : len(s)-1]
		if !strings.ContainsAny(inner, "\"'") {
			s = inner
		}
	}
	return strings.TrimSpace(s)
}

func NormalizeSuggestion(buf string, suggCmd string) string {
	rawCmd := strings.TrimSpace(suggCmd)
	suggCmd = CleanSuggestion(suggCmd)

	if strings.Contains(buf, "-m \"") || strings.Contains(buf, "-am \"") || strings.Contains(buf, "--message \"") {
		if !strings.HasPrefix(strings.ToLower(suggCmd), strings.ToLower(buf)) {
			for _, flag := range []string{"-m ", "-am ", "--message "} {
				idx := strings.Index(suggCmd, flag)
				if idx != -1 {
					afterFlag := suggCmd[idx+len(flag):]
					if !strings.HasPrefix(afterFlag, "\"") && !strings.HasPrefix(afterFlag, "'") {
						suggCmd = suggCmd[:idx+len(flag)] + "\"" + afterFlag + "\""
						break
					}
				}
			}
		}
	}

	if buf != "" {
		if strings.HasPrefix(strings.ToLower(suggCmd), strings.ToLower(buf)) && len(suggCmd) >= len(buf) {
			suggCmd = buf + suggCmd[len(buf):]
		} else if fields := strings.Fields(buf); len(fields) > 0 && len(suggCmd) > 0 {
			firstWord := strings.ToLower(fields[0])
			suggLow := strings.ToLower(suggCmd)
			if !strings.HasPrefix(suggLow, firstWord) && !strings.HasPrefix(suggLow, "sudo ") {
				delta := suggCmd
				if strings.HasPrefix(rawCmd, "-") || strings.HasPrefix(rawCmd, "\"") || strings.HasPrefix(rawCmd, "'") {
					delta = rawCmd
				}
				if strings.HasSuffix(buf, " ") || strings.HasSuffix(buf, "\"") || strings.HasSuffix(buf, "'") || strings.HasSuffix(buf, "=") || strings.HasSuffix(buf, "/") {
					suggCmd = buf + delta
				} else {
					suggCmd = buf + " " + delta
				}
			}
		}
	}

	return suggCmd
}

func ShouldOverwrite(originalBuf string, currentBuf string, newSugg *spec.Suggestion, currentConfidence int) bool {
	if newSugg == nil {
		return false
	}
	if !strings.HasPrefix(currentBuf, originalBuf) {
		return false
	}
	if !strings.HasPrefix(strings.ToLower(newSugg.Cmd), strings.ToLower(currentBuf)) {
		return false
	}
	return newSugg.Confidence > currentConfidence
}

func ExtractScriptsAndTargets(sb *strings.Builder, dir string, prefix string) {
	if data, err := os.ReadFile(filepath.Join(dir, "package.json")); err == nil {
		var pkg struct {
			Scripts map[string]string `json:"scripts"`
		}
		if err := json.Unmarshal(data, &pkg); err == nil && len(pkg.Scripts) > 0 {
			var scriptNames []string
			for name, cmd := range pkg.Scripts {
				scriptNames = append(scriptNames, fmt.Sprintf("%s: %s", name, cmd))
			}
			sort.Strings(scriptNames)
			if len(scriptNames) > 20 {
				scriptNames = append(scriptNames[:20], "... (truncated)")
			}
			label := "package.json"
			if prefix != "" {
				label = prefix + "/package.json"
			}
			fmt.Fprintf(sb, "Available %s scripts:\n%s\n\n", label, strings.Join(scriptNames, "\n"))
		}
	}
	if file, err := os.Open(filepath.Join(dir, "Makefile")); err == nil {
		defer func() { _ = file.Close() }()
		var targets []string
		seen := make(map[string]bool)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			// Skip variable assignments because operators like := or colons in values trick the colon parser into misclassifying variables as build targets
			if strings.Contains(line, "=") {
				continue
			}
			if idx := strings.Index(line, ":"); idx > 0 && !strings.HasPrefix(line, "\t") && !strings.HasPrefix(line, " ") {
				target := strings.TrimSpace(line[:idx])
				if target != "" && target != ".PHONY" && !strings.Contains(target, " ") && !seen[target] && !strings.HasPrefix(target, ".") {
					seen[target] = true
					targets = append(targets, target)
				}
			}
		}
		_ = scanner.Err()
		if len(targets) > 0 {
			sort.Strings(targets)
			// Cap at 10 targets to keep AI prompt short and avoid exceeding 6000 TPM limit (basedd on Groq api docs because Im using it now)
			if len(targets) > 10 {
				targets = append(targets[:10], "... (truncated)")
			}
			label := "Makefile"
			if prefix != "" {
				label = prefix + "/Makefile"
			}
			fmt.Fprintf(sb, "Available %s targets:\n%s\n\n", label, strings.Join(targets, ", "))
		}
	}
	openJustfile := func() (*os.File, error) {
		if f, err := os.Open(filepath.Join(dir, "justfile")); err == nil {
			return f, nil
		}
		return os.Open(filepath.Join(dir, "Justfile"))
	}
	if file, err := openJustfile(); err == nil {
		defer func() { _ = file.Close() }()
		var recipes []string
		seen := make(map[string]bool)
		scanner := bufio.NewScanner(file)
		recipeRegex := regexp.MustCompile(`^([a-zA-Z0-9_-]+):`)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "[") {
				continue
			}
			if matches := recipeRegex.FindStringSubmatch(line); len(matches) > 1 {
				recipe := matches[1]
				if !seen[recipe] {
					seen[recipe] = true
					recipes = append(recipes, recipe)
				}
			}
		}
		_ = scanner.Err()
		if len(recipes) > 0 {
			sort.Strings(recipes)
			// Cap at 10 recipes to keep AI prompt short and avoid exceeding 6000 TPM limit
			if len(recipes) > 10 {
				recipes = append(recipes[:10], "... (truncated)")
			}
			label := "justfile"
			if prefix != "" {
				label = prefix + "/justfile"
			}
			fmt.Fprintf(sb, "Available %s recipes:\n%s\n\n", label, strings.Join(recipes, ", "))
		}
	}
}
