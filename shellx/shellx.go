// Package shellx provides POSIX shell quoting — the stdlib has no answer to
// "interpolate this string into a shell line safely" (Python ships shlex).
package shellx

import "strings"

// Quote single-quotes s for safe interpolation into a POSIX shell command
// line: inside single quotes nothing expands, and embedded quotes are
// rendered as the standard '\” dance. Safe for any content except NUL
// (which no shell argument can carry anyway).
func Quote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
