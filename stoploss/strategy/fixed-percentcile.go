package strategy

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/stoploss"
)

var (
	ErrInvalidStopLossPercentage = errors.New("invalid stop loss percentage")
)

type fixedPercentStopLoss struct {
	LastPrice   decimal.Decimal
	stopLossPct decimal.Decimal
	stopLoss    decimal.Decimal
	triggered   bool
	callback    DefaultCallback
}

func NewFixedPercentStopLoss(entryPrice, stopLossPct decimal.Decimal, callback DefaultCallback) (stoploss.FixedStopLoss, error) {
	if stopLossPct.IsNegative() || stopLossPct.GreaterThan(decimal.NewFromInt(1)) {
		return nil, ErrInvalidStopLossPercentage
	}
	s := &fixedPercentStopLoss{
		LastPrice:   entryPrice,
		stopLossPct: stopLossPct,
		callback:    callback,
	}
	s.stopLoss = entryPrice.Mul(decimal.NewFromInt(1).Sub(stopLossPct))
	return s, nil
}

func (f *fixedPercentStopLoss) CalculateStopLoss(entryPrice decimal.Decimal) decimal.Decimal {
	return f.stopLoss
}

func (f *fixedPercentStopLoss) ShouldTriggerStopLoss(currentPrice decimal.Decimal) bool {
	if currentPrice.LessThanOrEqual(f.stopLoss) {
		f.OnTriggered("Fixed Percent Stop Loss Triggered")
		return true
	}
	return false
}

func (f *fixedPercentStopLoss) GetStopLoss() decimal.Decimal { return f.stopLoss }

func (f *fixedPercentStopLoss) ReSet(currentPrice decimal.Decimal) {
	f.LastPrice = currentPrice
	f.stopLoss = currentPrice.Mul(decimal.NewFromInt(1).Sub(f.stopLossPct))
}

func (f *fixedPercentStopLoss) OnTriggered(reason string) {
	if f.callback != nil {
		f.callback(reason)
	}
}
