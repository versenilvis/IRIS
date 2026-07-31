package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileGenerator(t *testing.T) {
	SetCWD("")

	// Setup mock files
	tmp := t.TempDir()
	_ = os.MkdirAll(filepath.Join(tmp, "src"), 0755)
	_ = os.WriteFile(filepath.Join(tmp, "main.go"), []byte(""), 0644)
	_ = os.WriteFile(filepath.Join(tmp, "README.md"), []byte(""), 0644)
	_ = os.WriteFile(filepath.Join(tmp, ".hidden"), []byte(""), 0644)
	_ = os.WriteFile(filepath.Join(tmp, "src/utils.go"), []byte(""), 0644)

	oldWd, _ := os.Getwd()
	_ = os.Chdir(tmp)
	defer func() { _ = os.Chdir(oldWd) }()


	t.Run("dirOnly shows only dirs", func(t *testing.T) {
		gen := FileGenerator("/")
		results := gen([]string{"cd", ""}, "cd ", "")
		for _, r := range results {
			if !strings.HasSuffix(r.Cmd, "/") {
				t.Errorf("Expected only directories, got %q", r.Cmd)
			}
		}
	})


	t.Run("Filter extension", func(t *testing.T) {
		gen := FileGenerator(".go")
		results := gen([]string{"ls", ""}, "ls ", "")
		foundMain := false
		for _, r := range results {
			if r.Cmd == "main.go" {
				foundMain = true
			}
			if r.Cmd == "README.md" {
				t.Errorf("Did not expect README.md when filtering for .go")
			}
		}
		if !foundMain {
			t.Errorf("FileGenerator(\".go\") did not suggest main.go")
		}
	})


	t.Run("Nested path", func(t *testing.T) {
		gen := FileGenerator()
		results := gen([]string{"ls", "src/u"}, "ls src/u", "src/u")
		foundUtils := false
		for _, r := range results {
			if r.Cmd == "src/utils.go" {
				foundUtils = true
			}
		}
		if !foundUtils {
			t.Errorf("Did not find src/utils.go for nested path src/u")
		}
	})


	t.Run("Deep scan 1 level", func(t *testing.T) {
		gen := FileGenerator()
		results := gen([]string{"ls", "src/"}, "ls src/", "src/")
		foundUtils := false
		for _, r := range results {
			if r.Cmd == "src/utils.go" {
				foundUtils = true
			}
		}
		if !foundUtils {
			t.Errorf("Deep scan did not find src/utils.go")
		}
	})




	t.Run("Hidden files are skipped", func(t *testing.T) {
		gen := FileGenerator()
		results := gen([]string{"ls", ""}, "ls ", "")
		for _, r := range results {
			if strings.HasPrefix(r.Cmd, ".") {
				t.Errorf("Hidden file %q should be skipped", r.Cmd)
			}
		}
	})
}

func TestGetCWDUsesShellReportedDirectory(t *testing.T) {
	launcherDir := t.TempDir()
	shellDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(shellDir, "from-shell.txt"), nil, 0644); err != nil {
		t.Fatal(err)
	}

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(launcherDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		SetCWD("")
		_ = os.Chdir(oldWd)
	})

	SetCWD(shellDir)
	if got := GetCWD(); got != shellDir {
		t.Fatalf("GetCWD() = %q, want shell-reported directory %q", got, shellDir)
	}

	results := FileGenerator()([]string{"cat", ""}, "cat ", "")
	found := false
	for _, result := range results {
		if result.Cmd == "from-shell.txt" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("FileGenerator() did not read suggestions from shell directory %q", shellDir)
	}
}

func TestSetCWDRejectsRelativePaths(t *testing.T) {
	SetCWD("")
	t.Cleanup(func() { SetCWD("") })

	SetCWD("relative/path")
	if got := GetCWD(); got == "relative/path" {
		t.Fatalf("GetCWD() accepted relative shell directory %q", got)
	}
}
