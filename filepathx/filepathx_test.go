package filepathx

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDirBytes(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "a"), make([]byte, 100), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "sub/b"), make([]byte, 23), 0o644); err != nil {
		t.Fatal(err)
	}
	// Symlinks are entries, not content — their target size must not count.
	if err := os.Symlink(filepath.Join(dir, "a"), filepath.Join(dir, "link")); err != nil {
		t.Fatal(err)
	}
	if got := DirBytes(dir); got != 123 {
		t.Fatalf("DirBytes = %d, want 123", got)
	}
	if got := DirBytes(filepath.Join(dir, "missing")); got != 0 {
		t.Fatalf("missing dir: %d", got)
	}
}
