// Package osx extends the standard os package: typed environment lookups
// with defaults (the config-from-env boilerplate every service re-writes)
// and an atomic WriteFile.
package osx

import (
	"os"
	"strconv"
	"time"
)

// EnvStr returns the value of the environment variable key, or def when the
// variable is unset or empty.
func EnvStr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// EnvBool returns the environment variable key parsed with
// strconv.ParseBool ("1", "t", "true", ... case-insensitively), or def when
// the variable is unset, empty, or unparsable.
func EnvBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}

// EnvInt returns the environment variable key parsed as an int, or def when
// the variable is unset, empty, or unparsable.
func EnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

// EnvInt64 returns the environment variable key parsed as an int64, or def
// when the variable is unset, empty, or unparsable.
func EnvInt64(key string, def int64) int64 {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return n
		}
	}
	return def
}

// EnvDuration returns the environment variable key parsed as a
// time.Duration, or def when the variable is unset, empty, or unparsable.
// A bare integer (no unit) is read as SECONDS — env files and Helm charts
// routinely write "30" where Go wants "30s".
func EnvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
		if n, err := strconv.Atoi(v); err == nil {
			return time.Duration(n) * time.Second
		}
	}
	return def
}
