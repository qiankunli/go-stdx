// Package stringsx extends the standard strings package with the log/CSV
// helpers services keep re-writing.
package stringsx

import "strings"

// Truncate returns the first max bytes of s (unchanged when it already
// fits). Byte semantics, not runes: this is a budget cap for logs and
// storage columns; a multi-byte rune at the boundary may be split.
func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
}

// TruncateEllipsis is Truncate plus a trailing "..." so readers can tell a
// cut string from a complete one. The ellipsis is added on top of max (the
// result may be max+3 bytes) — it marks truncation rather than tightening
// the budget.
func TruncateEllipsis(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
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
// Returns nil when nothing survives.
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
