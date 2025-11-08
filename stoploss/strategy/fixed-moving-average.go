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
	errMAInvalid     = errors.New("moving average value must be greater than 0")
	errOffsetInvalid = errors.New("offset percentage must be between 0 and 1")
)

// FixedMovingAverageStop represents a moving average based stop loss strategy
type FixedMovingAverageStop struct {
	stoploss.BaseResolver
	threshold     decimal.Decimal
	lastPrice     decimal.Decimal
	movingAverage decimal.Decimal
	offsetPercent decimal.Decimal
}

// FixedMovingAverageProfit represents a moving average based take profit strategy
type FixedMovingAverageProfit struct {
	stoploss.BaseResolver
	threshold     decimal.Decimal
	lastPrice     decimal.Decimal
	movingAverage decimal.Decimal
	offsetPercent decimal.Decimal
}

// NewFixedMovingAverageStop creates a FixedMAStopLoss based on Moving Average
func NewFixedMovingAverageStop(entryPrice, initialMA, offsetPercent decimal.Decimal, callback stoploss.DefaultCallback) (stoploss.FixedMAStopLoss, error) {
	if initialMA.LessThanOrEqual(decimal.Zero) {
		return nil, errMAInvalid
	}
	if offsetPercent.IsNegative() || offsetPercent.GreaterThan(decimal.NewFromInt(1)) {
		return nil, errOffsetInvalid
	}
	return &FixedMovingAverageStop{
		lastPrice:     entryPrice,
		movingAverage: initialMA,
		offsetPercent: offsetPercent,
		threshold:     initialMA.Mul(decimal.NewFromInt(1).Sub(offsetPercent)),
		BaseResolver: stoploss.BaseResolver{
			Active:   true,
			Callback: callback,
		},
	}, nil
}

// NewFixedMovingAverageProfit creates a FixedMATakeProfit based on Moving Average
func NewFixedMovingAverageProfit(entryPrice, initialMA, offsetPercent decimal.Decimal, callback stoploss.DefaultCallback) (stoploss.FixedMATakeProfit, error) {
	if initialMA.LessThanOrEqual(decimal.Zero) {
		return nil, errMAInvalid
	}
	if offsetPercent.IsNegative() || offsetPercent.GreaterThan(decimal.NewFromInt(1)) {
		return nil, errOffsetInvalid
	}
	return &FixedMovingAverageProfit{
		lastPrice:     entryPrice,
		movingAverage: initialMA,
		offsetPercent: offsetPercent,
		threshold:     initialMA.Mul(decimal.NewFromInt(1).Add(offsetPercent)),
		BaseResolver: stoploss.BaseResolver{
			Active:   true,
			Callback: callback,
		},
	}, nil
}

// SetMA sets the current moving average value for stop loss
func (ma *FixedMovingAverageStop) SetMA(value decimal.Decimal) {
	if !ma.Active {
		return
	}
	ma.movingAverage = value
}

// SetMA sets the current moving average value for take profit
func (ma *FixedMovingAverageProfit) SetMA(value decimal.Decimal) {
	if !ma.Active {
		return
	}
	ma.movingAverage = value
}

// CalculateStopLoss represents first update last price and calculate stop loss then update threshold
func (ma *FixedMovingAverageStop) CalculateStopLoss(currentPrice decimal.Decimal) (decimal.Decimal, error) {
	if !ma.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	ma.lastPrice = currentPrice
	ma.threshold = ma.movingAverage.Mul(decimal.NewFromInt(1).Sub(ma.offsetPercent))
	return ma.threshold, nil
}

// CalculateTakeProfit represents first update last price and calculate take profit then update threshold
func (ma *FixedMovingAverageProfit) CalculateTakeProfit(currentPrice decimal.Decimal) (decimal.Decimal, error) {
	if !ma.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	ma.lastPrice = currentPrice
	ma.threshold = ma.movingAverage.Mul(decimal.NewFromInt(1).Add(ma.offsetPercent))
	return ma.threshold, nil
}

// ShouldTriggerStopLoss checks if the stop loss should be triggered
func (ma *FixedMovingAverageStop) ShouldTriggerStopLoss(currentPrice decimal.Decimal) (bool, error) {
	if !ma.Active {
		return false, stoploss.ErrStatusInvalid
	}

	if currentPrice.LessThanOrEqual(ma.threshold) {
		err := ma.Trigger(stoploss.TRIGGERED_REASON_FIXED_MA_STOPLOSS)
		if err != nil {
			return true, stoploss.ErrCallBackFail
		}
		return true, nil
	}
	return false, nil
}

// ShouldTriggerTakeProfit checks if the take profit should be triggered
func (ma *FixedMovingAverageProfit) ShouldTriggerTakeProfit(currentPrice decimal.Decimal) (bool, error) {
	if !ma.Active {
		return false, stoploss.ErrStatusInvalid
	}
	if currentPrice.GreaterThanOrEqual(ma.threshold) {
		err := ma.Trigger(stoploss.TRIGGERED_REASON_FIXED_MA_TAKEPROFIT)
		if err != nil {
			return true, stoploss.ErrCallBackFail
		}
		return true, nil
	}
	return false, nil
}

// GetStopLoss returns the current stop loss threshold
func (ma *FixedMovingAverageStop) GetStopLoss() (decimal.Decimal, error) {
	if !ma.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return ma.threshold, nil
}

// GetTakeProfit returns the current take profit threshold
func (ma *FixedMovingAverageProfit) GetTakeProfit() (decimal.Decimal, error) {
	if !ma.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return ma.threshold, nil
}

// ReSetStopLosser resets the stop loss based on the current price
func (ma *FixedMovingAverageStop) ReSetStopLosser(currentPrice decimal.Decimal) error {
	if !ma.Active {
		return stoploss.ErrStatusInvalid
	}

	ma.lastPrice = currentPrice
	ma.threshold = ma.movingAverage.Mul(decimal.NewFromInt(1).Sub(ma.offsetPercent))
	ma.Active = true
	return nil
}

// ReSetTakeProfiter resets the take profit based on the current price
func (ma *FixedMovingAverageProfit) ReSetTakeProfiter(currentPrice decimal.Decimal) error {
	if !ma.Active {
		return stoploss.ErrStatusInvalid
	}
	ma.lastPrice = currentPrice
	ma.threshold = ma.movingAverage.Mul(decimal.NewFromInt(1).Add(ma.offsetPercent))
	ma.Active = true
	return nil
}
