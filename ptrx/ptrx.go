// Package ptrx holds the pointer helpers Go's optional-field style forces on
// every service: taking the address of a literal, and reading through a
// possibly-nil pointer. The stdlib has no answer (k8s.io/utils/ptr exists
// precisely because of that gap).
package ptrx

import "fmt"

// To returns a pointer to v — the "address of a literal" helper that makes
// optional struct fields writable in one expression: Enabled: ptrx.To(true).
func To[T any](v T) *T { return &v }

// Value returns *p, or T's zero value when p is nil.
func Value[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// FormatOr renders *p with fmt.Sprint, or fallback when p is nil — for logs
// where a nil optional should read as something ("default", "unset") instead
// of a pointer literal or a misleading zero.
func FormatOr[T any](p *T, fallback string) string {
	if p == nil {
		return fallback
	}
	return fmt.Sprint(*p)
}
