package core

import (
	"os"
	"path/filepath"
	"strings"
)

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
		dir := "."
		filePrefix := partial

		// check if the partial string contains path separators
		if i := strings.LastIndexAny(partial, "/\\"); i != -1 {
			dir = partial[:i+1]
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
				results = append(results, Suggestion{
					Cmd:  prefix + " " + fullPath + "/",
					Desc: "directory",
				})
			} else {
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
				results = append(results, Suggestion{
					Cmd:  prefix + " " + fullPath,
					Desc: fileDesc(name),
				})
			}
		}

		return results
	}
}

func fileDesc(name string) string {
	ext := strings.ToLower(filepath.Ext(name))

	// special cases without extensions
	if ext == "" {
		switch name {
		case "Dockerfile", "dockerfile":
			return "dockerfile"
		case "Makefile", "makefile":
			return "makefile"
		case "LICENSE", "license":
			return "license"
		}
	}

	switch ext {
	case ".go":
		return "go"
	case ".js":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".py":
		return "python"
	case ".rs":
		return "rust"
	case ".c", ".h":
		return "c source"
	case ".cpp", ".hpp", ".cc":
		return "c++ source"
	case ".java":
		return "java"
	case ".rb":
		return "ruby"
	case ".php":
		return "php"
	case ".kt":
		return "kotlin"
	case ".swift":
		return "swift"
	case ".sh":
		return "shell script"
	case ".lua":
		return "lua"
	case ".zig":
		return "zig"
	case ".ex", ".exs":
		return "elixir"

	case ".html", ".htm":
		return "html"
	case ".css":
		return "stylesheet"
	case ".scss", ".sass":
		return "sass"
	case ".jsx":
		return "react (js)"
	case ".tsx":
		return "react (ts)"
	case ".vue":
		return "vue"
	case ".svelte":
		return "svelte"

	case ".json":
		return "json"
	case ".yaml", ".yml":
		return "yaml"
	case ".toml":
		return "toml"
	case ".sql":
		return "sql"
	case ".xml":
		return "xml"
	case ".csv":
		return "csv"
	case ".env":
		return "env file"
	case ".lock":
		return "lock file"

	case ".md", ".markdown":
		return "markdown"
	case ".txt":
		return "text file"
	case ".pdf":
		return "pdf"
	case ".log":
		return "log file"

	case ".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp", ".ico":
		return "image"
	case ".mp4", ".mkv", ".avi", ".mov", ".webm":
		return "video"
	case ".mp3", ".wav", ".flac", ".ogg":
		return "audio"

	case ".dockerfile":
		return "dockerfile"
	case ".tf":
		return "terraform"
	case ".proto":
		return "protobuf"

	case ".zip", ".tar", ".gz", ".bz2", ".xz", ".7z", ".rar":
		return "archive"

	case ".mod":
		return "go module"
	case ".sum":
		return "go checksum"

	default:
		return "file"
	}
}
