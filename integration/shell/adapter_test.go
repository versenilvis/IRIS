package shell

import (
	"reflect"
	"testing"
)

func TestScanPosixAliases(t *testing.T) {

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

	got := ParseAliases(input)
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

		{"Single", "a='b'", []string{"a='b'"}},
		{"Multi", "a='b' c=\"d\"", []string{"a='b'", "c=\"d\""}},

		{"With Space", "ta='tmux a -t' l='ls -l'", []string{"ta='tmux a -t'", "l='ls -l'"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitAliasTokens(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("SplitAliasTokens(%q) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}
func TestReplaceLineMovesToEndBeforeKilling(t *testing.T) {
	// ctrl+u alone is only kill-whole-line under zsh's stock emacs keymap.
	// Under `bindkey '^U' backward-kill-line` (and bash's default
	// unix-line-discard) it kills backwards from the cursor, so anything to the
	// right of the cursor survived and collided with the text typed next.
	// Prefixing ctrl+e makes the three equivalent.
	got := ReplaceLine([]byte("cd config/"))

	if len(got) < 2 || got[0] != 0x05 || got[1] != 0x15 {
		t.Fatalf("ReplaceLine() = %q; want it to start with ctrl+e, ctrl+u", got)
	}
	if string(got[2:]) != "cd config/" {
		t.Errorf("ReplaceLine() payload = %q; want %q", got[2:], "cd config/")
	}
}
func TestPrepareSelectSequenceUsesReplaceLine(t *testing.T) {
	for _, a := range []Adapter{&ZshAdapter{}, &BashAdapter{}, &FishAdapter{}} {
		got := a.PrepareSelectSequence("git status")
		want := ReplaceLine([]byte("git status"))
		if string(got) != string(want) {
			t.Errorf("%s PrepareSelectSequence() = %q; want %q", a.GetName(), got, want)
		}
	}
}
