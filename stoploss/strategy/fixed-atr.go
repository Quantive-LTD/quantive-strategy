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
	errATRStopLossKInvalid = errors.New("stop loss k must be greater than 0")
	errATRInvalid          = errors.New("ATR must be greater than 0")
)

type ATR struct {
	stoploss.BaseResolver
	threshold  decimal.Decimal
	lastPrice  decimal.Decimal
	multiplier decimal.Decimal
	currentATR decimal.Decimal
}

// NewFixedATRStop creates a FixedVolatilityStopLoss based on ATR
func NewFixedATRStop(entryPrice, atr, k decimal.Decimal, callback stoploss.DefaultCallback) (stoploss.FixedVolatilityStopLoss, error) {
	if k.LessThanOrEqual(decimal.Zero) {
		return nil, errATRStopLossKInvalid
	}
	if atr.LessThanOrEqual(decimal.Zero) {
		return nil, errATRInvalid
	}

	return &ATR{
		lastPrice:  entryPrice,
		currentATR: atr,
		multiplier: k,
		threshold:  entryPrice.Sub(atr.Mul(k)),
		BaseResolver: stoploss.BaseResolver{
			Active:   true,
			Callback: callback,
		},
	}, nil
}

// NewFixedATRProfit creates a VolatilityTakeProfit based on ATR
func NewFixedATRProfit(entryPrice, atr, k decimal.Decimal, callback stoploss.DefaultCallback) (stoploss.FixedVolatilityTakeProfit, error) {
	if k.LessThanOrEqual(decimal.Zero) {
		return nil, errATRStopLossKInvalid
	}
	if atr.LessThanOrEqual(decimal.Zero) {
		return nil, errATRInvalid
	}
	return &ATR{
		lastPrice:  entryPrice,
		currentATR: atr,
		multiplier: k,
		threshold:  entryPrice.Add(atr.Mul(k)),
		BaseResolver: stoploss.BaseResolver{
			Active:   true,
			Callback: callback,
		},
	}, nil
}

// UpdateATR updates the current ATR value
func (a *ATR) UpdateATR(currentATR decimal.Decimal) error {
	if !a.Active {
		return stoploss.ErrStatusInvalid
	}
	a.currentATR = currentATR
	return nil
}

// CalculateStopLoss represents first update last price and calculate stop loss then update threshold
func (a *ATR) CalculateStopLoss(currentPrice decimal.Decimal) (decimal.Decimal, error) {
	if !a.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	a.lastPrice = currentPrice
	a.threshold = currentPrice.Sub(a.currentATR.Mul(a.multiplier))
	return a.threshold, nil
}

// CalculateTakeProfit represents first update last price and calculate take profit then update threshold
func (a *ATR) CalculateTakeProfit(currentPrice decimal.Decimal) (decimal.Decimal, error) {
	if !a.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	a.lastPrice = currentPrice
	a.threshold = currentPrice.Add(a.currentATR.Mul(a.multiplier))
	return a.threshold, nil
}

// ShouldTriggerStopLoss checks if the stop loss should be triggered
func (a *ATR) ShouldTriggerStopLoss(currentPrice decimal.Decimal) (bool, error) {
	if !a.Active {
		return false, stoploss.ErrStatusInvalid
	}
	if currentPrice.LessThanOrEqual(a.threshold) {
		err := a.Trigger(stoploss.TRIGGERED_REASON_FIXED_ATR_STOPLOSS)
		if err != nil {
			return true, stoploss.ErrCallBackFail
		}
		return true, nil
	}
	return false, nil
}

// ShouldTriggerTakeProfit checks if the take profit should be triggered
func (a *ATR) ShouldTriggerTakeProfit(currentPrice decimal.Decimal) (bool, error) {
	if !a.Active {
		return false, stoploss.ErrStatusInvalid
	}
	if currentPrice.GreaterThanOrEqual(a.threshold) {
		err := a.Trigger(stoploss.TRIGGERED_REASON_FIXED_ATR_TAKEPROFIT)
		if err != nil {
			return true, stoploss.ErrCallBackFail
		}
		return true, nil
	}
	return false, nil
}

// GetStopLoss returns the current stop loss threshold
func (a *ATR) GetStopLoss() (decimal.Decimal, error) {
	if !a.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return a.threshold, nil
}

// GetTakeProfit returns the current take profit threshold
func (a *ATR) GetTakeProfit() (decimal.Decimal, error) {
	if !a.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return a.threshold, nil
}

// ReSetStopLosser resets the stop loss based on the current price
func (a *ATR) ReSetStopLosser(currentPrice decimal.Decimal) error {
	if !a.Active {
		return stoploss.ErrStatusInvalid
	}
	a.lastPrice = currentPrice
	a.threshold = currentPrice.Sub(a.currentATR.Mul(a.multiplier))
	a.Active = true
	return nil
}

// ReSetTakeProfiter resets the take profit based on the current price
func (a *ATR) ReSetTakeProfiter(currentPrice decimal.Decimal) error {
	if !a.Active {
		return stoploss.ErrStatusInvalid
	}
	a.lastPrice = currentPrice
	a.threshold = currentPrice.Add(a.currentATR.Mul(a.multiplier))
	a.Active = true
	return nil
}
