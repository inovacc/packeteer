package safety

import (
	"context"
	"fmt"
	"sync/atomic"
)

// CaptureLimiter limits the number of concurrent live captures.
type CaptureLimiter struct {
	max     int32
	active  atomic.Int32
	sem     chan struct{}
}

// NewCaptureLimiter creates a limiter with the given max concurrent captures.
func NewCaptureLimiter(max int) *CaptureLimiter {
	if max <= 0 {
		max = 3
	}
	return &CaptureLimiter{
		max: int32(max),
		sem: make(chan struct{}, max),
	}
}

// Acquire attempts to acquire a capture slot. Returns a release function.
// Returns an error if the context is cancelled or all slots are occupied.
func (l *CaptureLimiter) Acquire(ctx context.Context) (release func(), err error) {
	select {
	case l.sem <- struct{}{}:
		l.active.Add(1)
		released := atomic.Bool{}
		return func() {
			if released.CompareAndSwap(false, true) {
				l.active.Add(-1)
				<-l.sem
			}
		}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return nil, fmt.Errorf("max concurrent captures reached (%d/%d active)", l.active.Load(), l.max)
	}
}

// Active returns the number of currently active captures.
func (l *CaptureLimiter) Active() int {
	return int(l.active.Load())
}

// Max returns the maximum allowed concurrent captures.
func (l *CaptureLimiter) Max() int {
	return int(l.max)
}
