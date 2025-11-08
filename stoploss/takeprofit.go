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

package stoploss

import "github.com/shopspring/decimal"

type HybridWithoutTime interface {
	Hybrid
	TakeProfitCond
	StopLossCond
}

type HybridWithTime interface {
	Hybrid
	TakeProfitCondT
	StopLossCondT
}

type Hybrid interface {
	Calculate(currentPrice decimal.Decimal) (decimal.Decimal, decimal.Decimal, error)
	Trigger(reason string) error
	ReSet(currentPrice decimal.Decimal) error
	Deactivate() error
}

// Time-based
type TimeBasedTakeProfit interface {
	TakeProfit
	TakeProfitCondT
}

// Fixed strategy
type FixedTakeProfit interface {
	TakeProfit
	TakeProfitCond
}

// Fixed-ATR
type FixedVolatilityTakeProfit interface {
	FixedTakeProfit
	UpdateATR(currentATR decimal.Decimal) error
}

// Fixed-Moving Average
type FixedMATakeProfit interface {
	FixedTakeProfit
	SetMA(value decimal.Decimal)
}

// general TakeProfit interface
type TakeProfit interface {
	CalculateTakeProfit(currentPrice decimal.Decimal) (decimal.Decimal, error)
	Trigger(reason string) error
	ReSetTakeProfiter(currentPrice decimal.Decimal) error
	Deactivate() error
}

// TakeProfit Condition
type TakeProfitCond interface {
	ShouldTriggerTakeProfit(currentPrice decimal.Decimal) (bool, error)
	GetTakeProfit() (decimal.Decimal, error)
}

// TakeProfit Condition with timestamp
type TakeProfitCondT interface {
	ShouldTriggerTakeProfit(currentPrice decimal.Decimal, timestamp int64) (bool, error)
	GetTakeProfit() (decimal.Decimal, error)
	GetTimeThreshold() (int64, error)
}
