// Package stringsx extends the standard strings package with the log/CSV
// helpers services keep re-writing.
package stringsx

import "strings"

// Truncate returns the first max runes of s (unchanged when it already fits).
// A negative max reads as 0 — a computed budget gone negative should produce
// an empty string, not a panic.
func Truncate(s string, max int) string {
	if max < 0 {
		max = 0
	}
	if max == 0 {
		return ""
	}
	n := 0
	for i := range s {
		if n == max {
			return s[:i]
		}
		n++
	}
	return s
}

// TruncateEllipsis is Truncate plus a trailing "..." so readers can tell a
// cut string from a complete one. The ellipsis is added on top of max — it
// marks truncation rather than tightening the rune budget.
func TruncateEllipsis(s string, max int) string {
	out := Truncate(s, max)
	if out == s {
		return s
	}
	return out + "..."
}

// FirstNonBlank returns the first value that contains any non-whitespace
// character, or "". Use stdlib cmp.Or for plain first-non-empty; this
// variant also skips values that are only spaces — the shape config
// fallback chains actually need.
func FirstNonBlank(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

// SplitAndTrim splits s on sep, trims whitespace from each part and drops
// empties — the tolerant reading of human-written lists ("a, b ,,c").
// Returns nil when nothing survives. sep follows strings.Split semantics
// (an empty sep splits per UTF-8 rune).
func SplitAndTrim(s, sep string) []string {
	if s == "" {
		return nil
	}
	var out []string
	for _, part := range strings.Split(s, sep) {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}
