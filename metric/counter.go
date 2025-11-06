package metric

import "sync/atomic"

// counterInt64Snapshot represents a static snapshot of a counter's value.
type counterInt64Snapshot int64

// Count returns the value of the counter's static snapshot for int64.
func (c counterInt64Snapshot) Count() int64 {
	return int64(c)
}

// counterInt64 is a thread-safe counter for int64 values.
type counterInt64 atomic.Int64

// NewCounterInt64 creates and returns a new instance of counterInt64.
func NewCounterInt64() *counterInt64 {
	return new(counterInt64)
}

// Inc increments the counter by the specified delta.
func (c *counterInt64) Inc(delta int64) {
	(*atomic.Int64)(c).Add(delta)
}

// Dec decrements the counter by the specified delta.
func (c *counterInt64) Dec(delta int64) {
	(*atomic.Int64)(c).Add(-delta)
}

// Clear resets the counter to zero.
func (c *counterInt64) Clear() {
	(*atomic.Int64)(c).Store(0)
}

// Fork creates a new counterInt64 instance but shares the same value.
func (c *counterInt64) Fork() *counterInt64 {
	var forked atomic.Int64
	forked.Store((*atomic.Int64)(c).Load())
	return (*counterInt64)(&forked)
}

// Snapshot returns a static snapshot of the current counter value
func (c *counterInt64) Snapshot() counterInt64Snapshot {
	return counterInt64Snapshot((*atomic.Int64)(c).Load())
}
