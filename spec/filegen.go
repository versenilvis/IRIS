package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/versenilvis/iris/internal/config"
)

// for tracking current working dir
var ShellPID int

// GetCWD returns the current working directory of the underlying shell
// on Linux, it reads it from /proc/[pid]/cwd
func GetCWD() string {
	if ShellPID > 0 {
		path := fmt.Sprintf("/proc/%d/cwd", ShellPID)
		cwd, err := os.Readlink(path)
		if err == nil {
			return cwd
		}
	}
	cwd, _ := os.Getwd()
	return cwd
}

// FileGenerator provides directory and file suggestions.
// It handles nested paths (e.g., 'src/main') by resolving the directory
// and filtering based on the last path component
//
// E.g, we have both src/main.go and src/main.py
// you type "cat src/mai", partial = "src/mai"
// dir = "src/", filePrefix = "mai"
// search for which matches "mai" case, and show it (whici is both main.go and main.py)
// but if you type "go run src/mai" -> it will only shows suggestion about src/main.go
// NOTE: it will skip hidden file start with dot prefix unless configured otherwise
func FileGenerator(filters ...string) GeneratorFunc {
	dirOnly := false
	filterSet := make(map[string]bool)
	for _, f := range filters {
		if f == "/" {
			dirOnly = true
			continue
		}
		filterSet[strings.ToLower(f)] = true
	}

	return func(tokens []string, prefix string, partial string) []Suggestion {
		base := GetCWD()
		dir := base
		filePrefix := partial

		// check if the partial string contains path separators
		if i := strings.LastIndexAny(partial, "/\\"); i != -1 {
			pathDir := partial[:i+1]
			if filepath.IsAbs(pathDir) || strings.HasPrefix(pathDir, "~") {
				if strings.HasPrefix(pathDir, "~") {
					home, _ := os.UserHomeDir()
					pathDir = filepath.Join(home, pathDir[1:])
				}
				dir = pathDir
			} else {
				dir = filepath.Join(base, pathDir)
			}
			filePrefix = partial[i+1:]
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil
		}

		// build the path prefix for the relative suggestions
		// e.g. if partial is "src/m", then pathPrefix is "src/"
		pathPrefix := ""
		if i := strings.LastIndexAny(partial, "/\\"); i != -1 {
			pathPrefix = partial[:i+1]
		}

		var results []Suggestion
		for _, entry := range entries {
			name := entry.Name()

			// skip hidden files if not configured to show them
			if !config.Get().UI.ShowHiddenFiles && strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			match := false
			if filePrefix == "" {
				match = true
			} else if dirOnly {
				match = strings.Contains(strings.ToLower(name), strings.ToLower(filePrefix))
			} else {
				match = HasPrefix(name, filePrefix)
			}

			if !match {
				continue
			}

			fullPath := pathPrefix + name

			isDir := entry.IsDir()
			if !isDir && entry.Type()&os.ModeSymlink != 0 {
				if info, err := os.Stat(filepath.Join(dir, name)); err == nil && info.IsDir() {
					isDir = true
				}
			}

			if isDir {
				if dirOnly || len(filterSet) == 0 {
					results = append(results, Suggestion{
						Cmd:      fullPath + "/",
						Desc:     "directory",
						Priority: 50,
					})
				} else {
					// scan only 1 level deeper if there is a filter
					subEntries, err := os.ReadDir(filepath.Join(dir, name))
					if err == nil {
						for _, subEntry := range subEntries {
							subIsDir := subEntry.IsDir()
							if !subIsDir && subEntry.Type()&os.ModeSymlink != 0 {
								if info, err := os.Stat(filepath.Join(dir, name, subEntry.Name())); err == nil && info.IsDir() {
									subIsDir = true
								}
							}
							if subIsDir {
								continue
							}
							subName := subEntry.Name()
							if strings.HasPrefix(subName, ".") {
								continue
							}
							ext := strings.ToLower(filepath.Ext(subName))
							if filterSet[ext] {
								results = append(results, Suggestion{
									Cmd:      fullPath + "/" + subName,
									Desc:     "file",
									Priority: 50,
								})
							}
						}
					}
				}
				continue
			}

			if dirOnly {
				continue
			}

			// if filters are set, only show matching extensions
			if len(filterSet) > 0 {
				ext := strings.ToLower(filepath.Ext(name))
				if !filterSet[ext] {
					continue
				}
			}
			desc := "file"
			if ext := strings.ToLower(filepath.Ext(name)); ext != "" {
				desc = strings.TrimPrefix(ext, ".")
			}
			results = append(results, Suggestion{
				Cmd:      fullPath,
				Desc:     desc,
				Priority: 50,
			})
		}

		sort.Slice(results, func(i, j int) bool {
			return strings.ToLower(results[i].Cmd) < strings.ToLower(results[j].Cmd)
		})

		return results
	}
}
