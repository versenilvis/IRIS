package shell_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/versenilvis/iris/integration/shell"
)

func TestSplitAliasTokens(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"g='git commit'", []string{"g='git commit'"}},
		{"a=b c=d", []string{"a=b", "c=d"}},
		{"ta='tmux a -t' l='ls -l'", []string{"ta='tmux a -t'", "l='ls -l'"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := shell.SplitAliasTokens(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("SplitAliasTokens(%q) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestScanPosixAliases(t *testing.T) {
	tmp := t.TempDir()
	
	aliasFile := filepath.Join(tmp, ".bashrc")
	content := `
# some comments
alias g='git'
alias gca='git commit -a'
alias multi="a" b="c"
`
	os.WriteFile(aliasFile, []byte(content), 0644)

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", oldHome)

	aliases := shell.ScanPosixAliases([]string{".bashrc"})
	
	expected := map[string]string{
		"g":     "git",
		"gca":   "git commit -a",
		"multi": "a",
		"b":     "c",
	}

	for k, v := range expected {
		if aliases[k] != v {
			t.Errorf("Expected alias %s=%s, got %s", k, v, aliases[k])
		}
	}
}
