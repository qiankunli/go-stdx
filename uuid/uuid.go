// Package uuid wraps github.com/google/uuid with the string / hex id shapes
// services actually pass around, so callers don't re-do .String() or the
// hex-encode dance at every call site. Generating the UUID itself is left to
// google/uuid — a mature, correct implementation is not worth reinventing.
package uuid

import guuid "github.com/google/uuid"

// V4 returns a random (version-4) UUID in canonical dashed form. On the
// practically impossible entropy failure it returns the zero UUID rather than
// panicking — these are correlation ids, where a crash is the worse trade.
func V4() string {
	u, _ := guuid.NewRandom()
	return u.String()
}
