package uuid

import (
	"regexp"
	"testing"
	"time"
)

var v7Re = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func TestV7Format(t *testing.T) {
	for i := 0; i < 10; i++ {
		got := V7()
		if !v7Re.MatchString(got) {
			t.Fatalf("V7() = %q, not a version-7 UUID", got)
		}
	}
}

func TestV7TimeOrdered(t *testing.T) {
	a := V7()
	time.Sleep(2 * time.Millisecond) // cross a millisecond boundary
	b := V7()
	if !(a < b) {
		t.Fatalf("V7 not time-ordered: %q then %q", a, b)
	}
}

func TestV7Unique(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		id := V7()
		if seen[id] {
			t.Fatalf("collision after %d draws: %s", i, id)
		}
		seen[id] = true
	}
}
