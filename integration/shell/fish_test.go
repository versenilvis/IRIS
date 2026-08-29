package shell

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestParseFishConfigAliasForms(t *testing.T) {
	input := `
alias clear 'clear -x'
alias g='git'
alias ls 'lsd --icon=always'
alias --save gwip "git add --all && git commit -am 'WIP'"
alias -- ll 'ls -l'
alias 'k=kubectl get'
alias gc git commit --verbose
# alias hidden 'not found'
alias c 'echo hi' # trailing comment
`
	expected := map[string]string{
		"clear": "clear -x",
		"g":     "git",
		"ls":    "lsd --icon=always",
		"gwip":  "git add --all && git commit -am 'WIP'",
		"ll":    "ls -l",
		"k":     "kubectl get",
		"gc":    "git commit --verbose",
		"c":     "echo hi",
	}

	got, _ := ParseFishConfig("", input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("ParseFishConfig() aliases = %#v; want %#v", got, expected)
	}
}

func TestParseFishConfigAbbrForms(t *testing.T) {
	input := `
abbr -a gl git log
abbr --add gco 'git checkout'
abbr gp git push
abbr -a -- xg4 'git commit -m'
abbr -a --position anywhere L '| less'
abbr -a --set-cursor gcm "git commit -m '%'"
abbr -a -g gd git diff
abbr -a --regex '[0-9]+' dynamic something
abbr -a --function expand_dots dots
abbr --erase gone
abbr --rename old new
abbr -a '!$' '$history[1]'
abbr -a --command git -- pull 'git pull'
abbr -a --command $tool -- \
    (string split -f1 -m1 ' ' -- $entry) \
    (string split -f2 -m1 ' ' -- $entry)
`
	expected := map[string]string{
		"gl":   "git log",
		"gco":  "git checkout",
		"gp":   "git push",
		"xg4":  "git commit -m",
		"L":    "| less",
		"gcm":  "git commit -m '%'",
		"gd":   "git diff",
		"!$":   "$history[1]",
		"pull": "git pull",
	}

	_, got := ParseFishConfig("", input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("ParseFishConfig() abbrs = %#v; want %#v", got, expected)
	}
}

func TestParseFishConfigPrefersFirstConditionalArm(t *testing.T) {
	// `end` closes every kind of fish block, so the inner `for` and `function`
	// must not pop the enclosing `if` and let the fallback arm win.
	input := `
if type -q lsd
    for f in a b
        echo $f
    end
    alias ls 'lsd --icons'
    alias lt 'lsd --tree'
else
    alias ls 'ls --color=auto'
    alias lt 'tree'
end

function greet
    echo hi
end

switch (uname)
    case Darwin
        alias o 'open'
    case Linux
        alias o 'xdg-open'
end

alias g 'git'
`
	expected := map[string]string{
		"ls": "lsd --icons",
		"lt": "lsd --tree",
		"o":  "open",
		"g":  "git",
	}

	got, _ := ParseFishConfig("", input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("ParseFishConfig() aliases = %#v; want %#v", got, expected)
	}
}

func TestParseFishConfigStatementSeparators(t *testing.T) {
	input := `
type -q eza; and alias ls 'eza --icons'
type -q bat; or alias cat 'bat --plain'
alias a 'echo a'; alias b 'echo b'
command alias c 'echo c'
`
	expected := map[string]string{
		"ls":  "eza --icons",
		"cat": "bat --plain",
		"a":   "echo a",
		"b":   "echo b",
		"c":   "echo c",
	}

	got, _ := ParseFishConfig("", input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("ParseFishConfig() aliases = %#v; want %#v", got, expected)
	}
}

func TestParseFishConfigReadsSavedAliasFunctions(t *testing.T) {
	// The exact shape `alias --save` writes, in both the space and equals forms
	// fish has used for the description.
	input := `
function gst --wraps='git status -sb' --description 'alias gst git status -sb'
    git status -sb $argv
end

function g --wraps='git status' --description 'alias g=git status'
    git status $argv
end

function lsd --description 'alias lsd lsd --icon=always'
    command lsd --icon=always $argv
end

function greet --description 'say hello'
    echo hello
end
`
	expected := map[string]string{
		"gst": "git status -sb",
		"g":   "git status",
		"lsd": "lsd --icon=always",
	}

	got, _ := ParseFishConfig("", input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("ParseFishConfig() aliases = %#v; want %#v", got, expected)
	}
}

func TestScanFishDefsCoversAutoloadedDirs(t *testing.T) {
	dir := t.TempDir()
	writeFish(t, filepath.Join(dir, "config.fish"), "alias g 'git'\nabbr -a gl git log\n")
	writeFish(t, filepath.Join(dir, "conf.d", "10-tools.fish"), "alias k 'kubectl'\nabbr -a kg 'kubectl get'\n")
	writeFish(t, filepath.Join(dir, "functions", "gst.fish"),
		"function gst --wraps='git status' --description 'alias gst git status'\n    git status $argv\nend\n")

	defs := scanFishDefs(dir)

	wantAliases := map[string]string{"g": "git", "k": "kubectl", "gst": "git status"}
	if !reflect.DeepEqual(defs.aliases, wantAliases) {
		t.Errorf("scanFishDefs() aliases = %#v; want %#v", defs.aliases, wantAliases)
	}
	wantAbbrs := map[string]string{"gl": "git log", "kg": "kubectl get"}
	if !reflect.DeepEqual(defs.abbrs, wantAbbrs) {
		t.Errorf("scanFishDefs() abbrs = %#v; want %#v", defs.abbrs, wantAbbrs)
	}
}

func TestScanFishDefsPrefersConfigOverAutoloadedFunction(t *testing.T) {
	// An autoloaded function never runs while the name is already defined, so
	// config.fish is what the shell actually ends up with.
	dir := t.TempDir()
	writeFish(t, filepath.Join(dir, "config.fish"), "alias gst 'git status -sb'\n")
	writeFish(t, filepath.Join(dir, "functions", "gst.fish"),
		"function gst --description 'alias gst git status'\n    git status $argv\nend\n")

	if got := scanFishDefs(dir).aliases["gst"]; got != "git status -sb" {
		t.Errorf("scanFishDefs() gst = %q; want %q", got, "git status -sb")
	}
}

func TestScanFishDefsFollowsSource(t *testing.T) {
	dir := t.TempDir()
	sourced := filepath.Join(dir, "aliases.fish")
	writeFish(t, sourced, "alias v 'nvim'\n")
	writeFish(t, filepath.Join(dir, "config.fish"),
		"source \""+sourced+"\"\nsource (fzf --fish | psub)\nsource \"$IRIS_MISSING_VAR/x.fish\"\n")

	got := scanFishDefs(dir).aliases
	if !reflect.DeepEqual(got, map[string]string{"v": "nvim"}) {
		t.Errorf("scanFishDefs() aliases = %#v; want v=nvim only", got)
	}
}

func TestScanFishDefsResolvesLoopSourcedFiles(t *testing.T) {
	// A config that names its parts in a `set` list and sources them from a
	// `for` loop reaches every alias only if both are resolved.
	dir := t.TempDir()
	writeFish(t, filepath.Join(dir, "config", "aliases.fish"), "alias g 'git'\n")
	writeFish(t, filepath.Join(dir, "config", "abbr.fish"), "abbr -a gl git log\n")
	writeFish(t, filepath.Join(dir, "config", "prompt.fish"), "alias p 'starship prompt'\n")
	writeFish(t, filepath.Join(dir, "config.fish"), `
source "$__fish_config_dir/config/prompt.fish"

function __source_fish_config_files
    set -l files \
        aliases \
        abbr

    for file in $files
        set -l config_file "$__fish_config_dir/config/$file.fish"
        path is -f -- "$config_file"; and source "$config_file"
    end
end

__source_fish_config_files
`)

	defs := scanFishDefs(dir)

	wantAliases := map[string]string{"g": "git", "p": "starship prompt"}
	if !reflect.DeepEqual(defs.aliases, wantAliases) {
		t.Errorf("scanFishDefs() aliases = %#v; want %#v", defs.aliases, wantAliases)
	}
	if !reflect.DeepEqual(defs.abbrs, map[string]string{"gl": "git log"}) {
		t.Errorf("scanFishDefs() abbrs = %#v; want gl=git log", defs.abbrs)
	}
}

func TestScanFishDefsInvalidatesOnNewConfDFile(t *testing.T) {
	dir := t.TempDir()
	writeFish(t, filepath.Join(dir, "config.fish"), "alias g 'git'\n")
	writeFish(t, filepath.Join(dir, "conf.d", "10-tools.fish"), "alias k 'kubectl'\n")

	if got := len(scanFishDefs(dir).aliases); got != 2 {
		t.Fatalf("scanFishDefs() = %d aliases; want 2", got)
	}

	// A file appearing in conf.d/ leaves every scanned file byte-identical, so
	// only the directory stamp can catch it.
	added := filepath.Join(dir, "conf.d", "20-more.fish")
	writeFish(t, added, "alias d 'docker'\n")
	if err := os.Chtimes(filepath.Join(dir, "conf.d"), time.Now().Add(time.Second), time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}

	got := scanFishDefs(dir).aliases
	if got["d"] != "docker" {
		t.Errorf("scanFishDefs() after adding %s = %#v; want d=docker", added, got)
	}
}

func TestScanFishDefsReturnsAnIndependentCopy(t *testing.T) {
	dir := t.TempDir()
	writeFish(t, filepath.Join(dir, "config.fish"), "alias g 'git'\nabbr -a gl git log\n")

	defs := scanFishDefs(dir)
	defs.aliases["injected"] = "nope"
	defs.abbrs["injected"] = "nope"

	again := scanFishDefs(dir)
	if len(again.aliases) != 1 || len(again.abbrs) != 1 {
		t.Errorf("cache returned a mutable reference: %#v", again)
	}
}

func TestFishAdapterImplementsAbbrScanner(t *testing.T) {
	if _, ok := Adapter(&FishAdapter{}).(AbbrScanner); !ok {
		t.Error("FishAdapter does not implement AbbrScanner")
	}
	for _, a := range []Adapter{&BashAdapter{}, &ZshAdapter{}} {
		if _, ok := a.(AbbrScanner); ok {
			t.Errorf("%s implements AbbrScanner; only fish has abbreviations", a.GetName())
		}
	}
}

func writeFish(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
