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
	for dir != "" {
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

// Normalize cleans and resolves symlinks on path
func Normalize(path string) string {
	if path == "" {
		return ""
	}
	clean := filepath.Clean(path)
	if real, err := filepath.EvalSymlinks(clean); err == nil {
		clean = real
	}
	return clean
}

var projectMarkers = []string{
	"go.mod",
	"package.json",
	"Cargo.toml",
	"justfile",
	"Justfile",
	"Makefile",
	"pyproject.toml",
	"pom.xml",
	"build.gradle",
}

func hasMarker(dir string) bool {
	for _, m := range projectMarkers {
		if _, err := os.Stat(filepath.Join(dir, m)); err == nil {
			return true
		}
	}
	return false
}

// DetectRoot finds the closest repository or project root for cwd.
// Note: WorkspaceInfo (from Detect) inspects signature files strictly in CWD for prompt icons/specs.
// In contrast, DetectRoot traverses upwards to establish scope boundaries for prediction.
func DetectRoot(cwd string) string {
	dir := Normalize(cwd)
	if dir == "" {
		return ""
	}
	home := ""
	if h, err := os.UserHomeDir(); err == nil && h != "" {
		home = Normalize(h)
	}
	var marker string
	for dir != home && dir != filepath.Dir(dir) {
		if _, err := os.Lstat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}
		if marker == "" && hasMarker(dir) {
			marker = dir
		}
		dir = filepath.Dir(dir)
	}
	return marker
}

// ProjectID returns canonical project identifier, resolving git worktrees to main repo
func ProjectID(root string) string {
	if root == "" {
		return ""
	}
	root = Normalize(root)
	gitPath := filepath.Join(root, ".git")
	fi, err := os.Lstat(gitPath)
	if err != nil || fi.IsDir() {
		return root
	}
	b, err := os.ReadFile(gitPath)
	if err != nil {
		return root
	}
	s := strings.TrimSpace(string(b))
	gitdir, ok := strings.CutPrefix(s, "gitdir:")
	if !ok {
		return root
	}
	gitdir = strings.TrimSpace(gitdir)
	if !filepath.IsAbs(gitdir) {
		gitdir = filepath.Join(root, gitdir)
	}
	cd, err := os.ReadFile(filepath.Join(gitdir, "commondir"))
	if err != nil {
		// submodule without commondir
		return root
	}
	common := strings.TrimSpace(string(cd))
	if !filepath.IsAbs(common) {
		common = filepath.Join(gitdir, common)
	}
	common = filepath.Clean(common)
	if filepath.Base(common) != ".git" {
		// bare repo or unexpected layout
		return root
	}
	id := filepath.Dir(common)
	return Normalize(id)
}

var (
	projIDCacheMu sync.RWMutex
	projIDCache   = make(map[string]string)
)

// DetectProjectIDCached returns cached project ID for cwd, avoiding repeated disk stats on keystrokes
func DetectProjectIDCached(cwd string) string {
	if cwd == "" {
		return ""
	}
	projIDCacheMu.RLock()
	id, ok := projIDCache[cwd]
	projIDCacheMu.RUnlock()
	if ok {
		return id
	}

	norm := Normalize(cwd)
	id = ProjectID(DetectRoot(norm))
	projIDCacheMu.Lock()
	projIDCache[cwd] = id
	projIDCacheMu.Unlock()
	return id
}

// InvalidateProjectIDCache clears the cached project IDs
func InvalidateProjectIDCache() {
	projIDCacheMu.Lock()
	projIDCache = make(map[string]string)
	projIDCacheMu.Unlock()
}
