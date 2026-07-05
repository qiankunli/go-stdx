package tarx

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"
)

func TestPackUnpackRoundTrip(t *testing.T) {
	src := t.TempDir()
	if err := os.MkdirAll(filepath.Join(src, "sub/deep"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "a.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "sub/deep/b.bin"), []byte{0, 1, 2}, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("a.txt", filepath.Join(src, "link")); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := PackDir(src, &buf, nil); err != nil {
		t.Fatalf("pack: %v", err)
	}

	dst := t.TempDir()
	if err := UnpackDir(&buf, dst); err != nil {
		t.Fatalf("unpack: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dst, "a.txt"))
	if err != nil || string(data) != "hello" {
		t.Fatalf("a.txt: %q err=%v", data, err)
	}
	fi, err := os.Stat(filepath.Join(dst, "sub/deep/b.bin"))
	if err != nil || fi.Mode().Perm() != 0o755 {
		t.Fatalf("b.bin mode: %v err=%v", fi.Mode(), err)
	}
	// Symlink survives and resolves inside the destination.
	got, err := os.Readlink(filepath.Join(dst, "link"))
	if err != nil || got != "a.txt" {
		t.Fatalf("link: %q err=%v", got, err)
	}
}

func TestPackDirSkip(t *testing.T) {
	src := t.TempDir()
	if err := os.MkdirAll(filepath.Join(src, "keep"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(src, "drop.local"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "keep/a"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "drop.local/b"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err := PackDir(src, &buf, func(rel string) bool { return rel == "drop.local" })
	if err != nil {
		t.Fatalf("pack: %v", err)
	}
	dst := t.TempDir()
	if err := UnpackDir(&buf, dst); err != nil {
		t.Fatalf("unpack: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dst, "keep/a")); err != nil {
		t.Fatalf("kept file missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dst, "drop.local")); !os.IsNotExist(err) {
		t.Fatalf("skipped dir leaked: %v", err)
	}
}

// mkTar builds a raw tar.gz with fully controlled entries (attack payloads).
func mkTar(t *testing.T, entries []tar.Header) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for i := range entries {
		if err := tw.WriteHeader(&entries[i]); err != nil {
			t.Fatal(err)
		}
	}
	tw.Close()
	gz.Close()
	return &buf
}

func TestUnpackRejectsEscapes(t *testing.T) {
	dst := t.TempDir()

	// Path traversal in the entry name.
	slip := mkTar(t, []tar.Header{{Name: "../evil.txt", Typeflag: tar.TypeReg, Mode: 0o644}})
	if err := UnpackDir(slip, dst); err == nil {
		t.Fatal("path traversal not rejected")
	}

	// Symlink whose target escapes the destination.
	link := mkTar(t, []tar.Header{{Name: "l", Typeflag: tar.TypeSymlink, Linkname: "../../etc/passwd"}})
	if err := UnpackDir(link, dst); err == nil {
		t.Fatal("escaping symlink not rejected")
	}
}
