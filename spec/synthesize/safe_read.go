package synthesize

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const defaultMaxFileSize = 2 * 1024 * 1024

type limitReadCloser struct {
	io.Reader
	closer io.Closer
}

func (l *limitReadCloser) Close() error {
	if l.closer != nil {
		return l.closer.Close()
	}
	return nil
}

func safeOpenFile(rootDir, relPath string, maxSize int64) (io.ReadCloser, error) {
	// prevent relative resolution from current working directory
	if rootDir == "" || !filepath.IsAbs(rootDir) {
		return nil, errors.New("source root directory must be an absolute path")
	}
	if maxSize <= 0 {
		maxSize = defaultMaxFileSize
	}

	cleanRel := filepath.Clean(relPath)
	sep := string(filepath.Separator)
	if cleanRel != relPath || cleanRel == ".." || strings.HasPrefix(cleanRel, ".."+sep) || filepath.IsAbs(cleanRel) {
		return nil, fmt.Errorf("invalid relative path %q", relPath)
	}

	canonicalRoot, err := filepath.EvalSymlinks(rootDir)
	if err != nil {
		return nil, fmt.Errorf("resolve root directory: %w", err)
	}

	fullPath := filepath.Join(canonicalRoot, cleanRel)
	canonicalTarget, err := filepath.EvalSymlinks(fullPath)
	if err != nil {
		return nil, fmt.Errorf("resolve target file: %w", err)
	}

	// ensure symlinks cannot escape the designated root
	relToRoot, err := filepath.Rel(canonicalRoot, canonicalTarget)
	if err != nil || relToRoot == ".." || strings.HasPrefix(relToRoot, ".."+sep) || filepath.IsAbs(relToRoot) {
		return nil, fmt.Errorf("symlink escape rejected: %q escapes source directory %q", fullPath, canonicalRoot)
	}

	info, err := os.Stat(canonicalTarget)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%q is not a regular file", canonicalTarget)
	}
	if info.Size() > maxSize {
		return nil, fmt.Errorf("file size %d exceeds limit of %d bytes", info.Size(), maxSize)
	}

	f, err := os.Open(canonicalTarget)
	if err != nil {
		return nil, err
	}

	return &limitReadCloser{
		Reader: io.LimitReader(f, maxSize),
		closer: f,
	}, nil
}
