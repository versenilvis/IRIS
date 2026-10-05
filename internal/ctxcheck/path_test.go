package ctxcheck

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPath_Cd(t *testing.T) {
	tmpDir := t.TempDir()
	existingDir := filepath.Join(tmpDir, "myfolder")
	_ = os.Mkdir(existingDir, 0755)

	existingFile := filepath.Join(tmpDir, "myfile.txt")
	_ = os.WriteFile(existingFile, []byte("hello"), 0644)

	// cd with no args, ~, -
	if v := ValidateCd([]string{"cd"}, tmpDir); v != Free {
		t.Errorf("expected Free for bare cd, got %v", v)
	}
	if v := ValidateCd([]string{"cd", "~"}, tmpDir); v != Free {
		t.Errorf("expected Free for cd ~, got %v", v)
	}
	if v := ValidateCd([]string{"cd", "-"}, tmpDir); v != Free {
		t.Errorf("expected Free for cd -, got %v", v)
	}

	// cd to existing dir
	if v := ValidateCd([]string{"cd", "myfolder"}, tmpDir); v != Valid {
		t.Errorf("expected Valid for cd myfolder, got %v", v)
	}
	if v := ValidateCd([]string{"cd", existingDir}, tmpDir); v != Valid {
		t.Errorf("expected Valid for cd existingDir, got %v", v)
	}

	// cd to existing regular file -> Invalid (not a dir)
	if v := ValidateCd([]string{"cd", "myfile.txt"}, tmpDir); v != Invalid {
		t.Errorf("expected Invalid for cd to regular file, got %v", v)
	}

	// cd to nonexistent dir
	if v := ValidateCd([]string{"cd", "nonexistent"}, tmpDir); v != Invalid {
		t.Errorf("expected Invalid for cd nonexistent, got %v", v)
	}

	// cd with CDPATH set and relative path nonexistent -> Unknown
	t.Setenv("CDPATH", "/some/path")
	if v := ValidateCd([]string{"cd", "somewhere"}, tmpDir); v != Unknown {
		t.Errorf("expected Unknown when CDPATH set and relative dir missing, got %v", v)
	}
}

func TestPath_ExplicitPathsAndInterpreters(t *testing.T) {
	tmpDir := t.TempDir()
	scriptFile := filepath.Join(tmpDir, "run.sh")
	_ = os.WriteFile(scriptFile, []byte("#!/bin/sh\necho hi\n"), 0755)

	pyFile := filepath.Join(tmpDir, "app.py")
	_ = os.WriteFile(pyFile, []byte("print('hi')\n"), 0644)

	// executable in command position
	if v := ValidatePathTokens([]string{"./run.sh"}, tmpDir); v != Valid {
		t.Errorf("expected Valid for ./run.sh, got %v", v)
	}
	if v := ValidatePathTokens([]string{"./missing.sh"}, tmpDir); v != Invalid {
		t.Errorf("expected Invalid for ./missing.sh, got %v", v)
	}

	// interpreter arguments
	if v := ValidatePathTokens([]string{"python", "app.py"}, tmpDir); v != Valid {
		t.Errorf("expected Valid for python app.py, got %v", v)
	}
	if v := ValidatePathTokens([]string{"python3", "missing.py"}, tmpDir); v != Invalid {
		t.Errorf("expected Invalid for python3 missing.py, got %v", v)
	}
	if v := ValidatePathTokens([]string{"sh", "run.sh"}, tmpDir); v != Valid {
		t.Errorf("expected Valid for sh run.sh, got %v", v)
	}

	// interpreter inline code or module flags -> free
	inlineCases := [][]string{
		{"python", "-m", "http.server"},
		{"python", "-c", "import sys"},
		{"node", "-e", "console.log(1)"},
		{"node", "-p", "process.version"},
		{"node", "--eval", "console.log(1)"},
		{"node", "--print", "process.version"},
		{"node", "-r", "ts-node/register", "app.ts"},
		{"node", "--require", "ts-node/register", "app.ts"},
		{"node", "--eval=console.log(1)"},
		{"node", "--print=process.version"},
		{"node", "--require=ts-node/register", "app.ts"},
		{"ruby", "-e", "puts 1"},
		{"bash", "-c", "echo 1"},
	}
	for _, tc := range inlineCases {
		if v := ValidatePathTokens(tc, tmpDir); v != Free {
			t.Errorf("expected Free for %v, got %v", tc, v)
		}
	}

	// other flags preserve script checking
	if v := ValidatePathTokens([]string{"python", "-u", "app.py"}, tmpDir); v != Valid {
		t.Errorf("expected Valid for python -u app.py, got %v", v)
	}
	if v := ValidatePathTokens([]string{"python", "-u", "missing.py"}, tmpDir); v != Invalid {
		t.Errorf("expected Invalid for python -u missing.py, got %v", v)
	}

	// go test ./... has '...' so it is treated as Free
	if v := ValidatePathTokens([]string{"go", "test", "./..."}, tmpDir); v != Free {
		t.Errorf("expected Free for go test ./..., got %v", v)
	}

	// arbitrary command without paths -> Free
	if v := ValidatePathTokens([]string{"git", "status"}, tmpDir); v != Free {
		t.Errorf("expected Free for git status, got %v", v)
	}
}
