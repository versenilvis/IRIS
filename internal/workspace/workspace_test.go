package workspace

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDetect_GitAndGoProject(t *testing.T) {
	tmp := t.TempDir()

	_ = os.Mkdir(filepath.Join(tmp, ".git"), 0755)
	_ = os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module test"), 0644)

	info := Detect(tmp)

	if !info.HasGit {
		t.Error("expected HasGit to be true")
	}
	if !info.HasGoProject {
		t.Error("expected HasGoProject to be true")
	}
	if info.HasNodeProject {
		t.Error("expected HasNodeProject to be false")
	}
	if info.HasRustProject {
		t.Error("expected HasRustProject to be false")
	}
	if len(info.SignatureFiles) != 2 {
		t.Errorf("expected 2 signature files, got %d: %v", len(info.SignatureFiles), info.SignatureFiles)
	}
}

func TestDetect_EmptyDirectory(t *testing.T) {
	tmp := t.TempDir()
	info := Detect(tmp)

	if info.HasGit || info.HasNodeProject || info.HasGoProject || info.HasRustProject || info.HasDockerfile || info.HasMakefile {
		t.Error("expected all flags to be false for empty directory")
	}
	if len(info.SignatureFiles) != 0 {
		t.Errorf("expected 0 signature files, got %d", len(info.SignatureFiles))
	}
}

func TestDetect_NodeProject(t *testing.T) {
	tmp := t.TempDir()
	_ = os.WriteFile(filepath.Join(tmp, "package.json"), []byte("{}"), 0644)
	_ = os.WriteFile(filepath.Join(tmp, "Dockerfile"), []byte("FROM node"), 0644)

	info := Detect(tmp)

	if !info.HasNodeProject {
		t.Error("expected HasNodeProject to be true")
	}
	if !info.HasDockerfile {
		t.Error("expected HasDockerfile to be true")
	}
	if info.HasGit {
		t.Error("expected HasGit to be false")
	}
}

func TestDetectCached_MidSessionFileCreation(t *testing.T) {
	tmp := t.TempDir()

	// first call: no go.mod exists
	info1 := DetectCached(tmp)
	if info1.HasGoProject {
		t.Fatal("expected HasGoProject to be false before creating go.mod")
	}

	// create go.mod mid-session (same cwd, no cd)
	_ = os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module test"), 0644)

	// second call: cache should invalidate because directory modtime changed
	info2 := DetectCached(tmp)
	if !info2.HasGoProject {
		t.Fatal("expected HasGoProject to be true after creating go.mod in same cwd")
	}
}

func TestDetectCached_CwdChange(t *testing.T) {
	dir1 := t.TempDir()
	dir2 := t.TempDir()

	_ = os.Mkdir(filepath.Join(dir1, ".git"), 0755)

	info1 := DetectCached(dir1)
	if !info1.HasGit {
		t.Fatal("expected HasGit for dir1")
	}

	info2 := DetectCached(dir2)
	if info2.HasGit {
		t.Fatal("expected no HasGit for dir2")
	}
}

func TestDetect_MultiEcosystems(t *testing.T) {
	tmp := t.TempDir()
	_ = os.WriteFile(filepath.Join(tmp, "justfile"), []byte("build:"), 0644)
	_ = os.WriteFile(filepath.Join(tmp, "pyproject.toml"), []byte(""), 0644)
	_ = os.WriteFile(filepath.Join(tmp, "Chart.yaml"), []byte("apiVersion: v2"), 0644)

	info := Detect(tmp)

	if !info.HasJustfile {
		t.Error("expected HasJustfile to be true")
	}
	if !info.HasPythonProject {
		t.Error("expected HasPythonProject to be true")
	}
	if !info.HasK8s {
		t.Error("expected HasK8s to be true")
	}
}

func TestDetectCached_BranchSwitchWithoutDirChange(t *testing.T) {
	tmp := t.TempDir()
	gitDir := filepath.Join(tmp, ".git")
	_ = os.Mkdir(gitDir, 0755)

	headPath := filepath.Join(gitDir, "HEAD")
	_ = os.WriteFile(headPath, []byte("ref: refs/heads/main\n"), 0644)

	dirInfoBefore, _ := os.Stat(tmp)

	info1 := DetectCached(tmp)
	if info1.GitBranch != "main" {
		t.Fatalf("expected branch 'main', got %q", info1.GitBranch)
	}

	// Simulate branch switch by updating .git/HEAD only
	// (this does not update the modtime of tmp on typical filesystems since tmp's direct children didn't change)
	_ = os.WriteFile(headPath, []byte("ref: refs/heads/feature\n"), 0644)
	
	// Force the modtime of HEAD to be distinct to avoid flakiness on low-res file systems
	infoAfter, _ := os.Stat(headPath)
	newMod := infoAfter.ModTime().Add(2 * time.Second)
	_ = os.Chtimes(headPath, newMod, newMod)

	dirInfoAfter, _ := os.Stat(tmp)
	if dirInfoBefore.ModTime() != dirInfoAfter.ModTime() {
		t.Log("Note: directory modtime changed automatically on this filesystem")
	}

	info2 := DetectCached(tmp)
	if info2.GitBranch != "feature" {
		t.Fatalf("expected branch 'feature', got %q", info2.GitBranch)
	}
}
