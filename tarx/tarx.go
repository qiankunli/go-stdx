// Package tarx packs a directory to a tar.gz stream and safely unpacks one —
// archive/tar gives the container format but leaves the two loops everyone
// re-writes (walk-and-add, extract-without-zip-slip) to the caller, and the
// extraction loop is famously easy to get wrong.
package tarx

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// PackDir writes dir as a tar.gz stream to w. Entries are relative to dir;
// regular files, directories and symlinks are kept (sockets/devices skipped —
// they're runtime state, not data). Entries for which skip returns true are
// omitted (nil = keep everything); a skipped directory prunes its subtree.
func PackDir(dir string, w io.Writer, skip func(rel string) bool) error {
	gz := gzip.NewWriter(w)
	tw := tar.NewWriter(gz)

	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if skip != nil && skip(rel) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		var link string
		if info.Mode()&fs.ModeSymlink != 0 {
			if link, err = os.Readlink(p); err != nil {
				return err
			}
		} else if !info.Mode().IsRegular() && !info.IsDir() {
			return nil // skip sockets, fifos, devices
		}
		hdr, err := tar.FileInfoHeader(info, link)
		if err != nil {
			return err
		}
		hdr.Name = filepath.ToSlash(rel)
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			f, err := os.Open(p)
			if err != nil {
				return err
			}
			_, err = io.Copy(tw, f)
			f.Close()
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if err := tw.Close(); err != nil {
		return err
	}
	return gz.Close()
}

// UnpackDir extracts a tar.gz stream into dir, refusing entries that would
// escape it (zip-slip) — treat any archive that has been at rest elsewhere
// as crossing a trust boundary. Symlinks are confined by what they RESOLVE
// to, not just their text.
func UnpackDir(r io.Reader, dir string) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("tarx: gzip: %w", err)
	}
	defer gz.Close()
	tr := tar.NewReader(gz)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("tarx: tar: %w", err)
		}
		target, err := confine(dir, hdr.Name)
		if err != nil {
			return err
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, fs.FileMode(hdr.Mode)&fs.ModePerm); err != nil {
				return err
			}
		case tar.TypeSymlink:
			// The link target may be relative and escape-y; confine what it
			// resolves to, not just its text.
			resolved := hdr.Linkname
			if !filepath.IsAbs(resolved) {
				resolved = filepath.Join(filepath.Dir(target), resolved)
			}
			if !strings.HasPrefix(filepath.Clean(resolved)+string(os.PathSeparator), filepath.Clean(dir)+string(os.PathSeparator)) {
				return fmt.Errorf("tarx: symlink %q escapes destination", hdr.Name)
			}
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			_ = os.Remove(target)
			if err := os.Symlink(hdr.Linkname, target); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, fs.FileMode(hdr.Mode)&fs.ModePerm)
			if err != nil {
				return err
			}
			_, err = io.Copy(f, tr)
			f.Close()
			if err != nil {
				return err
			}
		default:
			// skip other types
		}
	}
}

// confine joins name under dir and rejects path escapes.
func confine(dir, name string) (string, error) {
	target := filepath.Join(dir, filepath.FromSlash(name))
	if !strings.HasPrefix(target+string(os.PathSeparator), filepath.Clean(dir)+string(os.PathSeparator)) &&
		target != filepath.Clean(dir) {
		return "", fmt.Errorf("tarx: entry %q escapes destination", name)
	}
	return target, nil
}
