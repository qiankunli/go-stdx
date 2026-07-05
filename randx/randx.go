// Package randx extends crypto/rand with the id shapes daemons actually
// mint: short random hex strings for correlation/resource ids where a full
// UUID is more ceremony than needed.
package randx

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"strconv"
	"time"
)

// Hex returns n random bytes as a 2n-character lowercase hex string.
// Like uuid.V4, entropy failure degrades to a clock-derived id instead of
// panicking — these are resource ids, where a collision is a nuisance and
// a crash is not an acceptable trade. n <= 0 returns "".
func Hex(n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		s := strconv.FormatUint(uint64(time.Now().UnixNano()), 16)
		for len(s) < 2*n {
			s += s
		}
		return s[:2*n]
	}
	return hex.EncodeToString(b)
}
