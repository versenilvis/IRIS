package tests

import (
	"reflect"
	"testing"

	"github.com/versenilvis/iris/integration/shell"
)

func TestScanPosixAliases(t *testing.T) {
	// REQUIREMENT: Parse single alias, multi alias, comments, value with space
	input := `
alias gca='git commit -a'
alias ta="tmux a -t" # this is a comment
# alias hidden="not found"
alias l='ls' ll='ls -l'
`
	expected := map[string]string{
		"gca": "git commit -a",
		"ta":  "tmux a -t",
		"l":   "ls",
		"ll":  "ls -l",
	}

	got := shell.ParseAliases(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("ScanPosixAliases() = %v; want %v", got, expected)
	}
}

func TestSplitAliasTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		// REQUIREMENT: Parse multi alias on one line
		{"Single", "a='b'", []string{"a='b'"}},
		{"Multi", "a='b' c=\"d\"", []string{"a='b'", "c=\"d\""}},
		// REQUIREMENT: Value with space in quote
		{"With Space", "ta='tmux a -t' l='ls -l'", []string{"ta='tmux a -t'", "l='ls -l'"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shell.SplitAliasTokens(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("SplitAliasTokens(%q) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}
