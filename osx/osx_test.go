package osx

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestEnvStr(t *testing.T) {
	t.Setenv("STDX_T", "v")
	if got := EnvStr("STDX_T", "d"); got != "v" {
		t.Fatalf("set: %q", got)
	}
	if got := EnvStr("STDX_UNSET", "d"); got != "d" {
		t.Fatalf("unset: %q", got)
	}
	t.Setenv("STDX_T", "")
	if got := EnvStr("STDX_T", "d"); got != "d" {
		t.Fatalf("empty: %q", got)
	}
}

func TestEnvBool(t *testing.T) {
	t.Setenv("STDX_T", "true")
	if !EnvBool("STDX_T", false) {
		t.Fatal("true not parsed")
	}
	t.Setenv("STDX_T", "0")
	if EnvBool("STDX_T", true) {
		t.Fatal("0 not parsed as false")
	}
	t.Setenv("STDX_T", "nope")
	if !EnvBool("STDX_T", true) {
		t.Fatal("unparsable should fall back to default")
	}
	if EnvBool("STDX_UNSET", false) {
		t.Fatal("unset should fall back to default")
	}
}

func TestEnvInt(t *testing.T) {
	t.Setenv("STDX_T", "42")
	if got := EnvInt("STDX_T", 7); got != 42 {
		t.Fatalf("set: %d", got)
	}
	t.Setenv("STDX_T", "nope")
	if got := EnvInt("STDX_T", 7); got != 7 {
		t.Fatalf("unparsable: %d", got)
	}
	if got := EnvInt("STDX_UNSET", 7); got != 7 {
		t.Fatalf("unset: %d", got)
	}
}

func TestEnvInt64(t *testing.T) {
	t.Setenv("STDX_T", "9000000000")
	if got := EnvInt64("STDX_T", 1); got != 9000000000 {
		t.Fatalf("set: %d", got)
	}
	if got := EnvInt64("STDX_UNSET", 5); got != 5 {
		t.Fatalf("unset: %d", got)
	}
}

func TestEnvDuration(t *testing.T) {
	t.Setenv("STDX_T", "90s")
	if got := EnvDuration("STDX_T", time.Minute); got != 90*time.Second {
		t.Fatalf("unit: %v", got)
	}
	t.Setenv("STDX_T", "30") // bare integer reads as seconds
	if got := EnvDuration("STDX_T", time.Minute); got != 30*time.Second {
		t.Fatalf("bare: %v", got)
	}
	t.Setenv("STDX_T", "nope")
	if got := EnvDuration("STDX_T", time.Minute); got != time.Minute {
		t.Fatalf("unparsable: %v", got)
	}
}

func TestWriteFileAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "f.json")

	if err := WriteFileAtomic(path, []byte("one"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := WriteFileAtomic(path, []byte("two"), 0o644); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "two" {
		t.Fatalf("content: %q err=%v", data, err)
	}
	fi, err := os.Stat(path)
	if err != nil || fi.Mode().Perm() != 0o644 {
		t.Fatalf("mode: %v err=%v", fi.Mode(), err)
	}
	// No temp litter left behind.
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 1 {
		t.Fatalf("litter: %v err=%v", entries, err)
	}
}
