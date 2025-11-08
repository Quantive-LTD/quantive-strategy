// Copyright 2025 Quantive. All rights reserved.

// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// https://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metric

import (
	"sync/atomic"

	"github.com/shopspring/decimal"
)

const (
	// DefaultCounterPrecision defines the default precision for decimal counters.
	DefaultCounterPrecision int32 = 8
)

// CounterPrecisionSnapshot represents a static snapshot of a counter's value with decimal precision.
type counterPrecisionSnapshot decimal.Decimal

// Count returns the value of the counter's static snapshot with decimal precision.
func (c counterPrecisionSnapshot) Count() decimal.Decimal {
	return decimal.Decimal(c)
}

// counterPrecision is a thread-safe counter for decimal.Decimal values with specified precision.
type counterPrecision struct {
	value atomic.Value // holds decimal.Decimal
	prec  int32        // precision for decimal operations
}

// NewCounterPrecision creates and returns a new instance of counterPrecision with the specified precision.
// If precision is 0, DefaultCounterPrecision is used.
func NewCounterPrecision(precision int32) *counterPrecision {
	if precision == 0 {
		precision = DefaultCounterPrecision
	}
	c := &counterPrecision{prec: precision}
	c.value.Store(decimal.NewFromInt(0).Round(precision))
	return c
}

// Inc increments the counter by the specified delta with decimal precision.
func (c *counterPrecision) Inc(delta decimal.Decimal) {
	for {
		current := c.value.Load().(decimal.Decimal)
		newValue := current.Add(delta).Round(c.prec)
		if c.value.CompareAndSwap(current, newValue) {
			return
		}
	}
}

// Dec decrements the counter by the specified delta with decimal precision.
func (c *counterPrecision) Dec(delta decimal.Decimal) {
	for {
		current := c.value.Load().(decimal.Decimal)
		newValue := current.Sub(delta).Round(c.prec)
		if c.value.CompareAndSwap(current, newValue) {
			return
		}
	}
}

// Clear resets the counter to zero with decimal precision.
func (c *counterPrecision) Clear() {
	c.value.Store(decimal.NewFromInt(0).Round(c.prec))
}

// Fork creates a new counterPrecision instance but shares the same value.
func (c *counterPrecision) Fork() *counterPrecision {
	forked := NewCounterPrecision(c.prec)
	forked.Inc(c.value.Load().(decimal.Decimal))
	return forked
}

// Snapshot returns a static snapshot of the current counter value with decimal precision.
func (c *counterPrecision) Snapshot() counterPrecisionSnapshot {
	return counterPrecisionSnapshot(c.value.Load().(decimal.Decimal))
}
