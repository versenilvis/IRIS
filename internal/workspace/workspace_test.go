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

func TestDetectRoot_Monorepo(t *testing.T) {
	tmp := Normalize(t.TempDir())
	repoRoot := filepath.Join(tmp, "my-repo")
	backendDir := filepath.Join(repoRoot, "backend", "cmd")
	_ = os.MkdirAll(backendDir, 0755)
	_ = os.Mkdir(filepath.Join(repoRoot, ".git"), 0755)
	_ = os.WriteFile(filepath.Join(repoRoot, "backend", "go.mod"), []byte("module backend"), 0644)

	// .git at repoRoot must take precedence over inner go.mod
	got := DetectRoot(backendDir)
	if got != repoRoot {
		t.Fatalf("DetectRoot(%q) = %q, want %q", backendDir, got, repoRoot)
	}
}

func TestDetectRoot_NoGitFallbackToMarker(t *testing.T) {
	tmp := Normalize(t.TempDir())
	projDir := filepath.Join(tmp, "standalone-project")
	subDir := filepath.Join(projDir, "src", "pkg")
	_ = os.MkdirAll(subDir, 0755)
	_ = os.WriteFile(filepath.Join(projDir, "package.json"), []byte("{}"), 0644)

	got := DetectRoot(subDir)
	if got != projDir {
		t.Fatalf("DetectRoot(%q) = %q, want %q", subDir, got, projDir)
	}
}

func TestDetectRoot_NeverHomeOrRoot(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("no home dir")
	}
	home = Normalize(home)

	// cwd at home must never return home
	if got := DetectRoot(home); got != "" {
		t.Fatalf("DetectRoot(home) = %q, want empty", got)
	}

	// cwd at root / must never return root
	if got := DetectRoot("/"); got != "" {
		t.Fatalf("DetectRoot('/') = %q, want empty", got)
	}
}

func TestDetectRoot_Symlink(t *testing.T) {
	tmp := Normalize(t.TempDir())
	realRepo := filepath.Join(tmp, "real-repo")
	_ = os.MkdirAll(filepath.Join(realRepo, "sub"), 0755)
	_ = os.Mkdir(filepath.Join(realRepo, ".git"), 0755)

	symlinkPath := filepath.Join(tmp, "symlink-repo")
	if err := os.Symlink(realRepo, symlinkPath); err != nil {
		t.Skip("symlink not supported")
	}

	got := DetectRoot(filepath.Join(symlinkPath, "sub"))
	if got != realRepo {
		t.Fatalf("DetectRoot(symlink/sub) = %q, want %q", got, realRepo)
	}
}

func TestProjectID_Worktree(t *testing.T) {
	tmp := Normalize(t.TempDir())
	mainRepo := filepath.Join(tmp, "main-repo")
	mainGit := filepath.Join(mainRepo, ".git")
	_ = os.MkdirAll(filepath.Join(mainGit, "worktrees", "wt1"), 0755)

	// create worktree directory
	wtDir := filepath.Join(tmp, "wt-branch")
	_ = os.MkdirAll(wtDir, 0755)

	// worktree .git file
	gitdir := filepath.Join(mainGit, "worktrees", "wt1")
	_ = os.WriteFile(filepath.Join(wtDir, ".git"), []byte("gitdir: "+gitdir+"\n"), 0644)
	// commondir inside gitdir pointing to ../..
	_ = os.WriteFile(filepath.Join(gitdir, "commondir"), []byte("../..\n"), 0644)

	id := ProjectID(wtDir)
	if id != mainRepo {
		t.Fatalf("ProjectID(worktree) = %q, want %q", id, mainRepo)
	}
}

func TestProjectID_Submodule(t *testing.T) {
	tmp := Normalize(t.TempDir())
	submoduleDir := filepath.Join(tmp, "my-submodule")
	_ = os.MkdirAll(submoduleDir, 0755)

	// submodule has .git file pointing to module gitdir without commondir
	fakeGitDir := filepath.Join(tmp, "main", ".git", "modules", "subm")
	_ = os.MkdirAll(fakeGitDir, 0755)
	_ = os.WriteFile(filepath.Join(submoduleDir, ".git"), []byte("gitdir: "+fakeGitDir+"\n"), 0644)

	id := ProjectID(submoduleDir)
	if id != submoduleDir {
		t.Fatalf("ProjectID(submodule) = %q, want %q", id, submoduleDir)
	}
}

func TestDetectRoot_NeverHomeOrRoot_WithEnv(t *testing.T) {
	fakeHome := Normalize(t.TempDir())
	t.Setenv("HOME", fakeHome)

	// dotfiles repo directly at $HOME
	_ = os.Mkdir(filepath.Join(fakeHome, ".git"), 0755)
	_ = os.WriteFile(filepath.Join(fakeHome, "package.json"), []byte("{}"), 0644)

	downloads := filepath.Join(fakeHome, "Downloads")
	_ = os.MkdirAll(downloads, 0755)

	// cwd at $HOME must return "" even with .git and package.json
	if got := DetectRoot(fakeHome); got != "" {
		t.Fatalf("DetectRoot(fakeHome) = %q, want empty", got)
	}

	// cwd in non-project subdir under $HOME must return ""
	if got := DetectRoot(downloads); got != "" {
		t.Fatalf("DetectRoot(fakeHome/Downloads) = %q, want empty", got)
	}

	// project under $HOME should still resolve correctly
	proj := filepath.Join(fakeHome, "projects", "iris")
	_ = os.MkdirAll(filepath.Join(proj, ".git"), 0755)
	if got := DetectRoot(proj); got != proj {
		t.Fatalf("DetectRoot(proj) = %q, want %q", got, proj)
	}
}

func TestDetectRoot_Case2_MonorepoRegression(t *testing.T) {
	tmp := Normalize(t.TempDir())
	repo := filepath.Join(tmp, "repo")
	backend := filepath.Join(repo, "backend")
	_ = os.MkdirAll(backend, 0755)
	_ = os.Mkdir(filepath.Join(repo, ".git"), 0755)
	_ = os.WriteFile(filepath.Join(backend, "go.mod"), []byte("module backend"), 0644)

	rootBackend := DetectRoot(backend)
	rootRepo := DetectRoot(repo)

	if rootBackend != repo || rootRepo != repo {
		t.Fatalf("Case 2 broken: DetectRoot(backend)=%q, DetectRoot(repo)=%q, want both %q", rootBackend, rootRepo, repo)
	}
}

func TestDetectProjectIDCached(t *testing.T) {
	tmp := Normalize(t.TempDir())
	repo := filepath.Join(tmp, "cached-repo")
	_ = os.MkdirAll(filepath.Join(repo, "sub"), 0755)
	_ = os.Mkdir(filepath.Join(repo, ".git"), 0755)

	id1 := DetectProjectIDCached(filepath.Join(repo, "sub"))
	id2 := DetectProjectIDCached(filepath.Join(repo, "sub"))

	if id1 != repo || id2 != repo {
		t.Fatalf("DetectProjectIDCached = %q, %q, want %q", id1, id2, repo)
	}
}
