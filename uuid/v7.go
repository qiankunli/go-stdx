package uuid

import (
	"encoding/hex"

	guuid "github.com/google/uuid"
)

// New returns a time-ordered resource id in "prefix_<v7hex>" form. An empty
// prefix returns the bare V7 hex value.
func New(prefix string) string {
	id := V7Hex()
	if prefix == "" {
		return id
	}
	return prefix + "_" + id
}

// V7 returns a time-ordered (version-7) UUID in canonical dashed form.
// google/uuid keeps ordering even within a single millisecond, so these sort
// by creation time — the property DB keys and sequential resource ids want,
// which V4's uniform randomness scatters. Zero-UUID fallback on entropy
// failure, like V4.
func V7() string {
	u, _ := guuid.NewV7()
	return u.String()
}

// V7Hex returns a version-7 UUID as a 32-character lowercase hex string with
// no dashes — the compact id shape services use for DB keys, pod-name
// suffixes and correlation ids, where the canonical dashed form is noise.
// Same ordering and fallback as V7; always 32 hex characters.
func V7Hex() string {
	u, _ := guuid.NewV7()
	return hex.EncodeToString(u[:])
}
