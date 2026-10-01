package ctxcheck

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type MakeManifest struct {
	Targets        map[string]bool
	HasInclude     bool
	HasPatternRule bool
}

func findMakefile(cwd string) (string, error) {
	candidates := []string{"GNUmakefile", "makefile", "Makefile"}
	for _, name := range candidates {
		p := filepath.Join(cwd, name)
		if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
			return p, nil
		}
	}
	return "", os.ErrNotExist
}

func parseMakefile(path string) (*MakeManifest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	manifest := &MakeManifest{
		Targets: make(map[string]bool),
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		rawLine := scanner.Text()
		if strings.HasPrefix(rawLine, "\t") || strings.HasPrefix(rawLine, " ") {
			continue
		}
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "include ") || strings.HasPrefix(line, "-include ") || strings.HasPrefix(line, "sinclude ") {
			manifest.HasInclude = true
			continue
		}

		colonIdx := strings.IndexByte(line, ':')
		if colonIdx <= 0 {
			continue
		}
		if colonIdx+1 < len(line) && (line[colonIdx+1] == '=' || line[colonIdx+1] == ':') {
			// variable assignment: := or ::=
			continue
		}

		targetPart := strings.TrimSpace(line[:colonIdx])
		if strings.Contains(targetPart, "%") {
			manifest.HasPatternRule = true
			continue
		}

		for target := range strings.FieldsSeq(targetPart) {
			if isValidMakeTarget(target) {
				manifest.Targets[target] = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return manifest, nil
}

func isValidMakeTarget(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' && r != '.' {
			return false
		}
	}
	return true
}

func ValidateMake(tokens []string, cwd string) Verdict {
	if len(tokens) == 0 || tokens[0] != "make" {
		return Free
	}

	path, err := getCachedPath("make", cwd, findMakefile)
	if err != nil {
		return Invalid
	}

	args := tokens[1:]
	var target string

	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "-C" || arg == "-f" || arg == "--file" || arg == "--makefile" || arg == "--directory" ||
			strings.HasPrefix(arg, "--file=") || strings.HasPrefix(arg, "--makefile=") || strings.HasPrefix(arg, "--directory=") ||
			(strings.HasPrefix(arg, "-C") && len(arg) > 2) {
			return Unknown
		}
		if arg == "-j" || arg == "-l" || arg == "-o" || arg == "-W" || arg == "-I" {
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				i += 2
				continue
			}
			i++
			continue
		}
		if strings.HasPrefix(arg, "-") {
			// ignore standard flags
			i++
			continue
		}
		if strings.Contains(arg, "=") {
			// VAR=val
			i++
			continue
		}
		target = arg
		break
	}

	// bare make runs first target
	if target == "" {
		return Valid
	}

	data, err := getCachedManifest(path, func(p string) (any, error) {
		return parseMakefile(p)
	})
	if err != nil {
		return Unknown
	}
	manifest, ok := data.(*MakeManifest)
	if !ok || manifest == nil {
		return Unknown
	}

	if manifest.Targets[target] {
		return Valid
	}

	if manifest.HasInclude || manifest.HasPatternRule {
		return Unknown
	}

	return Invalid
}
