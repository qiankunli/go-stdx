package osx

import (
	"os"
	"path/filepath"
)

// WriteFileAtomic writes data to path via a same-directory temp file and
// rename, so readers never observe a truncated or half-written file — the
// guarantee os.WriteFile does not give when the writer crashes mid-write.
// The temp file lives next to path (rename must not cross filesystems).
// No fsync: this guards against partial writes, not power loss; callers
// needing power-loss durability should sync the file and directory themselves.
func WriteFileAtomic(path string, data []byte, perm os.FileMode) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(name)
		return err
	}
	// CreateTemp opens 0600; apply the caller's mode before the file becomes
	// visible under its real name.
	if err := tmp.Chmod(perm); err != nil {
		tmp.Close()
		os.Remove(name)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(name)
		return err
	}
	if err := os.Rename(name, path); err != nil {
		os.Remove(name)
		return err
	}
	return nil
}
