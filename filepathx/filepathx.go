// Package filepathx extends path/filepath with the walk-derived measurements
// stdlib leaves as an exercise.
package filepathx

import (
	"io/fs"
	"path/filepath"
)

// DirBytes sums regular-file sizes under dir. Best-effort: unreadable
// entries are skipped rather than failing the whole measurement — callers
// use this for quota/GC accounting where an approximate answer beats an
// error.
func DirBytes(dir string) int64 {
	var n int64
	_ = filepath.WalkDir(dir, func(_ string, d fs.DirEntry, err error) error {
		// Regular files only: a symlink's lstat size (its target string) is
		// not data and would skew quota math.
		if err != nil || !d.Type().IsRegular() {
			return nil
		}
		if info, err := d.Info(); err == nil {
			n += info.Size()
		}
		return nil
	})
	return n
}
