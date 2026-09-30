package ctxcheck

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMake_NoMakefile(t *testing.T) {
	tmpDir := t.TempDir()
	if v := ValidateMake([]string{"make", "build"}, tmpDir); v != Invalid {
		t.Fatalf("expected Invalid when no makefile in cwd, got %v", v)
	}
}

func TestMake_InCwdOnly(t *testing.T) {
	tmpDir := t.TempDir()
	// Makefile in parent directory
	_ = os.WriteFile(filepath.Join(tmpDir, "Makefile"), []byte("build:\n\techo build\n"), 0644)

	subDir := filepath.Join(tmpDir, "sub")
	_ = os.MkdirAll(subDir, 0755)

	// make does not search upward: should be Invalid in subDir
	if v := ValidateMake([]string{"make", "build"}, subDir); v != Invalid {
		t.Fatalf("expected Invalid when makefile is only in parent directory, got %v", v)
	}

	// in tmpDir, it is Valid
	if v := ValidateMake([]string{"make", "build"}, tmpDir); v != Valid {
		t.Fatalf("expected Valid in directory containing Makefile, got %v", v)
	}
}

func TestMake_TargetsAndBareMake(t *testing.T) {
	tmpDir := t.TempDir()
	content := `
all: build test

build:
	echo building

test:
	echo testing
`
	_ = os.WriteFile(filepath.Join(tmpDir, "GNUmakefile"), []byte(content), 0644)

	if v := ValidateMake([]string{"make"}, tmpDir); v != Valid {
		t.Errorf("expected bare make to be Valid, got %v", v)
	}
	if v := ValidateMake([]string{"make", "all"}, tmpDir); v != Valid {
		t.Errorf("expected 'all' target to be Valid, got %v", v)
	}
	if v := ValidateMake([]string{"make", "missing"}, tmpDir); v != Invalid {
		t.Errorf("expected missing target to be Invalid, got %v", v)
	}
}

func TestMake_IncludeAndPatternRules(t *testing.T) {
	tmpDir := t.TempDir()
	content := `
include common.mk

%.o: %.c
	gcc -c $<

build:
	echo building
`
	_ = os.WriteFile(filepath.Join(tmpDir, "Makefile"), []byte(content), 0644)

	// known target
	if v := ValidateMake([]string{"make", "build"}, tmpDir); v != Valid {
		t.Errorf("expected 'build' to be Valid, got %v", v)
	}

	// target not in file, but file has include / pattern rules -> Unknown (not Invalid)
	if v := ValidateMake([]string{"make", "other"}, tmpDir); v != Unknown {
		t.Errorf("expected missing target with include/pattern to be Unknown, got %v", v)
	}
}

func TestMake_Flags(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tmpDir, "Makefile"), []byte("build:\n\techo build\n"), 0644)

	if v := ValidateMake([]string{"make", "-C", "sub", "build"}, tmpDir); v != Unknown {
		t.Errorf("expected -C flag to be Unknown, got %v", v)
	}
	if v := ValidateMake([]string{"make", "-Csub", "build"}, tmpDir); v != Unknown {
		t.Errorf("expected -Csub flag to be Unknown, got %v", v)
	}
	if v := ValidateMake([]string{"make", "-f", "other.mk", "build"}, tmpDir); v != Unknown {
		t.Errorf("expected -f flag to be Unknown, got %v", v)
	}
	if v := ValidateMake([]string{"make", "--file=other.mk", "build"}, tmpDir); v != Unknown {
		t.Errorf("expected --file= flag to be Unknown, got %v", v)
	}
	if v := ValidateMake([]string{"make", "--makefile=other.mk", "build"}, tmpDir); v != Unknown {
		t.Errorf("expected --makefile= flag to be Unknown, got %v", v)
	}
	if v := ValidateMake([]string{"make", "--directory=sub", "build"}, tmpDir); v != Unknown {
		t.Errorf("expected --directory= flag to be Unknown, got %v", v)
	}

	for _, flag := range []string{"-j", "-l", "-o", "-W", "-I"} {
		if v := ValidateMake([]string{"make", flag, "4", "build"}, tmpDir); v != Valid {
			t.Errorf("expected %s with separate value and 'build' to be Valid, got %v", flag, v)
		}
		if v := ValidateMake([]string{"make", flag, "4", "missing"}, tmpDir); v != Invalid {
			t.Errorf("expected %s with separate value and 'missing' to be Invalid, got %v", flag, v)
		}
	}
}
