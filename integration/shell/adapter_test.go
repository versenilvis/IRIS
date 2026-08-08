package shell

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
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

func TestParseAliasesPrefersFirstConditionalArm(t *testing.T) {
	// A dotfile that picks a tool with a fallback defines the same alias twice.
	// The first arm is the preferred one; letting the fallback win just because
	// it is written last points suggestions at the wrong command.
	input := `
if which eza &>/dev/null; then
  alias ls='eza --icons'
  alias lt='eza -T'
else
  alias ls='lsd'
  alias lt='lsd --tree'
fi
alias g='git'
`
	expected := map[string]string{
		"ls": "eza --icons",
		"lt": "eza -T",
		"g":  "git",
	}

	got := ParseAliases(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("ParseAliases() = %v; want %v", got, expected)
	}
}

func TestScanPosixAliasesFollowsSource(t *testing.T) {
	dir := t.TempDir()

	aliasFile := filepath.Join(dir, "aliases.zsh")
	if err := os.WriteFile(aliasFile, []byte("alias gst='git status'\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	nestedFile := filepath.Join(dir, "nested.zsh")
	if err := os.WriteFile(nestedFile, []byte("alias k='kubectl'\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	// The guarded form and the `.` spelling both appear in real configs, as does
	// a source line that can't be resolved statically -- that one must be
	// skipped rather than guessed at.
	rc := filepath.Join(dir, ".zshrc")
	content := "alias v='nvim'\n" +
		"source \"$IRIS_TEST_DIR/aliases.zsh\"\n" +
		"[ -f \"" + nestedFile + "\" ] && . \"" + nestedFile + "\"\n" +
		"for f in $dir/*.zsh; do source \"$f\"; done\n" +
		"source <(fzf --zsh)\n"
	if err := os.WriteFile(rc, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("IRIS_TEST_DIR", dir)

	got := ScanPosixAliases([]string{rc})
	expected := map[string]string{
		"v":   "nvim",
		"gst": "git status",
		"k":   "kubectl",
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("ScanPosixAliases() = %v; want %v", got, expected)
	}
}

func TestScanPosixAliasesCacheInvalidation(t *testing.T) {
	dir := t.TempDir()
	rc := filepath.Join(dir, ".zshrc")
	sourced := filepath.Join(dir, "aliases.zsh")

	if err := os.WriteFile(rc, []byte("source \""+sourced+"\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	// The sourced file does not exist yet. Its later appearance has to
	// invalidate the cache even though the file that sources it never changes.
	if got := ScanPosixAliases([]string{rc}); len(got) != 0 {
		t.Fatalf("ScanPosixAliases() = %v; want empty", got)
	}

	if err := os.WriteFile(sourced, []byte("alias gst='git status'\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	got := ScanPosixAliases([]string{rc})
	if !reflect.DeepEqual(got, map[string]string{"gst": "git status"}) {
		t.Fatalf("after create: ScanPosixAliases() = %v; want gst", got)
	}

	// Editing a sourced file must be picked up as well. Size differs here, but
	// keep the mtime bump explicit so the check doesn't rest on that alone.
	if err := os.WriteFile(sourced, []byte("alias gst='git status -sb'\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(sourced, time.Now().Add(time.Second), time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	got = ScanPosixAliases([]string{rc})
	if !reflect.DeepEqual(got, map[string]string{"gst": "git status -sb"}) {
		t.Fatalf("after edit: ScanPosixAliases() = %v; want updated gst", got)
	}

	// A cached result must not be shared by reference; a caller mutating it
	// would otherwise poison every later scan.
	got["injected"] = "nope"
	if again := ScanPosixAliases([]string{rc}); len(again) != 1 {
		t.Fatalf("cache returned a mutable reference: %v", again)
	}
}

func TestScanPosixAliasesHandlesSourceCycle(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.zsh")
	b := filepath.Join(dir, "b.zsh")

	if err := os.WriteFile(a, []byte("alias a='echo a'\nsource \""+b+"\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(b, []byte("alias b='echo b'\nsource \""+a+"\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	got := ScanPosixAliases([]string{a})
	expected := map[string]string{"a": "echo a", "b": "echo b"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("ScanPosixAliases() = %v; want %v", got, expected)
	}
}
func TestParseAliasesKeepsNestedQuotes(t *testing.T) {
	// strings.Trim with a cutset strips greedily, so a value ending in a quote
	// lost it. With core.expand-alias on the target is typed back into the
	// prompt, so this left an unbalanced quote on the command line.
	input := `
alias gwip="git add --all && git commit -am 'WIP'"
alias shrug="echo '¯\_(ツ)_/¯' | pbcopy"
alias plain=mcd
alias spaced='git status -sb'
`
	expected := map[string]string{
		"gwip":   "git add --all && git commit -am 'WIP'",
		"shrug":  `echo '¯\_(ツ)_/¯' | pbcopy`,
		"plain":  "mcd",
		"spaced": "git status -sb",
	}

	got := ParseAliases(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("ParseAliases() = %#v; want %#v", got, expected)
	}
}
