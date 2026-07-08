package stringsx

import (
	"reflect"
	"testing"
)

func TestTruncate(t *testing.T) {
	if got := Truncate("hello", 3); got != "hel" {
		t.Fatalf("cut: %q", got)
	}
	if got := Truncate("你好abc", 3); got != "你好a" {
		t.Fatalf("rune cut: %q", got)
	}
	if got := Truncate("hi", 3); got != "hi" {
		t.Fatalf("fits: %q", got)
	}
	if got := Truncate("", 3); got != "" {
		t.Fatalf("empty: %q", got)
	}
	if got := Truncate("hello", -1); got != "" {
		t.Fatalf("negative max: %q", got)
	}
}

func TestTruncateEllipsis(t *testing.T) {
	if got := TruncateEllipsis("hello world", 5); got != "hello..." {
		t.Fatalf("cut: %q", got)
	}
	if got := TruncateEllipsis("你好abc", 3); got != "你好a..." {
		t.Fatalf("rune cut: %q", got)
	}
	if got := TruncateEllipsis("hi", 5); got != "hi" {
		t.Fatalf("fits: %q", got)
	}
	if got := TruncateEllipsis("hello", -1); got != "..." {
		t.Fatalf("negative max: %q", got)
	}
	if got := TruncateEllipsis("", -1); got != "" {
		t.Fatalf("negative max empty: %q", got)
	}
}

func TestFirstNonBlank(t *testing.T) {
	if got := FirstNonBlank("", "  ", "\t", "x", "y"); got != "x" {
		t.Fatalf("got %q", got)
	}
	if got := FirstNonBlank("", "   "); got != "" {
		t.Fatalf("all blank: %q", got)
	}
	if got := FirstNonBlank(); got != "" {
		t.Fatalf("no args: %q", got)
	}
}

func TestSplitAndTrim(t *testing.T) {
	if got := SplitAndTrim("a, b ,,c", ","); !reflect.DeepEqual(got, []string{"a", "b", "c"}) {
		t.Fatalf("got %v", got)
	}
	if got := SplitAndTrim("  ", ","); got != nil {
		t.Fatalf("blank: %v", got)
	}
	if got := SplitAndTrim("", ","); got != nil {
		t.Fatalf("empty: %v", got)
	}
}
