package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
// NOTE: it will skip hidden file start with dot prefix
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
			dir = filepath.Join(base, partial[:i+1])
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

			// skip hidden files
			if strings.HasPrefix(name, ".") {
				continue
			}

			// filter by filePrefix
			if filePrefix != "" && !hasPrefix(name, filePrefix) {
				continue
			}

			fullPath := pathPrefix + name
			if entry.IsDir() {
				if dirOnly || len(filterSet) == 0 {
					results = append(results, Suggestion{
						Cmd:  fullPath + "/",
						Desc: "directory",
					})
				}
				continue
			}
			// if filters are set, only show matching extensions
			if len(filterSet) > 0 {
				ext := strings.ToLower(filepath.Ext(name))
				if !filterSet[ext] {
					continue
				}
			}
			results = append(results, Suggestion{
				Cmd:  fullPath,
				Desc: "file",
			})
		}

		return results
	}
}
