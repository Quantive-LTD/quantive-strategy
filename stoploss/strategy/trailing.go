package strategy

import (
	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/stoploss"
)

type DefaultCallback func(reason string)

type trailing struct {
	stopLossRate decimal.Decimal
	lastAdjusted decimal.Decimal
	stopLoss     decimal.Decimal
	callback     DefaultCallback
}

func NewTrailing(entryPrice, stopLossRate decimal.Decimal, callback DefaultCallback) stoploss.FixedStopLoss {
	return &trailing{
		stopLossRate: stopLossRate,
		lastAdjusted: entryPrice,
		stopLoss:     entryPrice.Mul(decimal.NewFromInt(1).Sub(stopLossRate)),
		callback:     callback,
	}
}

func (t *trailing) CalculateStopLoss(currentPrice decimal.Decimal) decimal.Decimal {
	if currentPrice.GreaterThan(t.lastAdjusted) {
		t.lastAdjusted = currentPrice
		t.stopLoss = currentPrice.Mul(decimal.NewFromInt(1).Sub(t.stopLossRate))
	}
	return t.stopLoss
}

func (t *trailing) ShouldTriggerStopLoss(currentPrice decimal.Decimal) bool {
	if currentPrice.LessThanOrEqual(t.stopLoss) {
		t.OnTriggered("Trailing Stop Loss Triggered")
		return true
	}
	return false
}

func (t *trailing) GetStopLoss() decimal.Decimal { return t.stopLoss }

func (t *trailing) ReSet(currentPrice decimal.Decimal) {
	t.lastAdjusted = currentPrice
	t.stopLoss = currentPrice.Mul(decimal.NewFromInt(1).Sub(t.stopLossRate))
}

func (t *trailing) OnTriggered(reason string) {
	if t.callback != nil {
		t.callback(reason)
	}
}
