package shellx

import (
	"os/exec"
	"runtime"
	"testing"
)

func TestQuote(t *testing.T) {
	cases := map[string]string{
		"plain":        "'plain'",
		"has space":    "'has space'",
		"don't":        `'don'\''t'`,
		"$HOME `x` !h": "'$HOME `x` !h'",
		"":             "''",
	}
	for in, want := range cases {
		if got := Quote(in); got != want {
			t.Fatalf("Quote(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestQuoteRoundTrip pushes hostile strings through a real shell: what comes
// out of echo must be byte-identical to what went in.
func TestQuoteRoundTrip(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX shells only")
	}
	for _, s := range []string{"don't", "$HOME", "a;rm -rf /", "`date`", `back\slash`, "两个 words"} {
		out, err := exec.Command("sh", "-c", "printf %s "+Quote(s)).Output()
		if err != nil {
			t.Fatalf("sh: %v", err)
		}
		if string(out) != s {
			t.Fatalf("round trip: %q -> %q", s, out)
		}
	}
}
