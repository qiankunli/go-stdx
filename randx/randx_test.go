package randx

import (
	"regexp"
	"testing"
)

func TestHex(t *testing.T) {
	hexRe := regexp.MustCompile(`^[0-9a-f]+$`)
	for _, n := range []int{1, 6, 8, 16} {
		got := Hex(n)
		if len(got) != 2*n || !hexRe.MatchString(got) {
			t.Fatalf("Hex(%d) = %q", n, got)
		}
	}
	if got := Hex(0); got != "" {
		t.Fatalf("Hex(0) = %q", got)
	}
	if got := Hex(-1); got != "" {
		t.Fatalf("Hex(-1) = %q", got)
	}
	a, b := Hex(8), Hex(8)
	if a == b {
		t.Fatal("two draws collided — entropy source broken?")
	}
}
