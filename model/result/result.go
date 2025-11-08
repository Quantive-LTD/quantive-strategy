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

package result

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type StrategyStat struct {
	PriceThreshold decimal.Decimal
}

type StrategyGeneralResult struct {
	StrategyName  string
	StrategyType  string
	Triggered     bool
	TriggerType   string
	LastPrice     decimal.Decimal
	Stat          StrategyStat
	LastTime      time.Time
	TimeThreshold time.Duration
	Error         error
}

type StrategyHybridResult struct {
	StrategyName  string
	StrategyType  string
	Triggered     bool
	TriggerType   string
	LastTime      time.Time
	LastPrice     decimal.Decimal
	TimeThreshold time.Duration
	stopStat      StrategyStat
	profitStat    StrategyStat
	Error         error
}

func (sr *StrategyGeneralResult) Marshall() map[string]interface{} {
	return map[string]interface{}{
		"StrategyName":  sr.StrategyName,
		"StrategyType":  sr.StrategyType,
		"Triggered":     sr.Triggered,
		"TriggerType":   sr.TriggerType,
		"Stat":          sr.Stat,
		"LastTime":      sr.LastTime,
		"TimeThreshold": sr.TimeThreshold,
		"Error":         sr.Error,
	}
}

func (sr *StrategyHybridResult) Marshall() map[string]interface{} {
	return map[string]interface{}{
		"StrategyName":  sr.StrategyName,
		"StrategyType":  sr.StrategyType,
		"Triggered":     sr.Triggered,
		"TriggerType":   sr.TriggerType,
		"LastTime":      sr.LastTime,
		"TimeThreshold": sr.TimeThreshold,
		"StopStat":      sr.stopStat,
		"ProfitStat":    sr.profitStat,
		"Error":         sr.Error,
	}
}

func NewGeneral(strategyName, strategyType, triggerType string, lastPrice, priceThreshold decimal.Decimal, lastTime time.Time, timeThreshold time.Duration) *StrategyGeneralResult {
	return &StrategyGeneralResult{
		StrategyName:  strategyName,
		StrategyType:  strategyType,
		TriggerType:   triggerType,
		LastTime:      lastTime,
		TimeThreshold: timeThreshold,
		LastPrice:     lastPrice,
		Stat: StrategyStat{
			PriceThreshold: priceThreshold,
		},
	}
}

func NewHybrid(strategyName, strategyType, triggerType string, LastPrice, stopPriceThreshold, profitPriceThreshold decimal.Decimal, lastTime time.Time, timeThreshold time.Duration) *StrategyHybridResult {
	return &StrategyHybridResult{
		StrategyName:  strategyName,
		StrategyType:  strategyType,
		TriggerType:   triggerType,
		LastTime:      lastTime,
		TimeThreshold: timeThreshold,
		LastPrice:     LastPrice,
		stopStat: StrategyStat{
			PriceThreshold: stopPriceThreshold,
		},
		profitStat: StrategyStat{
			PriceThreshold: profitPriceThreshold,
		},
	}
}

func (sr *StrategyGeneralResult) SetTriggered(triggered bool) {
	sr.Triggered = triggered
}

func (sr *StrategyGeneralResult) SetError(err error) {
	sr.Error = err
}

func (sr *StrategyHybridResult) SetTriggered(triggered bool) {
	sr.Triggered = triggered
}

func (sr *StrategyHybridResult) SetError(err error) {
	sr.Error = err
}

func (sr *StrategyGeneralResult) String() string {
	errorStr := "<nil>"
	if sr.Error != nil {
		errorStr = sr.Error.Error()
	}
	return fmt.Sprintf("StrategyGeneralResult{StrategyName: %s, StrategyType: %s, Triggered: %t, TriggerType: %s, LastPrice: %s, PriceThreshold: %s, LastTime: %s, TimeThreshold: %s, Error: %v}",
		sr.StrategyName,
		sr.StrategyType,
		sr.Triggered,
		sr.TriggerType,
		sr.LastPrice.String(),
		sr.Stat.PriceThreshold.String(),
		sr.LastTime.String(),
		sr.TimeThreshold.String(),
		errorStr,
	)
}

func (sr *StrategyHybridResult) String() string {
	errorStr := "<nil>"
	if sr.Error != nil {
		errorStr = sr.Error.Error()
	}
	return fmt.Sprintf("StrategyHybridResult{StrategyName: %s, StrategyType: %s, Triggered: %t, TriggerType: %s, LastPrice: %s, StopPriceThreshold: %s, ProfitPriceThreshold: %s, LastTime: %s, TimeThreshold: %s, Error: %s}",
		sr.StrategyName,
		sr.StrategyType,
		sr.Triggered,
		sr.TriggerType,
		sr.LastPrice.String(),
		sr.stopStat.PriceThreshold.String(),
		sr.profitStat.PriceThreshold.String(),
		sr.LastTime.String(),
		sr.TimeThreshold.String(),
		errorStr,
	)
}

func (sr *StrategyGeneralResult) JSONMarshall() string {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(sr.Marshall()); err != nil {
		return ""
	}
	return buf.String()
}

func (sr *StrategyHybridResult) JSONMarshall() string {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(sr.Marshall()); err != nil {
		return ""
	}
	return buf.String()
}
