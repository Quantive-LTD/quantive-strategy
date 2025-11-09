// Copyright (C) 2025 Quantive
//
// SPDX-License-Identifier: MIT OR AGPL-3.0-or-later
//
// This file is part of the Decision Engine project.
// You may choose to use this file under the terms of either
// the MIT License or the GNU Affero General Public License v3.0 or later.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the LICENSE files for more details.

package strategy

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/stoploss"
)

var (
	errTimeThresholdInvalid = errors.New("time threshold must be greater than zero")
)

// TrailingTimeBasedStop represents a time-based trailing stop loss strategy
type TrailingTimeBasedStop struct {
	stoploss.BaseResolver
	tolerancePct   decimal.Decimal
	lastPrice      decimal.Decimal
	priceThreshold decimal.Decimal
	timeThreshold  int64
	triggerTime    int64
}

// TrailingTimeBasedProfit represents a time-based trailing take profit strategy
type TrailingTimeBasedProfit struct {
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
	return &TrailingTimeBasedStop{
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
	return &TrailingTimeBasedProfit{
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
func (t *TrailingTimeBasedStop) CalculateStopLoss(currentPrice decimal.Decimal) (decimal.Decimal, error) {
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
func (t *TrailingTimeBasedProfit) CalculateTakeProfit(currentPrice decimal.Decimal) (decimal.Decimal, error) {
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
func (t *TrailingTimeBasedStop) ShouldTriggerStopLoss(currentPrice decimal.Decimal, currentTimestamp int64) (bool, error) {
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
func (t *TrailingTimeBasedProfit) ShouldTriggerTakeProfit(currentPrice decimal.Decimal, currentTimestamp int64) (bool, error) {
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
func (t *TrailingTimeBasedStop) GetStopLoss() (decimal.Decimal, error) {
	if !t.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return t.priceThreshold, nil
}

// GetTakeProfit returns the current take profit threshold
func (t *TrailingTimeBasedProfit) GetTakeProfit() (decimal.Decimal, error) {
	if !t.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return t.priceThreshold, nil
}

// GetTimeThreshold returns the time threshold for stop loss
func (t *TrailingTimeBasedStop) GetTimeThreshold() (int64, error) {
	if !t.Active {
		return 0, stoploss.ErrStatusInvalid
	}
	return t.timeThreshold, nil
}

// GetTimeThreshold returns the time threshold for take profit
func (t *TrailingTimeBasedProfit) GetTimeThreshold() (int64, error) {
	if !t.Active {
		return 0, stoploss.ErrStatusInvalid
	}
	return t.timeThreshold, nil
}

// ReSetStopLosser resets the stop loss based on the current price
func (t *TrailingTimeBasedStop) ReSetStopLosser(currentPrice decimal.Decimal) error {
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
func (t *TrailingTimeBasedProfit) ReSetTakeProfiter(currentPrice decimal.Decimal) error {
	if !t.Active {
		return stoploss.ErrStatusInvalid
	}
	t.lastPrice = currentPrice
	t.priceThreshold = currentPrice.Mul(decimal.NewFromInt(1).Add(t.tolerancePct))
	t.triggerTime = 0
	t.Active = true
	return nil
}
