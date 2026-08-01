package shell

import (
	"context"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ReplaceLine builds the byte sequence that clears the shell's current input
// line and types text in its place.
//
// The clear is ctrl+e then ctrl+u, not ctrl+u alone. ctrl+u is only
// kill-whole-line under zsh's stock emacs keymap; `bindkey '^U'
// backward-kill-line` is a common preference, and it is also what bash does by
// default (unix-line-discard). Those kill from the cursor backwards, so with
// the cursor mid-line everything to its right survived and collided with the
// text typed next -- selecting a suggestion left the tail of the old buffer
// interleaved through the new one.
//
// Moving to end-of-line first makes all three bindings equivalent: there is
// nothing to the right left to preserve.
func ReplaceLine(text []byte) []byte {
	out := make([]byte, 0, len(text)+2)
	out = append(out, 0x05, 0x15) // ctrl+e (end-of-line), ctrl+u (kill)
	return append(out, text...)
}

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
	return ReplaceLine([]byte(selected))
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
	return ReplaceLine([]byte(selected))
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
	
	// Fallback: ask zsh directly in case ZDOTDIR is set in ~/.zshenv without export
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	cmd := exec.CommandContext(ctx, "zsh", "-c", "echo $ZDOTDIR")
	out, err := cmd.Output()
	if err == nil {
		zdotdir := strings.TrimSpace(string(out))
		if zdotdir != "" {
			return zdotdir
		}
	}
	
	home, _ := os.UserHomeDir()
	return home
}

// FishAdapter implementation
type FishAdapter struct{}

func (f *FishAdapter) GetName() string      { return "fish" }
func (f *FishAdapter) GetShellPath() string { return "fish" }
func (f *FishAdapter) GetEnv(fd int, pid int) []string {
	return append(os.Environ(), "IRIS_FD="+fmt.Sprint(fd), "IRIS_PID="+fmt.Sprint(pid))
}
func (f *FishAdapter) PrepareSelectSequence(selected string) []byte {
	return ReplaceLine([]byte(selected))
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

func ScanPosixAliases(files []string) map[string]string {
	aliases := make(map[string]string)
	home, err := os.UserHomeDir()
	if err != nil {
		return aliases
	}

	for _, f := range files {
		path := f
		if !filepath.IsAbs(f) {
			path = filepath.Join(home, f)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		maps.Copy(aliases, ParseAliases(string(data)))
	}
	return aliases
}

func ParseAliases(data string) map[string]string {
	aliases := make(map[string]string)
	lines := strings.SplitSeq(data, "\n")
	for line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "alias ") {
			continue
		}
		body := strings.TrimSpace(strings.TrimPrefix(line, "alias"))
		if body == "" {
			continue
		}

		pairs := SplitAliasTokens(body)
		for _, pair := range pairs {
			before, after, ok := strings.Cut(pair, "=")
			if !ok {
				continue
			}
			key := strings.TrimSpace(before)
			val := strings.Trim(strings.TrimSpace(after), `"'`)
			if key != "" && val != "" {
				aliases[key] = val
			}
		}
	}
	return aliases
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
