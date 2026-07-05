package uuid

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"time"
)

// v7Bytes builds the 16-byte RFC-9562 version-7 layout: a 48-bit Unix
// millisecond timestamp, the version/variant bits, and random rest — so the
// bytes sort roughly by creation time. On the (practically impossible)
// entropy failure it fills the random span from the nanosecond clock instead
// of bailing, so every shape built on top stays a well-formed v7 value (just
// less unique) rather than a panic or a malformed id.
func v7Bytes() [16]byte {
	var b [16]byte
	// Stamp before the entropy read: a slow rand.Reader must not push the
	// recorded time past the actual call time, or ordering loosens.
	ms := uint64(time.Now().UnixMilli())
	if _, err := io.ReadFull(rand.Reader, b[6:]); err != nil {
		ns := uint64(time.Now().UnixNano())
		for i := 6; i < 16; i++ {
			b[i] = byte(ns >> (8 * (uint(i) % 8)))
		}
	}
	b[0] = byte(ms >> 40)
	b[1] = byte(ms >> 32)
	b[2] = byte(ms >> 24)
	b[3] = byte(ms >> 16)
	b[4] = byte(ms >> 8)
	b[5] = byte(ms)
	b[6] = (b[6] & 0x0f) | 0x70 // version 7
	b[8] = (b[8] & 0x3f) | 0x80 // variant 1
	return b
}

// V7 returns a time-ordered RFC-9562 version-7 UUID in canonical dashed form:
// a 48-bit Unix millisecond timestamp followed by random bits, so ids sort
// roughly by creation time — the property services want for DB keys and
// sequential resource ids, where V4's uniform randomness scatters index
// pages. Entropy failure degrades like V4 (never a panic).
func V7() string {
	b := v7Bytes()
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// V7Hex returns a version-7 UUID as a 32-character lowercase hex string with
// no dashes — the compact id shape services use for DB keys, pod-name
// suffixes and correlation ids, where the canonical dashed form is noise.
// Same time-ordering and entropy-failure behavior as V7; the output is
// always 32 hex characters, so callers can rely on the width.
func V7Hex() string {
	b := v7Bytes()
	return hex.EncodeToString(b[:])
}
