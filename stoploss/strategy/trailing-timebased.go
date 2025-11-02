package strategy

import (
	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/stoploss"
)

type trailingTimeBased struct {
	stopLossRate  decimal.Decimal
	lastAdjusted  decimal.Decimal
	stopLoss      decimal.Decimal
	timeThreshold int64
	triggerTime   int64
	callback      DefaultCallback
}

func NewTrailingTimeBased(entryPrice, stopLossRate decimal.Decimal, timeThreshold int64, callback DefaultCallback) stoploss.PeriodicStopLoss {
	return &trailingTimeBased{
		stopLossRate:  stopLossRate,
		callback:      callback,
		lastAdjusted:  entryPrice,
		stopLoss:      entryPrice.Mul(decimal.NewFromInt(1).Sub(stopLossRate)),
		timeThreshold: timeThreshold,
	}
}

func (t *trailingTimeBased) CalculateStopLoss(currentPrice decimal.Decimal) decimal.Decimal {
	if currentPrice.GreaterThan(t.lastAdjusted) {
		t.lastAdjusted = currentPrice
		t.stopLoss = currentPrice.Mul(decimal.NewFromInt(1).Sub(t.stopLossRate))
		t.triggerTime = 0
	}
	return t.stopLoss
}

func (t *trailingTimeBased) ShouldTriggerStopLoss(currentPrice decimal.Decimal, currentTimestamp int64) bool {
	if currentPrice.LessThanOrEqual(t.stopLoss) {
		if t.triggerTime == 0 {
			t.triggerTime = currentTimestamp
		} else if currentTimestamp-t.triggerTime >= t.timeThreshold {
			t.OnTriggered("Trailing Time-Based Stop Loss Triggered")
			return true
		}
	} else {
		t.triggerTime = 0
	}
	return false
}

func (t *trailingTimeBased) GetStopLoss() decimal.Decimal { return t.stopLoss }

func (t *trailingTimeBased) ReSet(currentPrice decimal.Decimal) {
	t.lastAdjusted = currentPrice
	t.stopLoss = currentPrice.Mul(decimal.NewFromInt(1).Sub(t.stopLossRate))
	t.triggerTime = 0
}

func (t *trailingTimeBased) OnTriggered(reason string) {
	if t.callback != nil {
		t.callback(reason)
	}
}
