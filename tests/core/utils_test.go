package core_test

import (
	"reflect"
	"testing"

	"github.com/versenilvis/iris/commands/core"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"", []string{""}},
		{"git ", []string{"git", ""}},
		{"git  add", []string{"git", "add"}},
		{"git commit -m \"hello world\"", []string{"git", "commit", "-m", "hello world"}},
		{"git commit -m \"hello", []string{"git", "commit", "-m", "hello"}},
		{"git commit -m 'hello world'", []string{"git", "commit", "-m", "hello world"}},
		{"git commit -m 'hello", []string{"git", "commit", "-m", "hello"}},
		{"ls -l \"file name\"", []string{"ls", "-l", "file name"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := core.Tokenize(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Tokenize(%q) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestHasPrefix(t *testing.T) {
	tests := []struct {
		s      string
		prefix string
		want   bool
	}{
		{"Hello", "hel", true},
		{"Thử nghiệm", "thử", true},
		{"Iris", "Iris-Longer", false},
		{"Iris", "", true},
		{"", "a", false},
	}

	for _, tt := range tests {
		t.Run(tt.s+"_"+tt.prefix, func(t *testing.T) {
			if got := core.HasPrefix(tt.s, tt.prefix); got != tt.want {
				t.Errorf("HasPrefix(%q, %q) = %v; want %v", tt.s, tt.prefix, got, tt.want)
			}
		})
	}
}
