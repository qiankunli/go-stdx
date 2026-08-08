// Package timeline records the steps of an in-process operation.
//
// A Timeline only collects data. Callers decide whether a Snapshot is exported
// to logs, OpenTelemetry, a database, or somewhere else.
package timeline

import (
	"context"
	"sync"
	"time"
)

// Field is a key-value pair attached to a timeline or one of its steps.
type Field struct {
	Key   string
	Value any
}

// Step is one recorded point in an operation.
// Duration is measured from the previous step by Step, or from the caller's
// supplied start time by StepSince.
type Step struct {
	Time     time.Time
	Message  string
	Fields   []Field
	Duration time.Duration
}

// Snapshot is a point-in-time copy of a Timeline. EndTime is the snapshot time
// while the operation is running and the stable finish time after Finish. The
// Field slices are copied, while values held by Field.Value are copied as-is.
type Snapshot struct {
	Operation string
	Fields    []Field
	StartTime time.Time
	EndTime   time.Time
	Finished  bool
	Steps     []Step
}

// Duration returns the operation duration represented by the snapshot.
func (s Snapshot) Duration() time.Duration {
	return s.EndTime.Sub(s.StartTime)
}

// Timeline records one operation and its steps. Its methods are safe for
// concurrent use. Finish freezes the timeline; later mutations are ignored.
type Timeline struct {
	mu sync.RWMutex

	operation string
	fields    []Field
	startTime time.Time
	endTime   time.Time
	steps     []Step
}

// New starts a timeline for operation.
func New(operation string, fields ...Field) *Timeline {
	return &Timeline{
		operation: operation,
		fields:    cloneFields(fields),
		startTime: time.Now(),
	}
}

// AddField adds fields to the operation. A field with an existing key replaces
// the previous value.
func (t *Timeline) AddField(fields ...Field) {
	if t == nil || len(fields) == 0 {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.endTime.IsZero() {
		return
	}
	t.fields = mergeFields(t.fields, fields)
}

// Step records a step and measures its duration from the previous step, or
// from the start of the operation when it is the first step. It returns false
// when the timeline has already finished.
func (t *Timeline) Step(message string, fields ...Field) (Step, bool) {
	if t == nil {
		return Step{}, false
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.endTime.IsZero() {
		return Step{}, false
	}

	now := time.Now()
	startedAt := t.startTime
	if len(t.steps) > 0 {
		startedAt = t.steps[len(t.steps)-1].Time
	}
	return t.appendStep(now, now.Sub(startedAt), message, fields), true
}

// StepSince records a step and measures its duration from startedAt. It is
// useful for nested or concurrent work where the previous step is not the
// beginning of this step. It returns false when the timeline has already
// finished.
func (t *Timeline) StepSince(startedAt time.Time, message string, fields ...Field) (Step, bool) {
	if t == nil {
		return Step{}, false
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.endTime.IsZero() {
		return Step{}, false
	}
	now := time.Now()
	return t.appendStep(now, now.Sub(startedAt), message, fields), true
}

func (t *Timeline) appendStep(at time.Time, duration time.Duration, message string, fields []Field) Step {
	step := Step{
		Time:     at,
		Message:  message,
		Fields:   cloneFields(fields),
		Duration: duration,
	}
	t.steps = append(t.steps, step)
	step.Fields = cloneFields(step.Fields)
	return step
}

// Snapshot returns a point-in-time copy without finishing the operation.
func (t *Timeline) Snapshot() Snapshot {
	if t == nil {
		return Snapshot{}
	}

	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.snapshotLocked(time.Now())
}

// Finish freezes the timeline and returns its final snapshot. Fields are
// merged atomically only by the first Finish call; subsequent calls return the
// same final snapshot.
func (t *Timeline) Finish(fields ...Field) Snapshot {
	if t == nil {
		return Snapshot{}
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	if t.endTime.IsZero() {
		t.fields = mergeFields(t.fields, fields)
		t.endTime = time.Now()
	}
	return t.snapshotLocked(t.endTime)
}

func (t *Timeline) snapshotLocked(snapshotTime time.Time) Snapshot {
	endTime := t.endTime
	if endTime.IsZero() {
		endTime = snapshotTime
	}
	steps := make([]Step, len(t.steps))
	for i, step := range t.steps {
		steps[i] = step
		steps[i].Fields = cloneFields(step.Fields)
	}
	return Snapshot{
		Operation: t.operation,
		Fields:    cloneFields(t.fields),
		StartTime: t.startTime,
		EndTime:   endTime,
		Finished:  !t.endTime.IsZero(),
		Steps:     steps,
	}
}

func cloneFields(fields []Field) []Field {
	if fields == nil {
		return nil
	}
	cloned := make([]Field, len(fields))
	copy(cloned, fields)
	return cloned
}

func mergeFields(existing, added []Field) []Field {
	for _, field := range added {
		updated := false
		for i := range existing {
			if existing[i].Key == field.Key {
				existing[i] = field
				updated = true
				break
			}
		}
		if !updated {
			existing = append(existing, field)
		}
	}
	return existing
}

type contextKey struct{}

// NewContext returns a child context carrying t.
func NewContext(ctx context.Context, t *Timeline) context.Context {
	return context.WithValue(ctx, contextKey{}, t)
}

// FromContext returns the Timeline carried by ctx, if any.
func FromContext(ctx context.Context) (*Timeline, bool) {
	if ctx == nil {
		return nil, false
	}
	t, ok := ctx.Value(contextKey{}).(*Timeline)
	return t, ok && t != nil
}
