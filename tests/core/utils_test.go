package tests

import (
	"reflect"
	"testing"

	"github.com/versenilvis/iris/commands/core"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		// REQUIREMENT: Empty input
		{"Empty input", "", []string{""}},
		// REQUIREMENT: Trailing space("git " -> 2 tokens, the last token is "")
		{"Trailing space", "git ", []string{"git", ""}},
		// REQUIREMENT: Multi-space("git add")
		{"Multi-space", "git  add", []string{"git", "add"}},
		// REQUIREMENT: Quoted string("git commit -m \"hello world\"")
		{"Quoted string", "git commit -m \"hello world\"", []string{"git", "commit", "-m", "hello world"}},
		// REQUIREMENT: Quote not closed
		{"Quote not closed", "git commit -m \"hello", []string{"git", "commit", "-m", "hello"}},
		// REQUIREMENT: Single quote vs double quote
		{"Single quote", "git commit -m 'hello world'", []string{"git", "commit", "-m", "hello world"}},
		// REQUIREMENT: Backslash escape
		{"Backslash escape", "git commit -m \"hello\\ world\"", []string{"git", "commit", "-m", "hello world"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := core.Tokenize(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Tokenize(%q) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestHasPrefix(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		prefix string
		want   bool
	}{
		// REQUIREMENT: Case insensitive match
		{"Case insensitive", "Hello", "hel", true},
		// REQUIREMENT: Unicode support (Vietnamese)
		{"Unicode support", "Thử nghiệm", "thử", true},
		// REQUIREMENT: Prefix longer than string -> false
		{"Prefix longer", "Iris", "Iris-Longer", false},
		// REQUIREMENT: Empty prefix -> true
		{"Empty prefix", "Iris", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := core.HasPrefix(tt.s, tt.prefix); got != tt.want {
				t.Errorf("HasPrefix(%q, %q) = %v; want %v", tt.s, tt.prefix, got, tt.want)
			}
		})
	}
}
