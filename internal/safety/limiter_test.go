package safety

import (
	"context"
	"sync"
	"testing"
)

func TestCaptureLimiter_Basic(t *testing.T) {
	lim := NewCaptureLimiter(2)

	if lim.Active() != 0 {
		t.Fatalf("expected 0 active, got %d", lim.Active())
	}
	if lim.Max() != 2 {
		t.Fatalf("expected max 2, got %d", lim.Max())
	}

	r1, err := lim.Acquire(context.Background())
	if err != nil {
		t.Fatalf("acquire 1 failed: %v", err)
	}
	if lim.Active() != 1 {
		t.Fatalf("expected 1 active, got %d", lim.Active())
	}

	r2, err := lim.Acquire(context.Background())
	if err != nil {
		t.Fatalf("acquire 2 failed: %v", err)
	}
	if lim.Active() != 2 {
		t.Fatalf("expected 2 active, got %d", lim.Active())
	}

	// Third should fail.
	_, err = lim.Acquire(context.Background())
	if err == nil {
		t.Fatal("expected error on third acquire")
	}

	// Release one, then acquire should work.
	r1()
	if lim.Active() != 1 {
		t.Fatalf("expected 1 active after release, got %d", lim.Active())
	}

	r3, err := lim.Acquire(context.Background())
	if err != nil {
		t.Fatalf("acquire 3 failed: %v", err)
	}

	r2()
	r3()
	if lim.Active() != 0 {
		t.Fatalf("expected 0 active after all releases, got %d", lim.Active())
	}
}

func TestCaptureLimiter_DoubleRelease(t *testing.T) {
	lim := NewCaptureLimiter(1)
	r, _ := lim.Acquire(context.Background())
	r()
	r() // Should not panic or double-decrement.
	if lim.Active() != 0 {
		t.Fatalf("expected 0 active, got %d", lim.Active())
	}
}

func TestCaptureLimiter_Concurrent(t *testing.T) {
	lim := NewCaptureLimiter(5)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			release, err := lim.Acquire(context.Background())
			if err != nil {
				return
			}
			// Simulate work.
			release()
		}()
	}

	wg.Wait()
	if lim.Active() != 0 {
		t.Fatalf("expected 0 active after all goroutines, got %d", lim.Active())
	}
}

func TestCaptureLimiter_CancelledContext(t *testing.T) {
	lim := NewCaptureLimiter(1)
	r, _ := lim.Acquire(context.Background())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := lim.Acquire(ctx)
	if err == nil {
		t.Fatal("expected error with cancelled context")
	}

	r()
}

func TestCaptureLimiter_DefaultMax(t *testing.T) {
	lim := NewCaptureLimiter(0)
	if lim.Max() != 3 {
		t.Fatalf("expected default max 3, got %d", lim.Max())
	}
}
