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
	"encoding/json"
	"sync"
)

// GaugeInfoSnapshot represents a snapshot of the GaugeInfo at a point in time
type GaugeInfoSnapshot GaugeInfoValue

// Value returns the GaugeInfoValue representation of the snapshot
func (g GaugeInfoSnapshot) Value() GaugeInfoValue { return GaugeInfoValue(g) }

// GaugeInfoValue represents key-value pairs stored in GaugeInfo
type GaugeInfoValue map[string]string

// String returns the JSON representation of the GaugeInfoValue
func (g GaugeInfoValue) String() string { data, _ := json.Marshal(g); return string(data) }

// GaugeInfo is a thread-safe structure for storing key-value pairs
type GaugeInfo struct {
	mutex sync.Mutex
	value GaugeInfoValue
}

// NewGaugeInfo creates and returns a new GaugeInfo instance
func NewGaugeInfo() *GaugeInfo {
	return &GaugeInfo{
		value: make(GaugeInfoValue),
	}
}

// Snapshot creates a snapshot of the current GaugeInfo state
func (g *GaugeInfo) Snapshot() GaugeInfoSnapshot {
	return GaugeInfoSnapshot(g.value)
}

// Set sets a key-value pair in the GaugeInfo
func (g *GaugeInfo) Set(key, val string) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.value[key] = val
}
