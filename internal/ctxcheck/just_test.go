package ctxcheck

import (
	"os"
	"path/filepath"
	"testing"
)

func TestJust_NoJustfile(t *testing.T) {
	tmpDir := t.TempDir()
	v := ValidateJust([]string{"just", "test"}, tmpDir)
	if v != Invalid {
		t.Fatalf("expected Invalid when no justfile exists, got %v", v)
	}
}

func TestJust_BareJust(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tmpDir, "justfile"), []byte("default:\n\techo hello\n"), 0644)
	v := ValidateJust([]string{"just"}, tmpDir)
	if v != Valid {
		t.Fatalf("expected Valid for bare just, got %v", v)
	}
}

func TestJust_RecipeHeadersAndParams(t *testing.T) {
	tmpDir := t.TempDir()
	content := `
# comment
set shell := ["bash", "-c"]
export VAR := "val"
foo := "ignored_var"

build:
	echo building

@test target:
	echo testing {{target}}

deploy env='prod' +services:
	echo deploying

run *args:
	echo running

_private_recipe:
	echo internal
`
	_ = os.WriteFile(filepath.Join(tmpDir, "justfile"), []byte(content), 0644)

	cases := []struct {
		cmd      []string
		expected Verdict
	}{
		{[]string{"just", "build"}, Valid},
		{[]string{"just", "test"}, Valid},
		{[]string{"just", "test", "my-target"}, Valid},
		{[]string{"just", "deploy"}, Valid},
		{[]string{"just", "run", "arg1", "arg2"}, Valid},
		{[]string{"just", "_private_recipe"}, Valid},
		{[]string{"just", "nonexistent"}, Invalid},
		{[]string{"just", "foo"}, Invalid}, // foo is a variable, not a recipe
		{[]string{"just", "VAR"}, Invalid}, // VAR is an export, not a recipe
	}

	for _, tc := range cases {
		v := ValidateJust(tc.cmd, tmpDir)
		if v != tc.expected {
			t.Errorf("cmd %v: expected %v, got %v", tc.cmd, tc.expected, v)
		}
	}
}

func TestJust_Alias(t *testing.T) {
	tmpDir := t.TempDir()
	content := `
test:
	echo test

alias t := test
alias check := test
`
	_ = os.WriteFile(filepath.Join(tmpDir, "Justfile"), []byte(content), 0644)

	if v := ValidateJust([]string{"just", "t"}, tmpDir); v != Valid {
		t.Errorf("expected alias 't' to be Valid, got %v", v)
	}
	if v := ValidateJust([]string{"just", "check"}, tmpDir); v != Valid {
		t.Errorf("expected alias 'check' to be Valid, got %v", v)
	}
	if v := ValidateJust([]string{"just", "other"}, tmpDir); v != Invalid {
		t.Errorf("expected 'other' to be Invalid, got %v", v)
	}
}

func TestJust_ImportAndMod(t *testing.T) {
	tmpDir := t.TempDir()
	content := `
import "common.just"
mod sub "submodules/sub"

build:
	echo build
`
	_ = os.WriteFile(filepath.Join(tmpDir, ".justfile"), []byte(content), 0644)

	// known recipe in current file
	if v := ValidateJust([]string{"just", "build"}, tmpDir); v != Valid {
		t.Errorf("expected 'build' to be Valid, got %v", v)
	}

	// token matches mod name -> Unknown
	if v := ValidateJust([]string{"just", "sub", "task"}, tmpDir); v != Unknown {
		t.Errorf("expected mod name 'sub' to be Unknown, got %v", v)
	}

	// missing recipe when import exists -> Unknown (not Invalid)
	if v := ValidateJust([]string{"just", "imported_recipe"}, tmpDir); v != Unknown {
		t.Errorf("expected missing recipe to be Unknown when imports exist, got %v", v)
	}
}

func TestJust_SetFallback(t *testing.T) {
	tmpDir := t.TempDir()
	content := `
set fallback := true

local:
	echo local
`
	_ = os.WriteFile(filepath.Join(tmpDir, "justfile"), []byte(content), 0644)

	if v := ValidateJust([]string{"just", "local"}, tmpDir); v != Valid {
		t.Errorf("expected 'local' to be Valid, got %v", v)
	}

	// missing recipe with set fallback -> Unknown
	if v := ValidateJust([]string{"just", "parent_recipe"}, tmpDir); v != Unknown {
		t.Errorf("expected missing recipe with fallback to be Unknown, got %v", v)
	}
}

func TestJust_UpwardSearch(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tmpDir, "justfile"), []byte("root_task:\n\techo root\n"), 0644)

	subDir := filepath.Join(tmpDir, "src", "pkg", "deep")
	_ = os.MkdirAll(subDir, 0755)

	if v := ValidateJust([]string{"just", "root_task"}, subDir); v != Valid {
		t.Errorf("expected upward search to find 'root_task', got %v", v)
	}
	if v := ValidateJust([]string{"just", "missing"}, subDir); v != Invalid {
		t.Errorf("expected missing recipe to be Invalid from subfolder, got %v", v)
	}
}

func TestJust_FlagsAndComplexTokens(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tmpDir, "justfile"), []byte("build:\n\techo build\n"), 0644)

	// regular flags stripped before recipe
	if v := ValidateJust([]string{"just", "-q", "build"}, tmpDir); v != Valid {
		t.Errorf("expected 'just -q build' to be Valid, got %v", v)
	}
	if v := ValidateJust([]string{"just", "--quiet", "build"}, tmpDir); v != Valid {
		t.Errorf("expected 'just --quiet build' to be Valid, got %v", v)
	}

	// redirect flags -> Unknown
	if v := ValidateJust([]string{"just", "-f", "other.just", "build"}, tmpDir); v != Unknown {
		t.Errorf("expected -f flag to be Unknown, got %v", v)
	}
	if v := ValidateJust([]string{"just", "-d", "other_dir", "build"}, tmpDir); v != Unknown {
		t.Errorf("expected -d flag to be Unknown, got %v", v)
	}

	// colon or slash in target -> Unknown
	if v := ValidateJust([]string{"just", "sub::build"}, tmpDir); v != Unknown {
		t.Errorf("expected '::' in target to be Unknown, got %v", v)
	}
	if v := ValidateJust([]string{"just", "dir/build"}, tmpDir); v != Unknown {
		t.Errorf("expected '/' in target to be Unknown, got %v", v)
	}

	// VAR=val token -> Unknown
	if v := ValidateJust([]string{"just", "FOO=bar", "build"}, tmpDir); v != Unknown {
		t.Errorf("expected VAR=val token to be Unknown, got %v", v)
	}
}
