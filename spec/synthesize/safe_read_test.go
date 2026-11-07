package synthesize

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSafeOpenFile_ValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "sample.txt")
	if err := os.WriteFile(filePath, []byte("hello safe world"), 0644); err != nil {
		t.Fatal(err)
	}

	rc, err := safeOpenFile(tmpDir, "sample.txt", 1024)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if string(data) != "hello safe world" {
		t.Errorf("got %q, want 'hello safe world'", string(data))
	}
}

func TestSafeOpenFile_RejectsRelativeRoot(t *testing.T) {
	_, err := safeOpenFile("relative/path", "file.txt", 1024)
	if err == nil {
		t.Fatal("expected error for relative root path")
	}
}

func TestSafeOpenFile_RejectsTraversal(t *testing.T) {
	tmpDir := t.TempDir()
	_, err := safeOpenFile(tmpDir, "../escape.txt", 1024)
	if err == nil {
		t.Fatal("expected error for path traversal")
	}

	_, err = safeOpenFile(tmpDir, "/etc/passwd", 1024)
	if err == nil {
		t.Fatal("expected error for absolute subpath")
	}
}

func TestSafeOpenFile_SymlinkEscape(t *testing.T) {
	outsideDir := t.TempDir()
	outsideFile := filepath.Join(outsideDir, "secret.txt")
	if err := os.WriteFile(outsideFile, []byte("secret"), 0644); err != nil {
		t.Fatal(err)
	}

	sourceDir := t.TempDir()
	symlinkPath := filepath.Join(sourceDir, "escapelink")
	if err := os.Symlink(outsideFile, symlinkPath); err != nil {
		t.Fatal(err)
	}

	_, err := safeOpenFile(sourceDir, "escapelink", 1024)
	if err == nil {
		t.Fatal("expected error when symlink points outside source root")
	}
	if !strings.Contains(err.Error(), "symlink escape rejected") {
		t.Errorf("expected symlink escape error message, got: %v", err)
	}
}

func TestSafeOpenFile_SymlinkInsideAllowed(t *testing.T) {
	sourceDir := t.TempDir()
	targetFile := filepath.Join(sourceDir, "real.txt")
	if err := os.WriteFile(targetFile, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	linkPath := filepath.Join(sourceDir, "link.txt")
	if err := os.Symlink(targetFile, linkPath); err != nil {
		t.Fatal(err)
	}

	rc, err := safeOpenFile(sourceDir, "link.txt", 1024)
	if err != nil {
		t.Fatalf("expected symlink within source root to succeed: %v", err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "content" {
		t.Errorf("got %q, want 'content'", string(data))
	}
}

func TestSafeOpenFile_FileSizeLimit(t *testing.T) {
	tmpDir := t.TempDir()
	largeFile := filepath.Join(tmpDir, "large.txt")
	content := make([]byte, 2048)
	if err := os.WriteFile(largeFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	_, err := safeOpenFile(tmpDir, "large.txt", 1024)
	if err == nil {
		t.Fatal("expected error when file exceeds max size limit")
	}
	if !strings.Contains(err.Error(), "exceeds limit") {
		t.Errorf("expected size limit error, got: %v", err)
	}
}

func TestSafeOpenFile_AllowsFileStartingWithDotDot(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "..dotdotfile.txt")
	if err := os.WriteFile(filePath, []byte("valid content"), 0644); err != nil {
		t.Fatal(err)
	}

	rc, err := safeOpenFile(tmpDir, "..dotdotfile.txt", 1024)
	if err != nil {
		t.Fatalf("expected file starting with .. to be allowed, got: %v", err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "valid content" {
		t.Errorf("got %q, want 'valid content'", string(data))
	}
}
