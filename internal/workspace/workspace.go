package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type WorkspaceInfo struct {
	HasGit           bool
	GitBranch        string
	HasNodeProject   bool
	HasGoProject     bool
	HasRustProject   bool
	HasPythonProject bool
	HasDockerfile    bool
	HasMakefile      bool
	HasJustfile      bool
	HasK8s           bool
	SignatureFiles   []string
}

var signatureChecks = []struct {
	path  string
	field func(*WorkspaceInfo)
}{
	{".git", func(w *WorkspaceInfo) { w.HasGit = true }},
	{"package.json", func(w *WorkspaceInfo) { w.HasNodeProject = true }},
	{"go.mod", func(w *WorkspaceInfo) { w.HasGoProject = true }},
	{"Cargo.toml", func(w *WorkspaceInfo) { w.HasRustProject = true }},
	{"Dockerfile", func(w *WorkspaceInfo) { w.HasDockerfile = true }},
	{"Makefile", func(w *WorkspaceInfo) { w.HasMakefile = true }},
	{"justfile", func(w *WorkspaceInfo) { w.HasJustfile = true }},
	{"pyproject.toml", func(w *WorkspaceInfo) { w.HasPythonProject = true }},
	{"requirements.txt", func(w *WorkspaceInfo) { w.HasPythonProject = true }},
	{"Chart.yaml", func(w *WorkspaceInfo) { w.HasK8s = true }},
	{"k8s", func(w *WorkspaceInfo) { w.HasK8s = true }},
	{"kubernetes", func(w *WorkspaceInfo) { w.HasK8s = true }},
	{"docker-compose.yml", func(w *WorkspaceInfo) { w.HasDockerfile = true }},
	{"docker-compose.yaml", func(w *WorkspaceInfo) { w.HasDockerfile = true }},
	{"Taskfile.yml", nil},
	{"pom.xml", nil},
	{"build.gradle", nil},
	{"CMakeLists.txt", nil},
}

// Detect scans the given directory for signature files and returns workspace metadata
func Detect(cwd string) WorkspaceInfo {
	var info WorkspaceInfo

	for _, check := range signatureChecks {
		fullPath := filepath.Join(cwd, check.path)
		if _, err := os.Stat(fullPath); err == nil {
			info.SignatureFiles = append(info.SignatureFiles, check.path)
			if check.field != nil {
				check.field(&info)
			}
		}
	}

	info.HasGit, info.GitBranch = detectGitInfo(cwd)

	return info
}

func resolveGitHeadPath(cwd string) (hasGit bool, headPath string) {
	dir := cwd
	for {
		if dir == "" {
			break
		}
		gitPath := filepath.Join(dir, ".git")
		info, err := os.Stat(gitPath)
		if err == nil {
			if info.IsDir() {
				return true, filepath.Join(gitPath, "HEAD")
			}
			content, errRead := os.ReadFile(gitPath)
			if errRead == nil {
				s := strings.TrimSpace(string(content))
				if after, ok := strings.CutPrefix(s, "gitdir: "); ok {
					gitDir := strings.TrimSpace(after)
					if !filepath.IsAbs(gitDir) {
						gitDir = filepath.Join(dir, gitDir)
					}
					return true, filepath.Join(gitDir, "HEAD")
				}
			}
			return true, "" // found .git but couldn't resolve HEAD
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return false, ""
}

func detectGitInfo(cwd string) (hasGit bool, branch string) {
	hasGit, headPath := resolveGitHeadPath(cwd)
	if headPath != "" {
		if data, errHead := os.ReadFile(headPath); errHead == nil {
			s := strings.TrimSpace(string(data))
			if after, ok := strings.CutPrefix(s, "ref: refs/heads/"); ok {
				return hasGit, after
			}
		}
	}
	return hasGit, ""
}

type cacheEntry struct {
	key  string // cwd + "|" + dirModTime
	info WorkspaceInfo
}

var (
	wsCache   *cacheEntry
	wsCacheMu sync.Mutex
)

// DetectCached returns cached workspace info, invalidating when directory modtime changes
// or when the Git HEAD file changes (to catch branch switches that don't affect cwd modtime).
func DetectCached(cwd string) WorkspaceInfo {
	dirInfo, err := os.Stat(cwd)
	if err != nil {
		return Detect(cwd)
	}
	key := cwd + "|" + dirInfo.ModTime().String()

	// Incorporate Git HEAD modtime into cache key to catch branch switches
	if _, headPath := resolveGitHeadPath(cwd); headPath != "" {
		if headInfo, err := os.Stat(headPath); err == nil {
			key += "|HEAD:" + headInfo.ModTime().String()
		}
	}

	wsCacheMu.Lock()
	defer wsCacheMu.Unlock()

	if wsCache != nil && wsCache.key == key {
		return wsCache.info
	}

	info := Detect(cwd)
	wsCache = &cacheEntry{key: key, info: info}
	return info
}
