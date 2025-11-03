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
	stoploss.BaseStopLoss
	stopLossRate  decimal.Decimal
	lastAdjusted  decimal.Decimal
	stopLoss      decimal.Decimal
	timeThreshold int64
	triggerTime   int64
}

func NewTrailingTimeBased(entryPrice, stopLossRate decimal.Decimal, timeThreshold int64, callback stoploss.DefaultCallback) (stoploss.PeriodicStopLoss, error) {
	if stopLossRate.IsNegative() || stopLossRate.GreaterThan(decimal.NewFromInt(1)) {
		return nil, errStopLossRateInvalid
	}
	if timeThreshold <= 0 {
		return nil, errTimeThresholdInvalid
	}
	return &trailingTimeBased{
		stopLossRate:  stopLossRate,
		lastAdjusted:  entryPrice,
		stopLoss:      entryPrice.Mul(decimal.NewFromInt(1).Sub(stopLossRate)),
		timeThreshold: timeThreshold,
		BaseStopLoss: stoploss.BaseStopLoss{
			Active:   true,
			Callback: callback,
		},
	}, nil
}

func (t *trailingTimeBased) CalculateStopLoss(currentPrice decimal.Decimal) (decimal.Decimal, error) {
	if !t.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	if currentPrice.GreaterThan(t.lastAdjusted) {
		t.lastAdjusted = currentPrice
		t.stopLoss = currentPrice.Mul(decimal.NewFromInt(1).Sub(t.stopLossRate))
		t.triggerTime = 0
	}
	return t.stopLoss, nil
}

func (t *trailingTimeBased) ShouldTriggerStopLoss(currentPrice decimal.Decimal, currentTimestamp int64) (bool, error) {
	if !t.Active {
		return false, stoploss.ErrStatusInvalid
	}
	if currentPrice.LessThanOrEqual(t.stopLoss) {
		if t.triggerTime == 0 {
			t.triggerTime = currentTimestamp
		} else if currentTimestamp-t.triggerTime >= t.timeThreshold {
			err := t.Trigger("Trailing Time-Based Stop Loss Triggered")
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

func (t *trailingTimeBased) GetStopLoss() (decimal.Decimal, error) {
	if !t.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return t.stopLoss, nil
}

func (t *trailingTimeBased) ReSet(currentPrice decimal.Decimal) error {
	if !t.Active {
		return stoploss.ErrStatusInvalid
	}
	t.lastAdjusted = currentPrice
	t.stopLoss = currentPrice.Mul(decimal.NewFromInt(1).Sub(t.stopLossRate))
	t.triggerTime = 0
	t.Active = true
	return nil
}
