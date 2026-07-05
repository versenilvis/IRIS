package spec

import (
	"os"
	"path/filepath"
	"sync"
)

var (
	pathCmds     = make(map[string]bool)
	pathCmdsOnce sync.Once
)

// scanExternalCommands populates pathCmds with system executables
func scanExternalCommands() {
	pathCmdsOnce.Do(func() {
		scanPath()
	})
}

// scanPath populates pathCmds with all executable files found in path environment variable
func scanPath() {
	dirs := filepath.SplitList(os.Getenv("PATH"))
	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, f := range files {
			if !f.IsDir() {
				info, err := f.Info()
				if err == nil && info.Mode()&0111 != 0 {
					pathCmds[f.Name()] = true
				}
			}
		}
	}
}
