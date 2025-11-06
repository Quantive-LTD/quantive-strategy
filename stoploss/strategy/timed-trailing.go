// Copyright 2024 Perry. All rights reserved.

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
	"errors"

	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/stoploss"
)

var (
	errTimeThresholdInvalid = errors.New("time threshold must be greater than zero")
)

type trailingTimeBased struct {
	stoploss.BaseResolver
	tolerancePct   decimal.Decimal
	lastPrice      decimal.Decimal
	priceThreshold decimal.Decimal
	timeThreshold  int64
	triggerTime    int64
}

// NewTrailingTimeBasedStop creates a TimeBasedTrailingStopLoss based on fixed percentage and time threshold
func NewTrailingTimeBasedStop(entryPrice, stopLossRate decimal.Decimal, timeThreshold int64, callback stoploss.DefaultCallback) (stoploss.TimeBasedStopLoss, error) {
	if stopLossRate.IsNegative() || stopLossRate.GreaterThan(decimal.NewFromInt(1)) {
		return nil, errStopLossRateInvalid
	}
	if timeThreshold <= 0 {
		return nil, errTimeThresholdInvalid
	}
	return &trailingTimeBased{
		tolerancePct:   stopLossRate,
		lastPrice:      entryPrice,
		priceThreshold: entryPrice.Mul(decimal.NewFromInt(1).Sub(stopLossRate)),
		timeThreshold:  timeThreshold,
		BaseResolver: stoploss.BaseResolver{
			Active:   true,
			Callback: callback,
		},
	}, nil
}

// NewTrailingTimeBasedProfit creates a TimeBasedTrailingTakeProfit based on fixed percentage and time threshold
func NewTrailingTimeBasedProfit(entryPrice, takeProfitRate decimal.Decimal, timeThreshold int64, callback stoploss.DefaultCallback) (stoploss.TimeBasedTakeProfit, error) {
	if takeProfitRate.IsNegative() || takeProfitRate.GreaterThan(decimal.NewFromInt(1)) {
		return nil, errTakeProfitRateInvalid
	}
	if timeThreshold <= 0 {
		return nil, errTimeThresholdInvalid
	}
	return &trailingTimeBased{
		tolerancePct:   takeProfitRate,
		lastPrice:      entryPrice,
		priceThreshold: entryPrice.Mul(decimal.NewFromInt(1).Add(takeProfitRate)),
		timeThreshold:  timeThreshold,
		BaseResolver: stoploss.BaseResolver{
			Active:   true,
			Callback: callback,
		},
	}, nil
}

// CalculateStopLoss represents first update last price and calculate stop loss then update threshold
func (t *trailingTimeBased) CalculateStopLoss(currentPrice decimal.Decimal) (decimal.Decimal, error) {
	if !t.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	if currentPrice.GreaterThan(t.lastPrice) {
		t.lastPrice = currentPrice
		t.priceThreshold = currentPrice.Mul(decimal.NewFromInt(1).Sub(t.tolerancePct))
		t.triggerTime = 0
	}
	t.lastPrice = currentPrice
	return t.priceThreshold, nil
}

// CalculateTakeProfit represents first update last price and calculate take profit then update threshold
func (t *trailingTimeBased) CalculateTakeProfit(currentPrice decimal.Decimal) (decimal.Decimal, error) {
	if !t.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	if currentPrice.LessThan(t.lastPrice) {
		t.lastPrice = currentPrice
		t.priceThreshold = currentPrice.Mul(decimal.NewFromInt(1).Add(t.tolerancePct))
		t.triggerTime = 0
	}
	t.lastPrice = currentPrice
	return t.priceThreshold, nil
}

// ShouldTriggerStopLoss checks if the stop loss should be triggered based on time threshold and price threshold
func (t *trailingTimeBased) ShouldTriggerStopLoss(currentPrice decimal.Decimal, currentTimestamp int64) (bool, error) {
	if !t.Active {
		return false, stoploss.ErrStatusInvalid
	}
	if currentPrice.LessThanOrEqual(t.priceThreshold) {
		if t.triggerTime == 0 {
			t.triggerTime = currentTimestamp
		} else if currentTimestamp-t.triggerTime >= t.timeThreshold {
			err := t.Trigger(stoploss.TRIGGERED_REASON_TIMED_TRAILING_STOPLOSS)
			if err != nil {
				return true, stoploss.ErrCallBackFail
			}
			return true, nil
		}
	} else {
		t.triggerTime = 0
	}
	return false, nil
}

// ShouldTriggerTakeProfit checks if the take profit should be triggered based on time threshold and price threshold
func (t *trailingTimeBased) ShouldTriggerTakeProfit(currentPrice decimal.Decimal, currentTimestamp int64) (bool, error) {
	if !t.Active {
		return false, stoploss.ErrStatusInvalid
	}
	if currentPrice.GreaterThanOrEqual(t.priceThreshold) {
		if t.triggerTime == 0 {
			t.triggerTime = currentTimestamp
		} else if currentTimestamp-t.triggerTime >= t.timeThreshold {
			err := t.Trigger(stoploss.TRIGGERED_REASON_TIMED_TRAILING_TAKEPROFIT)
			if err != nil {
				return true, stoploss.ErrCallBackFail
			}
			return true, nil
		}
	} else {
		t.triggerTime = 0
	}
	return false, nil
}

// GetStopLoss returns the current stop loss threshold
func (t *trailingTimeBased) GetStopLoss() (decimal.Decimal, error) {
	if !t.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return t.priceThreshold, nil
}

// GetTakeProfit returns the current take profit threshold
func (t *trailingTimeBased) GetTakeProfit() (decimal.Decimal, error) {
	if !t.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return t.priceThreshold, nil
}

// ReSetStopLosser resets the stop loss based on the current price
func (t *trailingTimeBased) ReSetStopLosser(currentPrice decimal.Decimal) error {
	if !t.Active {
		return stoploss.ErrStatusInvalid
	}
	t.lastPrice = currentPrice
	t.priceThreshold = currentPrice.Mul(decimal.NewFromInt(1).Sub(t.tolerancePct))
	t.triggerTime = 0
	t.Active = true
	return nil
}

// ReSetTakeProfiter resets the take profit based on the current price
func (t *trailingTimeBased) ReSetTakeProfiter(currentPrice decimal.Decimal) error {
	if !t.Active {
		return stoploss.ErrStatusInvalid
	}
	t.lastPrice = currentPrice
	t.priceThreshold = currentPrice.Mul(decimal.NewFromInt(1).Add(t.tolerancePct))
	t.triggerTime = 0
	t.Active = true
	return nil
}
