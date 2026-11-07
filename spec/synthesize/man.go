package synthesize

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/versenilvis/iris/spec"
)

var (
	reFlagLong     = regexp.MustCompile(`--[a-zA-Z0-9][a-zA-Z0-9_.+-]*`)
	reFlagShort    = regexp.MustCompile(`^-[a-zA-Z0-9]$`)
	reValidFlag    = regexp.MustCompile(`^-{1,2}[A-Za-z0-9][A-Za-z0-9_.+-]*$`)
	reANSI         = regexp.MustCompile(`\x1b(?:\[[0-9;?]*[a-zA-Z]|\].*?(?:\x07|\x1b\\)|[PX^_].*?\x1b\\|[()][AB012])`)
	reRoffFont     = regexp.MustCompile(`\\f(?:\[[^\]]*\]|\([A-Za-z0-9]{2}|[A-Za-z0-9()])`)
	reRoffDevice   = regexp.MustCompile(`\\X'[^']*'`)
	reRoffSize     = regexp.MustCompile(`\\s[+-]?[0-9]+`)
	reRoffSpecial  = regexp.MustCompile(`\\[^a-zA-Z0-9]`)
	reTopSection   = regexp.MustCompile(`^\.(?:SH|Sh)\s+(?:"([^"]+)"|([^\s]+.*))`)
	reSubSection   = regexp.MustCompile(`^\.(?:SS|Ss)\s+(?:"([^"]+)"|([^\s]+.*))`)
	reMdocItemFlag = regexp.MustCompile(`^\.It\s+(?:Fl\s+)+`)
)

const (
	maxLineLength = 64 * 1024
	maxScanLines  = 10000
)

type countedReader struct {
	r     io.Reader
	limit int64
	read  int64
}

func (c *countedReader) Read(p []byte) (int, error) {
	if c.limit > 0 && c.read >= c.limit {
		return 0, fmt.Errorf("stream exceeds maximum size limit of %d bytes", c.limit)
	}
	maxToRead := len(p)
	if c.limit > 0 && int64(maxToRead) > (c.limit-c.read) {
		maxToRead = int(c.limit - c.read + 1)
	}
	n, err := c.r.Read(p[:maxToRead])
	c.read += int64(n)
	if c.limit > 0 && c.read > c.limit {
		return n, fmt.Errorf("stream exceeds maximum size limit of %d bytes", c.limit)
	}
	return n, err
}

func sanitizeText(s string) string {
	s = reANSI.ReplaceAllString(s, "")
	var b strings.Builder
	for _, r := range s {
		// drop C0, C1, DEL, and ESC control characters
		if r < 0x20 || (r >= 0x7f && r <= 0x9f) || r == 0x1b {
			if r == '\t' || r == '\n' || r == '\r' {
				b.WriteByte(' ')
			}
			continue
		}
		b.WriteRune(r)
	}
	return strings.TrimSpace(b.String())
}

func truncateRunes(s string, maxRunes int) string {
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	return string(runes[:maxRunes])
}

func cleanRoffText(s string) string {
	s = reRoffDevice.ReplaceAllString(s, "")
	s = reRoffFont.ReplaceAllString(s, "")
	s = reRoffSize.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, `\-`, "-")
	s = strings.ReplaceAll(s, `\(hy`, "-")
	s = strings.ReplaceAll(s, `\(en`, "-")
	s = strings.ReplaceAll(s, `\(em`, "-")
	s = strings.ReplaceAll(s, `\(aq`, "'")
	s = strings.ReplaceAll(s, `\(dq`, `"`)
	s = strings.ReplaceAll(s, `\e`, `\`)
	s = strings.ReplaceAll(s, `\ `, " ")
	s = strings.ReplaceAll(s, `\^`, "")
	s = strings.ReplaceAll(s, `\&`, "")
	s = strings.ReplaceAll(s, `\|`, "")
	s = strings.ReplaceAll(s, `\c`, "")
	s = reRoffSpecial.ReplaceAllString(s, "")
	return sanitizeText(s)
}

func expandRoffTag(line string) string {
	trimmed := strings.TrimSpace(line)
	if after, ok := strings.CutPrefix(trimmed, ".It"); ok {
		// mdoc option macro conversion
		rest := after
		rest = strings.ReplaceAll(rest, "Fl Fl ", "--")
		rest = strings.ReplaceAll(rest, "Fl Fl", "--")
		rest = strings.ReplaceAll(rest, "Fl ", "-")
		rest = strings.ReplaceAll(rest, "Fl", "-")
		return cleanRoffText(rest)
	}

	for _, macro := range []string{".BR", ".B", ".BI", ".IR", ".IB", ".I", ".RB", ".RI"} {
		if after, ok := strings.CutPrefix(trimmed, macro+" "); ok {
			args := parseMacroArgs(after)
			return cleanRoffText(strings.Join(args, " "))
		}
	}

	return cleanRoffText(trimmed)
}

func parseMacroArgs(s string) []string {
	var args []string
	var current strings.Builder
	inQuote := false

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '"' {
			inQuote = !inQuote
			continue
		}
		if ch == ' ' && !inQuote {
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
			continue
		}
		current.WriteByte(ch)
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}

func findShortFlags(s string) []string {
	var out []string
	for i := 0; i+1 < len(s); i++ {
		if s[i] != '-' {
			continue
		}
		ch := s[i+1]
		isAlphaNum := (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')
		if !isAlphaNum {
			continue
		}
		if i > 0 {
			prev := s[i-1]
			if prev == '-' || (prev >= 'a' && prev <= 'z') || (prev >= 'A' && prev <= 'Z') || (prev >= '0' && prev <= '9') || prev == '_' || prev == '.' {
				continue
			}
		}
		if i+2 < len(s) {
			next := s[i+2]
			if next == '-' || (next >= 'a' && next <= 'z') || (next >= 'A' && next <= 'Z') || (next >= '0' && next <= '9') || next == '_' {
				continue
			}
		}
		out = append(out, s[i:i+2])
		i++
	}
	return out
}

func extractFlags(line string) (string, []string) {
	longMatches := reFlagLong.FindAllString(line, -1)
	remaining := reFlagLong.ReplaceAllString(line, " ")
	shortMatches := findShortFlags(remaining)

	var uniqueLong []string
	seen := make(map[string]bool)
	for _, f := range longMatches {
		if reValidFlag.MatchString(f) && !seen[f] {
			seen[f] = true
			uniqueLong = append(uniqueLong, f)
		}
	}

	var uniqueShort []string
	for _, f := range shortMatches {
		if reValidFlag.MatchString(f) && !seen[f] {
			seen[f] = true
			uniqueShort = append(uniqueShort, f)
		}
	}

	if len(uniqueLong) == 0 && len(uniqueShort) == 0 {
		return "", nil
	}

	if len(uniqueLong) > 0 {
		name := uniqueLong[0]
		aliases := append(uniqueLong[1:], uniqueShort...)
		return name, aliases
	}

	name := uniqueShort[0]
	return name, uniqueShort[1:]
}

func cleanDescription(lines []string) string {
	var cleanParts []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, `.\"`) || strings.HasPrefix(l, `'\"`) {
			continue
		}
		if strings.HasPrefix(l, ".") {
			if strings.HasPrefix(l, ".B ") || strings.HasPrefix(l, ".I ") || strings.HasPrefix(l, ".BR ") {
				parts := parseMacroArgs(l[strings.Index(l, " ")+1:])
				cleanParts = append(cleanParts, cleanRoffText(strings.Join(parts, " ")))
			}
			continue
		}
		clean := cleanRoffText(l)
		if clean != "" {
			cleanParts = append(cleanParts, clean)
		}
	}

	full := strings.Join(cleanParts, " ")
	full = strings.Join(strings.Fields(full), " ")
	if full == "" {
		return ""
	}

	if dotIdx := strings.Index(full, ". "); dotIdx != -1 {
		full = full[:dotIdx]
	}
	// strip trailing punctuation for clean suggestion display
	full = strings.TrimRight(full, ". \t\r\n,;:")
	full = truncateRunes(full, 160)
	full = strings.TrimRight(full, ". \t\r\n,;:")
	return sanitizeText(full)
}

func parseManReader(r io.Reader, cmdName string, maxSize int64) (string, []spec.UserOptionJSON, error) {
	if maxSize <= 0 {
		maxSize = defaultMaxFileSize
	}

	buf := bufio.NewReader(r)
	// detect gzip compression without external dependencies
	magic, err := buf.Peek(2)
	var stream io.Reader = buf
	if err == nil && len(magic) == 2 && magic[0] == 0x1f && magic[1] == 0x8b {
		gz, gzErr := gzip.NewReader(buf)
		if gzErr != nil {
			return "", nil, gzErr
		}
		defer gz.Close()
		stream = gz
	}

	// guard decompressed streams against unbounded memory usage
	stream = &countedReader{r: stream, limit: maxSize}

	scanner := bufio.NewScanner(stream)
	scanner.Buffer(make([]byte, 4096), maxLineLength)

	var rawLines []string
	for scanner.Scan() {
		if len(rawLines) >= maxScanLines {
			return "", nil, fmt.Errorf("manpage line count exceeds limit of %d", maxScanLines)
		}
		rawLines = append(rawLines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", nil, err
	}

	var cmdDesc string
	var currentSection string
	var currentSubSection string
	var inNameSection bool
	var options []spec.UserOptionJSON
	seenNames := make(map[string]bool)

	type pendingOption struct {
		name    string
		aliases []string
		lines   []string
	}
	var current *pendingOption

	commitCurrent := func() {
		if current == nil {
			return
		}
		name := sanitizeText(current.name)
		if !reValidFlag.MatchString(name) {
			current = nil
			return
		}
		var cleanAliases []string
		for _, a := range current.aliases {
			ca := sanitizeText(a)
			if reValidFlag.MatchString(ca) && ca != name {
				cleanAliases = append(cleanAliases, ca)
			}
		}

		if !seenNames[name] {
			seenNames[name] = true
			desc := cleanDescription(current.lines)
			options = append(options, spec.UserOptionJSON{
				Name:        name,
				Aliases:     cleanAliases,
				Description: desc,
			})
		}
		current = nil
	}

	for i := 0; i < len(rawLines); i++ {
		line := rawLines[i]
		trimmed := strings.TrimSpace(line)

		if sm := reTopSection.FindStringSubmatch(trimmed); sm != nil {
			commitCurrent()
			sec := sm[1]
			if sec == "" {
				sec = sm[2]
			}
			currentSection = strings.ToUpper(strings.TrimSpace(sec))
			currentSubSection = ""
			inNameSection = currentSection == "NAME"
			continue
		}

		if sm := reSubSection.FindStringSubmatch(trimmed); sm != nil {
			commitCurrent()
			sub := sm[1]
			if sub == "" {
				sub = sm[2]
			}
			currentSubSection = strings.ToUpper(strings.TrimSpace(sub))
			continue
		}

		if inNameSection && cmdDesc == "" && !strings.HasPrefix(trimmed, ".") {
			cleaned := cleanRoffText(trimmed)
			if idx := strings.LastIndex(cleaned, " - "); idx != -1 {
				cmdDesc = strings.TrimSpace(cleaned[idx+3:])
			} else if _, after, ok := strings.Cut(cleaned, "- "); ok {
				cmdDesc = strings.TrimSpace(after)
			}
			continue
		}

		if inNameSection && strings.HasPrefix(trimmed, ".Nd ") {
			cmdDesc = cleanRoffText(strings.TrimPrefix(trimmed, ".Nd "))
			continue
		}

		isOptionSection := strings.Contains(currentSection, "OPTION") ||
			strings.Contains(currentSection, "DESCRIPTION") ||
			strings.Contains(currentSection, "FLAGS") ||
			strings.Contains(currentSubSection, "OPTION") ||
			strings.Contains(currentSubSection, "FLAGS")

		if !isOptionSection {
			continue
		}

		if strings.HasPrefix(trimmed, ".TP") || strings.HasPrefix(trimmed, ".HP") {
			commitCurrent()
			for i+1 < len(rawLines) {
				nextLine := strings.TrimSpace(rawLines[i+1])
				if nextLine == "" || nextLine == "." || strings.HasPrefix(nextLine, `.\"`) || strings.HasPrefix(nextLine, `'\"`) {
					i++
					continue
				}
				break
			}
			if i+1 < len(rawLines) {
				i++
				tagLine := expandRoffTag(rawLines[i])
				if name, aliases := extractFlags(tagLine); name != "" {
					current = &pendingOption{name: name, aliases: aliases}
					continue
				}
			}
			continue
		}

		if strings.HasPrefix(trimmed, ".IP") {
			commitCurrent()
			tagLine := expandRoffTag(trimmed)
			if name, aliases := extractFlags(tagLine); name != "" {
				current = &pendingOption{name: name, aliases: aliases}
				continue
			}
		}

		if reMdocItemFlag.MatchString(trimmed) {
			commitCurrent()
			tagLine := expandRoffTag(trimmed)
			if name, aliases := extractFlags(tagLine); name != "" {
				current = &pendingOption{name: name, aliases: aliases}
				continue
			}
		}

		// check for standalone lines formatted as option headers
		if strings.HasPrefix(trimmed, `\fB\-`) || strings.HasPrefix(trimmed, `.B \-`) || strings.HasPrefix(trimmed, `.BR \-`) {
			commitCurrent()
			tagLine := expandRoffTag(trimmed)
			if name, aliases := extractFlags(tagLine); name != "" {
				current = &pendingOption{name: name, aliases: aliases}
				continue
			}
		}

		if current != nil {
			current.lines = append(current.lines, trimmed)
		}
	}

	commitCurrent()
	cmdDesc = sanitizeText(cmdDesc)
	cmdDesc = truncateRunes(cmdDesc, 160)
	cmdDesc = strings.TrimRight(cmdDesc, ". \t\r\n")

	return cmdDesc, options, nil
}

func parseManFile(dir, cmdName string, maxSize int64) (string, []spec.UserOptionJSON, error) {
	cleanName := filepath.Base(cmdName)
	candidates := []string{
		cleanName + ".1",
		cleanName + ".1.gz",
		cleanName + ".8",
		cleanName + ".8.gz",
	}

	var lastErr error
	for _, cand := range candidates {
		reader, err := safeOpenFile(dir, cand, maxSize)
		if err != nil {
			lastErr = err
			continue
		}
		defer reader.Close()

		desc, opts, err := parseManReader(reader, cleanName, maxSize)
		if err != nil {
			return "", nil, err
		}
		if len(opts) > 0 || desc != "" {
			return desc, opts, nil
		}
	}

	if lastErr != nil {
		return "", nil, lastErr
	}
	return "", nil, nil
}
