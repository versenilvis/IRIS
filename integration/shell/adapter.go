package shell

import (
	"context"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"
)

// Adapter defines the behavior for different shell environments
type Adapter interface {
	GetName() string
	GetShellPath() string
	GetEnv(fd int, pid int) []string
	PrepareSelectSequence(selected string) []byte
	// ScanAliases returns a map of alias name to target command
	ScanAliases() map[string]string
}

// Current shell instance
var Current Adapter

func Init(name string) {
	switch name {
	case "zsh":
		Current = &ZshAdapter{}
	case "fish":
		Current = &FishAdapter{}
	default:
		Current = &BashAdapter{}
	}
}

// BashAdapter implementation
type BashAdapter struct{}

func (b *BashAdapter) GetName() string      { return "bash" }
func (b *BashAdapter) GetShellPath() string { return "bash" }
func (b *BashAdapter) GetEnv(fd int, pid int) []string {
	return append(os.Environ(), "IRIS_FD="+fmt.Sprint(fd), "IRIS_PID="+fmt.Sprint(pid))
}
func (b *BashAdapter) PrepareSelectSequence(selected string) []byte {
	return append([]byte{0x15}, []byte(selected)...)
}
func (b *BashAdapter) ScanAliases() map[string]string {
	return ScanPosixAliases([]string{".bashrc", ".bash_profile", ".bash_aliases"})
}

// ZshAdapter implementation
type ZshAdapter struct{}

func (z *ZshAdapter) GetName() string      { return "zsh" }
func (z *ZshAdapter) GetShellPath() string { return "zsh" }
func (z *ZshAdapter) GetEnv(fd int, pid int) []string {
	return append(os.Environ(), "IRIS_FD="+fmt.Sprint(fd), "IRIS_PID="+fmt.Sprint(pid))
}
func (z *ZshAdapter) PrepareSelectSequence(selected string) []byte {
	return append([]byte{0x15}, []byte(selected)...)
}
func (z *ZshAdapter) ScanAliases() map[string]string {
	envSet := os.Getenv("ZDOTDIR") != ""
	zdotdir := GetZshConfigDir()
	home, _ := os.UserHomeDir()
	
	var files []string
	if !envSet && zdotdir != home {
		files = append(files, filepath.Join(home, ".zshenv"))
	}
	
	files = append(files,
		filepath.Join(zdotdir, ".zshenv"),
		filepath.Join(zdotdir, ".zprofile"),
		filepath.Join(zdotdir, ".zshrc"),
	)
	
	return ScanPosixAliases(files)
}

func GetZshConfigDir() string {
	if zdotdir := os.Getenv("ZDOTDIR"); zdotdir != "" {
		return zdotdir
	}

	// Fallback: ask zsh directly in case ZDOTDIR is set in ~/.zshenv without
	// export. ScanAliases runs on every keystroke, so this is resolved once per
	// process -- ZDOTDIR can't change under a running iris, and paying for a
	// zsh subprocess per keystroke costs several milliseconds of input latency.
	zdotdirOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
		cmd := exec.CommandContext(ctx, "zsh", "-c", "echo $ZDOTDIR")
		out, err := cmd.Output()
		if err == nil {
			zdotdirCached = strings.TrimSpace(string(out))
		}
		if zdotdirCached == "" {
			zdotdirCached, _ = os.UserHomeDir()
		}
	})

	return zdotdirCached
}

var (
	zdotdirOnce   sync.Once
	zdotdirCached string
)

// FishAdapter implementation
type FishAdapter struct{}

func (f *FishAdapter) GetName() string      { return "fish" }
func (f *FishAdapter) GetShellPath() string { return "fish" }
func (f *FishAdapter) GetEnv(fd int, pid int) []string {
	return append(os.Environ(), "IRIS_FD="+fmt.Sprint(fd), "IRIS_PID="+fmt.Sprint(pid))
}
func (f *FishAdapter) PrepareSelectSequence(selected string) []byte {
	return append([]byte{0x15}, []byte(selected)...)
}
func (f *FishAdapter) ScanAliases() map[string]string {
	// fish uses 'alias' command in config.fish or separate function files
	return ScanPosixAliases([]string{filepath.Join(GetFishConfigDir(), "config.fish")})
}

func GetFishConfigDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "fish")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "fish")
}

func GetFishDataDir() string {
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return filepath.Join(xdg, "fish")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "fish")
}

// maxSourceDepth caps how deep ScanPosixAliases will follow a chain of
// `source` lines. Layered dotfiles rarely nest more than two or three levels;
// the cap is a backstop for pathological configs, since cycles are already
// handled by the visited set.
const maxSourceDepth = 8

// aliasScanCache memoizes the last scan. Lookup re-scans aliases on every
// keystroke, and following `source` chains turns that into a walk over every
// config file in the tree -- work that is wasted while the files sit unchanged.
// A stat of each file already visited is enough to tell whether the cached
// result is still good, and is roughly two orders of magnitude cheaper.
var aliasScanCache struct {
	sync.Mutex
	entryPoints []string
	stamps      map[string]fileStamp
	aliases     map[string]string
}

type fileStamp struct {
	modTime time.Time
	size    int64
}

func statStamp(path string) (fileStamp, bool) {
	info, err := os.Stat(path)
	if err != nil {
		return fileStamp{}, false
	}
	return fileStamp{modTime: info.ModTime(), size: info.Size()}, true
}

func ScanPosixAliases(files []string) map[string]string {
	home, err := os.UserHomeDir()
	if err != nil {
		return map[string]string{}
	}

	paths := make([]string, 0, len(files))
	for _, f := range files {
		if filepath.IsAbs(f) {
			paths = append(paths, f)
			continue
		}
		paths = append(paths, filepath.Join(home, f))
	}

	aliasScanCache.Lock()
	defer aliasScanCache.Unlock()

	if cached := cachedAliases(paths); cached != nil {
		return cached
	}

	aliases := make(map[string]string)
	stamps := make(map[string]fileStamp)
	visited := make(map[string]bool)
	for _, path := range paths {
		scanAliasFile(path, aliases, visited, stamps, 0)
	}

	aliasScanCache.entryPoints = paths
	aliasScanCache.stamps = stamps
	aliasScanCache.aliases = aliases

	return maps.Clone(aliases)
}

// cachedAliases returns the previous scan when the same entry points were asked
// for and every file it touched is byte-for-byte where it was. A file that has
// since gone missing counts as a change, so a deleted include invalidates too.
func cachedAliases(paths []string) map[string]string {
	if aliasScanCache.aliases == nil || !slices.Equal(aliasScanCache.entryPoints, paths) {
		return nil
	}
	for path, want := range aliasScanCache.stamps {
		got, ok := statStamp(path)
		if want == (fileStamp{}) {
			// Was missing at scan time; it appearing is a change.
			if ok {
				return nil
			}
			continue
		}
		if !ok || got != want {
			return nil
		}
	}
	return maps.Clone(aliasScanCache.aliases)
}

// scanAliasFile records the aliases defined in path, following any file it
// sources. Most real configs keep aliases in a dedicated file pulled in with
// `source`, so reading only the entry points would miss nearly all of them.
func scanAliasFile(path string, aliases map[string]string, visited map[string]bool, stamps map[string]fileStamp, depth int) {
	if depth > maxSourceDepth {
		return
	}

	// Key the visited set on the resolved path so a dotfile symlinked into
	// $HOME isn't scanned twice under both names.
	key := path
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		key = resolved
	}
	if visited[key] {
		return
	}
	visited[key] = true

	// Stamp missing files too, as the zero value: a config that sources a file
	// which does not exist yet must re-scan once that file appears, and the
	// parent's own mtime won't change when it does.
	stamp, _ := statStamp(path)
	stamps[path] = stamp

	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	walkShellConfig(string(data), aliases, func(target string) {
		scanAliasFile(target, aliases, visited, stamps, depth+1)
	})
}

func ParseAliases(data string) map[string]string {
	aliases := make(map[string]string)
	walkShellConfig(data, aliases, nil)
	return aliases
}

// walkShellConfig reads shell config text top to bottom, recording alias
// definitions into aliases and handing each `source`/`.` target to onSource (if
// non-nil) at the point it appears, so a sourced file's definitions land in the
// same order the shell would apply them.
func walkShellConfig(data string, aliases map[string]string, onSource func(target string)) {
	// A config that defines the same alias in both arms of a conditional --
	// `if which eza; then alias ls='eza'; else alias ls='lsd'; fi` -- would
	// otherwise end up with the fallback, purely because it is written last.
	// The condition can't be evaluated from static text, so prefer the first
	// arm: it is conventionally the preferred tool, with later arms as
	// fallbacks. altArm tracks whether any enclosing `if` has moved past it.
	altArm := 0
	depth := 0

	lines := strings.SplitSeq(data, "\n")
	for line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		switch {
		case line == "if" || strings.HasPrefix(line, "if "):
			depth++
		case depth > 0 && (line == "else" || line == "elif" || strings.HasPrefix(line, "elif ")):
			altArm++
		case depth > 0 && (line == "fi" || strings.HasPrefix(line, "fi ") || strings.HasPrefix(line, "fi;")):
			depth--
			if altArm > 0 {
				altArm--
			}
		}

		if onSource != nil {
			for _, target := range sourceTargets(line) {
				onSource(target)
			}
		}

		if altArm > 0 {
			continue
		}
		parseAliasLine(line, aliases)
	}
}

func parseAliasLine(line string, aliases map[string]string) {
	if !strings.HasPrefix(line, "alias ") {
		return
	}
	body := strings.TrimSpace(strings.TrimPrefix(line, "alias"))
	if body == "" {
		return
	}

	pairs := SplitAliasTokens(body)
	for _, pair := range pairs {
		before, after, ok := strings.Cut(pair, "=")
		if !ok {
			continue
		}
		key := strings.TrimSpace(before)
		val := unquoteAliasValue(strings.TrimSpace(after))
		if key != "" && val != "" {
			aliases[key] = val
		}
	}
}

// sourceTargets returns the absolute paths a single config line sources. The
// line is split on shell separators first so guarded forms such as
// `[ -f "$X" ] && source "$X"` are picked up, not just lines that begin with
// `source`. Targets that can't be resolved to a literal path -- unset
// variables, command or process substitution, globs -- are dropped rather than
// guessed at.
func sourceTargets(line string) []string {
	var targets []string
	for _, segment := range splitShellSegments(line) {
		segment = strings.TrimSpace(segment)

		var rest string
		switch {
		case strings.HasPrefix(segment, "source "):
			rest = segment[len("source "):]
		case strings.HasPrefix(segment, ". "):
			rest = segment[len(". "):]
		default:
			continue
		}

		fields := SplitAliasTokens(strings.TrimSpace(rest))
		if len(fields) == 0 {
			continue
		}
		if path, ok := expandConfigPath(fields[0]); ok {
			targets = append(targets, path)
		}
	}
	return targets
}

// splitShellSegments breaks a line on `&&`, `||` and `;` outside of quotes, so
// each segment can be tested for a leading command word.
func splitShellSegments(line string) []string {
	var segments []string
	var cur strings.Builder
	inQuote := false
	var quoteChar rune

	runes := []rune(line)
	for i := 0; i < len(runes); i++ {
		c := runes[i]
		switch {
		case !inQuote && (c == '"' || c == '\''):
			inQuote = true
			quoteChar = c
			cur.WriteRune(c)
		case inQuote && c == quoteChar:
			inQuote = false
			cur.WriteRune(c)
		case !inQuote && c == ';':
			segments = append(segments, cur.String())
			cur.Reset()
		case !inQuote && i+1 < len(runes) && (c == '&' || c == '|') && runes[i+1] == c:
			segments = append(segments, cur.String())
			cur.Reset()
			i++
		default:
			cur.WriteRune(c)
		}
	}
	segments = append(segments, cur.String())
	return segments
}

// expandConfigPath turns a quoted, variable-bearing path from a config file
// into an absolute path. It reports false when the path can't be resolved
// statically, which is the common case for loop variables (`source "$file"`)
// and process substitution (`source <(fzf --zsh)`).
func expandConfigPath(raw string) (string, bool) {
	path := strings.Trim(strings.TrimSpace(raw), `"'`)
	if path == "" {
		return "", false
	}
	if strings.ContainsAny(path, "*?`") || strings.Contains(path, "$(") || strings.Contains(path, "<(") {
		return "", false
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", false
	}

	resolved := true
	path = os.Expand(path, func(name string) string {
		if name == "HOME" {
			return home
		}
		val, ok := os.LookupEnv(name)
		if !ok || val == "" {
			resolved = false
			return ""
		}
		return val
	})
	if !resolved || path == "" {
		return "", false
	}

	if path == "~" {
		path = home
	} else if strings.HasPrefix(path, "~/") {
		path = filepath.Join(home, path[2:])
	}

	if !filepath.IsAbs(path) {
		return "", false
	}
	return filepath.Clean(path), true
}

// unquoteAliasValue removes one matched pair of surrounding quotes.
//
// strings.Trim takes a cutset and strips greedily from both ends, so a value
// whose last character is itself a quote lost it:
//
//	alias gwip="git add --all && git commit -am 'WIP'"
//	  -> git add --all && git commit -am 'WIP
//
// That is not only a wrong description in the menu. With core.expand-alias on,
// the target is typed back into the prompt when the alias is expanded, so the
// user was left holding a command line with an unbalanced quote.
func unquoteAliasValue(s string) string {
	if len(s) < 2 {
		return s
	}
	if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
		return s[1 : len(s)-1]
	}
	return s
}

func SplitAliasTokens(s string) []string {
	var tokens []string
	var cur strings.Builder
	inQuote := false
	var quoteChar rune
	for _, c := range s {
		switch {
		case !inQuote && (c == '"' || c == '\''):
			inQuote = true
			quoteChar = c
			cur.WriteRune(c)
		case inQuote && c == quoteChar:
			inQuote = false
			cur.WriteRune(c)
		case !inQuote && c == ' ':
			if cur.Len() > 0 {
				tokens = append(tokens, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteRune(c)
		}
	}
	if cur.Len() > 0 {
		tokens = append(tokens, cur.String())
	}
	return tokens
}
