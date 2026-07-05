package ptrx

import "testing"

func TestTo(t *testing.T) {
	p := To(42)
	if p == nil || *p != 42 {
		t.Fatalf("To(42) = %v", p)
	}
	b := To(true)
	if b == nil || !*b {
		t.Fatalf("To(true) = %v", b)
	}
}

func TestValue(t *testing.T) {
	if got := Value(To("x")); got != "x" {
		t.Fatalf("Value = %q", got)
	}
	if got := Value[string](nil); got != "" {
		t.Fatalf("nil string: %q", got)
	}
	if got := Value[int](nil); got != 0 {
		t.Fatalf("nil int: %d", got)
	}
}

func TestFormatOr(t *testing.T) {
	if got := FormatOr[bool](nil, "default"); got != "default" {
		t.Fatalf("nil: %q", got)
	}
	if got := FormatOr(To(false), "default"); got != "false" {
		t.Fatalf("false: %q", got)
	}
	if got := FormatOr(To(int64(30)), "default"); got != "30" {
		t.Fatalf("int64: %q", got)
	}
}
