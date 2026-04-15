package core

import (
	"os"
	"path/filepath"
)

var pathCmds = make(map[string]bool)

// scanExternalCommands populates pathCmds with system executables
func scanExternalCommands() {
	scanPath()
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
