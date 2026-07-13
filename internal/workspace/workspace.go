package workspace

import (
	"os"
	"path/filepath"
	"sync"
)

type WorkspaceInfo struct {
	HasGit         bool
	HasNodeProject bool
	HasGoProject   bool
	HasRustProject bool
	HasDockerfile  bool
	HasMakefile    bool
	SignatureFiles []string
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
	{"docker-compose.yml", nil},
	{"docker-compose.yaml", nil},
	{"justfile", nil},
	{"Taskfile.yml", nil},
	{"pom.xml", nil},
	{"build.gradle", nil},
	{"CMakeLists.txt", nil},
	{"pyproject.toml", nil},
	{"requirements.txt", nil},
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

	return info
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
