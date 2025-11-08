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

package strategy

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/stoploss"
)

var (
	errStopLossRateInvalid   = errors.New("stop loss rate must be between 0 and 1")
	errTakeProfitRateInvalid = errors.New("take profit rate must be between 0 and 1")
)

// FixedPercentStop represents a stop loss strategy based on fixed percentage
type FixedPercentStop struct {
	stoploss.BaseResolver
	threshold    decimal.Decimal
	tolerancePct decimal.Decimal
	LastPrice    decimal.Decimal
}

// FixedPercentProfit represents a take profit strategy based on fixed percentage
type FixedPercentProfit struct {
	stoploss.BaseResolver
	threshold    decimal.Decimal
	tolerancePct decimal.Decimal
	LastPrice    decimal.Decimal
}

// NewFixedPercentStop creates a FixedPercentStopLoss base on fixed percentage
func NewFixedPercentStop(entryPrice, stopLossPct decimal.Decimal, callback stoploss.DefaultCallback) (stoploss.FixedStopLoss, error) {
	if stopLossPct.IsNegative() || stopLossPct.GreaterThan(decimal.NewFromInt(1)) {
		return nil, errStopLossRateInvalid
	}
	s := &FixedPercentStop{
		LastPrice:    entryPrice,
		threshold:    entryPrice.Mul(decimal.NewFromInt(1).Sub(stopLossPct)),
		tolerancePct: stopLossPct,
		BaseResolver: stoploss.BaseResolver{
			Active:   true,
			Callback: callback,
		},
	}
	return s, nil
}

// NewFixedPercentProfit creates a FixedPercentTakeProfit based on fixed percentage
func NewFixedPercentProfit(entryPrice, takeProfitPct decimal.Decimal, callback stoploss.DefaultCallback) (stoploss.FixedTakeProfit, error) {
	if takeProfitPct.IsNegative() || takeProfitPct.GreaterThan(decimal.NewFromInt(1)) {
		return nil, errTakeProfitRateInvalid
	}
	s := &FixedPercentProfit{
		LastPrice:    entryPrice,
		threshold:    entryPrice.Mul(decimal.NewFromInt(1).Add(takeProfitPct)),
		tolerancePct: takeProfitPct,
		BaseResolver: stoploss.BaseResolver{
			Active:   true,
			Callback: callback,
		},
	}
	return s, nil
}

// CalculateStopLoss represents first update last price and calculate stop loss then update threshold
func (f *FixedPercentStop) CalculateStopLoss(currentPrice decimal.Decimal) (decimal.Decimal, error) {
	if !f.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	f.LastPrice = currentPrice
	f.threshold = currentPrice.Mul(decimal.NewFromInt(1).Sub(f.tolerancePct))
	return f.threshold, nil
}

// ShouldTriggerStopLoss checks if the stop loss should be triggered
func (f *FixedPercentStop) ShouldTriggerStopLoss(currentPrice decimal.Decimal) (bool, error) {
	if !f.Active {
		return false, stoploss.ErrStatusInvalid
	}
	if currentPrice.LessThanOrEqual(f.threshold) {
		err := f.Trigger(stoploss.TRIGGERED_REASON_FIXED_PERCENTCILE_STOPLOSS)
		if err != nil {
			return true, stoploss.ErrCallBackFail
		}
		return true, nil
	}
	return false, nil
}

// GetStopLoss returns the current stop loss threshold
func (f *FixedPercentStop) GetStopLoss() (decimal.Decimal, error) {
	if !f.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return f.threshold, nil
}

// ReSetStopLosser resets the stop loss based on the current price
func (f *FixedPercentStop) ReSetStopLosser(currentPrice decimal.Decimal) error {
	if !f.Active {
		return stoploss.ErrStatusInvalid
	}
	f.LastPrice = currentPrice
	f.threshold = currentPrice.Mul(decimal.NewFromInt(1).Sub(f.tolerancePct))
	f.Active = true
	return nil
}

// CalculateTakeProfit represents first update last price and calculate take profit then update threshold
func (f *FixedPercentProfit) CalculateTakeProfit(currentPrice decimal.Decimal) (decimal.Decimal, error) {
	if !f.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	f.LastPrice = currentPrice
	f.threshold = currentPrice.Mul(decimal.NewFromInt(1).Add(f.tolerancePct))
	return f.threshold, nil
}

// ShouldTriggerTakeProfit checks if the take profit should be triggered
func (f *FixedPercentProfit) ShouldTriggerTakeProfit(currentPrice decimal.Decimal) (bool, error) {
	if !f.Active {
		return false, stoploss.ErrStatusInvalid
	}
	if currentPrice.GreaterThanOrEqual(f.threshold) {
		err := f.Trigger(stoploss.TRIGGERED_REASON_FIXED_PERCENTCILE_TAKEPROFIT)
		if err != nil {
			return true, stoploss.ErrCallBackFail
		}
		return true, nil
	}
	return false, nil
}

// GetTakeProfit returns the current take profit threshold
func (f *FixedPercentProfit) GetTakeProfit() (decimal.Decimal, error) {
	if !f.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return f.threshold, nil
}

// ReSetTakeProfiter resets the take profit based on the current price
func (f *FixedPercentProfit) ReSetTakeProfiter(currentPrice decimal.Decimal) error {
	if !f.Active {
		return stoploss.ErrStatusInvalid
	}
	f.LastPrice = currentPrice
	f.threshold = currentPrice.Mul(decimal.NewFromInt(1).Add(f.tolerancePct))
	f.Active = true
	return nil
}
