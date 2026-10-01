package ctxcheck

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type JustManifest struct {
	Recipes        map[string]bool
	Mods           map[string]bool
	HasFallback    bool
	HasImportOrMod bool
}

func findJustfile(cwd string) (string, error) {
	dir := filepath.Clean(cwd)
	candidates := []string{"justfile", "Justfile", ".justfile"}
	for {
		for _, name := range candidates {
			p := filepath.Join(dir, name)
			if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
				return p, nil
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", os.ErrNotExist
}

func parseJustfile(path string) (*JustManifest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	manifest := &JustManifest{
		Recipes: make(map[string]bool),
		Mods:    make(map[string]bool),
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		rawLine := scanner.Text()
		if strings.HasPrefix(rawLine, " ") || strings.HasPrefix(rawLine, "\t") {
			// indented line is recipe body
			continue
		}
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "[") {
			continue
		}

		if strings.HasPrefix(line, "set ") {
			if strings.Contains(line, "fallback") {
				manifest.HasFallback = true
			}
			continue
		}
		if strings.HasPrefix(line, "export ") {
			continue
		}

		if strings.HasPrefix(line, "import ") || strings.HasPrefix(line, "import?") {
			manifest.HasImportOrMod = true
			continue
		}

		if strings.HasPrefix(line, "mod ") || strings.HasPrefix(line, "mod?") {
			manifest.HasImportOrMod = true
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				modName := strings.Trim(fields[1], "\"'")
				if isValidRecipeName(modName) {
					manifest.Mods[modName] = true
				}
			}
			continue
		}

		if strings.HasPrefix(line, "alias ") {
			fields := strings.Fields(line)
			if len(fields) >= 4 && fields[0] == "alias" && fields[2] == ":=" {
				aliasName := fields[1]
				if isValidRecipeName(aliasName) {
					manifest.Recipes[aliasName] = true
				}
			}
			continue
		}

		// check for recipe header
		colonIdx := strings.IndexByte(line, ':')
		if colonIdx <= 0 {
			continue
		}
		if colonIdx+1 < len(line) && line[colonIdx+1] == '=' {
			// variable assignment: name := val
			continue
		}

		headerPart := strings.TrimSpace(line[:colonIdx])
		headerFields := strings.Fields(headerPart)
		if len(headerFields) > 0 {
			name := strings.TrimPrefix(headerFields[0], "@")
			if isValidRecipeName(name) {
				manifest.Recipes[name] = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return manifest, nil
}

func isValidRecipeName(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return false
		}
	}
	return true
}

func ValidateJust(tokens []string, cwd string) Verdict {
	if len(tokens) == 0 || tokens[0] != "just" {
		return Free
	}

	path, err := getCachedPath("just", cwd, findJustfile)
	if err != nil {
		return Invalid
	}

	args := tokens[1:]
	var recipeName string

	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "-f" || arg == "--justfile" || arg == "-d" || arg == "--working-directory" {
			return Unknown
		}
		if strings.Contains(arg, "::") || strings.Contains(arg, "/") {
			return Unknown
		}
		if strings.Contains(arg, "=") {
			return Unknown
		}

		if strings.HasPrefix(arg, "-") {
			// flags that take arguments
			if arg == "--color" || arg == "--command-color" || arg == "--chooser" ||
				arg == "--dotenv-filename" || arg == "--dotenv-path" || arg == "--list-heading" ||
				arg == "--list-prefix" || arg == "--shell" || arg == "--shell-arg" || arg == "--summary" {
				i += 2
				continue
			}
			i++
			continue
		}

		recipeName = arg
		break
	}

	// bare `just` runs default recipe
	if recipeName == "" {
		return Valid
	}

	data, err := getCachedManifest(path, func(p string) (any, error) {
		return parseJustfile(p)
	})
	if err != nil {
		return Unknown
	}
	manifest, ok := data.(*JustManifest)
	if !ok || manifest == nil {
		return Unknown
	}

	if manifest.Mods[recipeName] {
		return Unknown
	}

	if manifest.Recipes[recipeName] {
		return Valid
	}

	if manifest.HasFallback || manifest.HasImportOrMod {
		return Unknown
	}

	return Invalid
}
