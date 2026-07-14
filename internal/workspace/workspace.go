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

	info.GitBranch = detectGitBranch(cwd)
	if info.GitBranch != "" {
		info.HasGit = true
	}

	return info
}

func detectGitBranch(cwd string) string {
	dir := cwd
	for i := 0; i < 10 && dir != "" && dir != "/" && dir != "."; i++ {
		gitPath := filepath.Join(dir, ".git")
		info, err := os.Stat(gitPath)
		if err == nil {
			var headPath string
			if info.IsDir() {
				headPath = filepath.Join(gitPath, "HEAD")
			} else {
				content, errRead := os.ReadFile(gitPath)
				if errRead == nil {
					s := strings.TrimSpace(string(content))
					if strings.HasPrefix(s, "gitdir: ") {
						gitDir := strings.TrimSpace(strings.TrimPrefix(s, "gitdir: "))
						if !filepath.IsAbs(gitDir) {
							gitDir = filepath.Join(dir, gitDir)
						}
						headPath = filepath.Join(gitDir, "HEAD")
					}
				}
			}
			if headPath != "" {
				if data, errHead := os.ReadFile(headPath); errHead == nil {
					s := strings.TrimSpace(string(data))
					if strings.HasPrefix(s, "ref: refs/heads/") {
						return strings.TrimPrefix(s, "ref: refs/heads/")
					}
				}
			}
			return ""
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
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
// this handles mid-session file creation (e.g. go mod init) without requiring cd
func DetectCached(cwd string) WorkspaceInfo {
	dirInfo, err := os.Stat(cwd)
	if err != nil {
		return Detect(cwd)
	}
	key := cwd + "|" + dirInfo.ModTime().String()

	wsCacheMu.Lock()
	defer wsCacheMu.Unlock()

	if wsCache != nil && wsCache.key == key {
		return wsCache.info
	}

	info := Detect(cwd)
	wsCache = &cacheEntry{key: key, info: info}
	return info
}
