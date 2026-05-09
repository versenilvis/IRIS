package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	return ScanPosixAliases([]string{".zshrc", ".zshenv", ".zprofile"})
}

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
	return ScanPosixAliases([]string{filepath.Join(".config", "fish", "config.fish")})
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

		for k, v := range ParseAliases(string(data)) {
			aliases[k] = v
		}
	}
	return aliases
}

func ParseAliases(data string) map[string]string {
	aliases := make(map[string]string)
	lines := strings.Split(data, "\n")
	for _, line := range lines {
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
			eqIdx := strings.IndexByte(pair, '=')
			if eqIdx < 0 {
				continue
			}
			key := strings.TrimSpace(pair[:eqIdx])
			val := strings.Trim(strings.TrimSpace(pair[eqIdx+1:]), `"'`)
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
