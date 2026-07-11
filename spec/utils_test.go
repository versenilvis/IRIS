package spec

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"Empty input", "", []string{""}},
		{"Trailing space", "git ", []string{"git", ""}},
		{"Multi-space", "git  add", []string{"git", "add"}},
		{"Quoted string", "git commit -m \"hello world\"", []string{"git", "commit", "-m", "hello world"}},
		{"Quote not closed", "git commit -m \"hello", []string{"git", "commit", "-m", "hello"}},
		{"Single quote", "git commit -m 'hello world'", []string{"git", "commit", "-m", "hello world"}},
		{"Backslash escape", "git commit -m \"hello\\ world\"", []string{"git", "commit", "-m", "hello world"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Tokenize(tt.input)
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
		{"Case insensitive", "Hello", "hel", true},
		{"Unicode support", "Thử nghiệm", "thử", true},
		{"Prefix longer", "Iris", "Iris-Longer", false},
		{"Empty prefix", "Iris", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasPrefix(tt.s, tt.prefix); got != tt.want {
				t.Errorf("HasPrefix(%q, %q) = %v; want %v", tt.s, tt.prefix, got, tt.want)
			}
		})
	}
}
