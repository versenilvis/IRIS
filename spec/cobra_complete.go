package spec

import (
	"context"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type cobraCacheEntry struct {
	suggestions []Suggestion
}

var (
	cobraCache   = map[string]cobraCacheEntry{}
	cobraCacheMu sync.Mutex
)

func cobraBinKey(binName string) string {
	path, err := exec.LookPath(binName)
	if err != nil {
		return binName
	}
	info, err := os.Stat(path)
	if err != nil {
		return binName
	}
	return binName + "|" + info.ModTime().String()
}

// parseCobraOutput parses output from `<cmd> __complete <args>`.
// each line is "value\tdesc", last line is ":N" (ShellCompDirective bitmask).
// returns nil if output is not Cobra-style.
func parseCobraOutput(raw string, prefix string) []Suggestion {
	lines := strings.Split(strings.TrimRight(raw, "\n"), "\n")
	if len(lines) == 0 {
		return nil
	}
	lastLine := lines[len(lines)-1]
	if !strings.HasPrefix(lastLine, ":") {
		return nil
	}
	directive, err := strconv.Atoi(lastLine[1:])
	if err != nil {
		return nil
	}
	// ShellCompDirectiveError = 1
	if directive&1 != 0 {
		return nil
	}

	candidates := lines[:len(lines)-1]
	results := make([]Suggestion, 0, len(candidates))
	for _, line := range candidates {
		if line == "" {
			continue
		}
		value, desc, _ := strings.Cut(line, "\t")
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		cmd := value
		if prefix != "" {
			cmd = prefix + " " + value
		}
		results = append(results, Suggestion{
			Cmd:        cmd,
			Desc:       desc,
			Source:     "spec-inferred",
			Confidence: 50,
			Priority:   30,
		})
	}
	return results
}

func buildCobraCacheKey(binKey string, args []string, partial string) string {
	var sb strings.Builder
	sb.WriteString(binKey)
	for _, arg := range args {
		sb.WriteByte('\x00')
		sb.WriteString(arg)
	}
	sb.WriteByte('\x00')
	sb.WriteString(partial)
	return sb.String()
}

// newProbeCmd builds and isolates the `__complete` probe command.
// starts the child in its own session so it has no controlling terminal
// and therefore won't affect the user's tty in the case of programs that
// don't respect `__complete`.
func newProbeCmd(ctx context.Context, binName string, args []string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, binName, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	return cmd
}

// QueryCobraComplete calls `binName __complete <args> <partial>` and returns
// structured suggestions cached per binary mtime, args, and partial.
// returns nil if the binary is not Cobra-based or times out.
func QueryCobraComplete(binName string, args []string, partial string) []Suggestion {
	if strings.ContainsAny(binName, `/\`) {
		return nil
	}

	binKey := cobraBinKey(binName)
	argKey := buildCobraCacheKey(binKey, args, partial)

	cobraCacheMu.Lock()
	if entry, ok := cobraCache[argKey]; ok {
		cobraCacheMu.Unlock()
		return filterByPartial(entry.suggestions, partial)
	}
	cobraCacheMu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	cmdArgs := append([]string{"__complete"}, args...)
	cmdArgs = append(cmdArgs, partial)
	probe := newProbeCmd(ctx, binName, cmdArgs)
	out, err := probe.Output()
	if err != nil {
		return nil
	}

	prefixParts := append([]string{binName}, args...)
	prefix := strings.Join(prefixParts, " ")
	suggestions := parseCobraOutput(string(out), prefix)

	cobraCacheMu.Lock()
	cobraCache[argKey] = cobraCacheEntry{suggestions: suggestions}
	cobraCacheMu.Unlock()

	return filterByPartial(suggestions, partial)
}

func filterByPartial(suggestions []Suggestion, partial string) []Suggestion {
	if partial == "" {
		return suggestions
	}
	filtered := make([]Suggestion, 0, len(suggestions))
	for _, s := range suggestions {
		lastWord := s.Cmd
		if idx := strings.LastIndex(s.Cmd, " "); idx >= 0 {
			lastWord = s.Cmd[idx+1:]
		}
		if HasPrefix(lastWord, partial) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

// ResetCobraCache clears the completion cache — use in tests only
func ResetCobraCache() {
	cobraCacheMu.Lock()
	cobraCache = map[string]cobraCacheEntry{}
	cobraCacheMu.Unlock()
}
