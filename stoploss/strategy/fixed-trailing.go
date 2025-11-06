// Copyright 2025 Perry. All rights reserved.

// Licensed MIT License

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
	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/stoploss"
)

type trailing struct {
	stoploss.BaseResolver
	tolerancePct decimal.Decimal
	lastPrice    decimal.Decimal
	threshold    decimal.Decimal
}

// NewFixedTrailingStop creates a FixedTrailingStopLoss based on fixed percentage
func NewFixedTrailingStop(entryPrice, stopLossRate decimal.Decimal, callback stoploss.DefaultCallback) (stoploss.FixedStopLoss, error) {
	if stopLossRate.IsNegative() || stopLossRate.GreaterThan(decimal.NewFromInt(1)) {
		return nil, errStopLossRateInvalid
	}
	return &trailing{
		tolerancePct: stopLossRate,
		lastPrice:    entryPrice,
		threshold:    entryPrice.Mul(decimal.NewFromInt(1).Sub(stopLossRate)),
		BaseResolver: stoploss.BaseResolver{
			Active:   true,
			Callback: callback,
		},
	}, nil
}

// NewFixedTrailingProfit creates a FixedTrailingTakeProfit based on fixed percentage
func NewFixedTrailingProfit(entryPrice, takeProfitRate decimal.Decimal, callback stoploss.DefaultCallback) (stoploss.FixedTakeProfit, error) {
	if takeProfitRate.IsNegative() || takeProfitRate.GreaterThan(decimal.NewFromInt(1)) {
		return nil, errTakeProfitRateInvalid
	}
	return &trailing{
		tolerancePct: takeProfitRate,
		lastPrice:    entryPrice,
		threshold:    entryPrice.Mul(decimal.NewFromInt(1).Add(takeProfitRate)),
		BaseResolver: stoploss.BaseResolver{
			Active:   true,
			Callback: callback,
		},
	}, nil
}

// CalculateStopLoss represents first update last price and calculate stop loss then update threshold
func (t *trailing) CalculateStopLoss(currentPrice decimal.Decimal) (decimal.Decimal, error) {
	if !t.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	if currentPrice.GreaterThan(t.lastPrice) {
		t.threshold = currentPrice.Mul(decimal.NewFromInt(1).Sub(t.tolerancePct))
	}
	t.lastPrice = currentPrice
	return t.threshold, nil
}

// CalculateTakeProfit represents first update last price and calculate take profit then update threshold
func (t *trailing) CalculateTakeProfit(currentPrice decimal.Decimal) (decimal.Decimal, error) {
	if !t.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	if currentPrice.LessThan(t.lastPrice) {
		t.threshold = currentPrice.Mul(decimal.NewFromInt(1).Add(t.tolerancePct))
	}
	t.lastPrice = currentPrice
	return t.threshold, nil
}

// ShouldTriggerStopLoss checks if the stop loss should be triggered
func (t *trailing) ShouldTriggerStopLoss(currentPrice decimal.Decimal) (bool, error) {
	if !t.Active {
		return false, stoploss.ErrStatusInvalid
	}
	if currentPrice.LessThanOrEqual(t.threshold) {
		err := t.Trigger(stoploss.TRIGGERED_REASON_FIXED_TRAILING_STOPLOSS)
		if err != nil {
			return true, stoploss.ErrCallBackFail
		}
		return true, nil
	}
	return false, nil
}

// ShouldTriggerTakeProfit checks if the take profit should be triggered
func (t *trailing) ShouldTriggerTakeProfit(currentPrice decimal.Decimal) (bool, error) {
	if !t.Active {
		return false, stoploss.ErrStatusInvalid
	}
	if currentPrice.GreaterThanOrEqual(t.threshold) {
		err := t.Trigger(stoploss.TRIGGERED_REASON_FIXED_TRAILING_TAKEPROFIT)
		if err != nil {
			return true, stoploss.ErrCallBackFail
		}
		return true, nil
	}
	return false, nil
}

// GetStopLoss returns the current stop loss threshold
func (t *trailing) GetStopLoss() (decimal.Decimal, error) {
	if !t.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return t.threshold, nil
}

// GetTakeProfit returns the current take profit threshold
func (t *trailing) GetTakeProfit() (decimal.Decimal, error) {
	if !t.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return t.threshold, nil
}

// ReSetStopLosser resets the stop loss based on the current price
func (t *trailing) ReSetStopLosser(currentPrice decimal.Decimal) error {
	if !t.Active {
		return stoploss.ErrStatusInvalid
	}
	t.lastPrice = currentPrice
	t.threshold = currentPrice.Mul(decimal.NewFromInt(1).Sub(t.tolerancePct))
	t.Active = true
	return nil
}

// ReSetTakeProfiter resets the take profit based on the current price
func (t *trailing) ReSetTakeProfiter(currentPrice decimal.Decimal) error {
	if !t.Active {
		return stoploss.ErrStatusInvalid
	}
	t.lastPrice = currentPrice
	t.threshold = currentPrice.Mul(decimal.NewFromInt(1).Add(t.tolerancePct))
	t.Active = true
	return nil
}
