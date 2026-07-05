package uuid

import (
	"crypto/rand"
	"fmt"
	"io"
	"time"
)

// V7 returns a time-ordered RFC-9562 version-7 UUID: a 48-bit Unix
// millisecond timestamp followed by random bits, so ids sort roughly by
// creation time — the property services want for DB keys and sequential
// resource ids, where V4's uniform randomness scatters index pages.
// Entropy failure degrades like V4 (timestamp fallback, never a panic).
func V7() string {
	var b [16]byte
	// Stamp before the entropy read: a slow rand.Reader must not push the
	// recorded time past the actual call time, or ordering loosens.
	ms := uint64(time.Now().UnixMilli())
	if _, err := io.ReadFull(rand.Reader, b[6:]); err != nil {
		return fmt.Sprintf("fallback-%d", time.Now().UnixNano())
	}
	b[0] = byte(ms >> 40)
	b[1] = byte(ms >> 32)
	b[2] = byte(ms >> 24)
	b[3] = byte(ms >> 16)
	b[4] = byte(ms >> 8)
	b[5] = byte(ms)
	b[6] = (b[6] & 0x0f) | 0x70 // version 7
	b[8] = (b[8] & 0x3f) | 0x80 // variant 1
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
