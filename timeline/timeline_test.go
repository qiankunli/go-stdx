package timeline

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestTimelineStepAndSnapshot(t *testing.T) {
	tl := New("sandbox_start", Field{Key: "conversation_id", Value: "conv-1"})
	startedAt := time.Now().Add(-10 * time.Millisecond)
	tl.StepSince(startedAt, "prepare", Field{Key: "path", Value: "create"})

	snapshot := tl.Snapshot()
	if snapshot.Operation != "sandbox_start" || snapshot.Finished {
		t.Fatalf("snapshot = %+v", snapshot)
	}
	if snapshot.Duration() <= 0 {
		t.Fatalf("duration = %s, want positive", snapshot.Duration())
	}
	if len(snapshot.Steps) != 1 || snapshot.Steps[0].Message != "prepare" {
		t.Fatalf("steps = %+v", snapshot.Steps)
	}
	if snapshot.Steps[0].Duration < 10*time.Millisecond {
		t.Fatalf("step duration = %s, want at least 10ms", snapshot.Steps[0].Duration)
	}
}

func TestTimelineSnapshotIsIndependent(t *testing.T) {
	tl := New("sandbox_start", Field{Key: "result", Value: "running"})
	tl.Step("prepare", Field{Key: "outcome", Value: "created"})
	partial := tl.Snapshot()

	tl.AddField(Field{Key: "result", Value: "success"})
	tl.Step("ready")
	final := tl.Finish(Field{Key: "path", Value: "created"})

	if len(partial.Steps) != 1 || partial.Fields[0].Value != "running" {
		t.Fatalf("partial snapshot mutated: %+v", partial)
	}
	if !final.Finished || len(final.Steps) != 2 {
		t.Fatalf("final snapshot = %+v", final)
	}
	if got, ok := fieldValue(final.Fields, "result"); !ok || got != "success" {
		t.Fatalf("result field = %v, %v", got, ok)
	}
}

func TestTimelineFinishIsIdempotent(t *testing.T) {
	tl := New("sandbox_exec")
	first := tl.Finish(Field{Key: "result", Value: "success"})

	if _, ok := tl.Step("ignored"); ok {
		t.Fatal("Step after Finish returned ok")
	}
	tl.AddField(Field{Key: "result", Value: "failed"})
	second := tl.Finish(Field{Key: "result", Value: "failed"})
	third := tl.Snapshot()

	if !first.EndTime.Equal(second.EndTime) {
		t.Fatalf("finish time changed: %s != %s", first.EndTime, second.EndTime)
	}
	if !second.EndTime.Equal(third.EndTime) {
		t.Fatalf("snapshot time changed after finish: %s != %s", second.EndTime, third.EndTime)
	}
	if len(second.Steps) != 0 {
		t.Fatalf("steps after finish = %+v", second.Steps)
	}
	if got, _ := fieldValue(second.Fields, "result"); got != "success" {
		t.Fatalf("result field = %v, want success", got)
	}
}

func TestTimelineConcurrentUse(t *testing.T) {
	tl := New("sandbox_start")
	var wg sync.WaitGroup
	for i := range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tl.Step(fmt.Sprintf("step-%d", i))
			_ = tl.Snapshot()
		}()
	}
	wg.Wait()

	final := tl.Finish()
	if len(final.Steps) != 50 {
		t.Fatalf("step count = %d, want 50", len(final.Steps))
	}
}

func TestContext(t *testing.T) {
	tl := New("sandbox_start")
	ctx := NewContext(context.Background(), tl)

	got, ok := FromContext(ctx)
	if !ok || got != tl {
		t.Fatalf("FromContext() = %p, %v; want %p, true", got, ok, tl)
	}
	if got, ok := FromContext(context.Background()); got != nil || ok {
		t.Fatalf("empty FromContext() = %p, %v", got, ok)
	}
	if got, ok := FromContext(nil); got != nil || ok {
		t.Fatalf("nil FromContext() = %p, %v", got, ok)
	}
}

func fieldValue(fields []Field, key string) (any, bool) {
	for _, field := range fields {
		if field.Key == key {
			return field.Value, true
		}
	}
	return nil, false
}
