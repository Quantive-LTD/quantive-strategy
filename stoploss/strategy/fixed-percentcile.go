package strategy

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/stoploss"
)

var (
	errStopLossRateInvalid = errors.New("stop loss rate must be between 0 and 1")
)

type fixedPercent struct {
	stoploss.BaseResolver
	stopLoss    decimal.Decimal
	stopLossPct decimal.Decimal
	LastPrice   decimal.Decimal
}

func NewFixedPercent(entryPrice, stopLossPct decimal.Decimal, callback stoploss.DefaultCallback) (stoploss.FixedStopLoss, error) {
	if stopLossPct.IsNegative() || stopLossPct.GreaterThan(decimal.NewFromInt(1)) {
		return nil, errStopLossRateInvalid
	}
	s := &fixedPercent{
		LastPrice:   entryPrice,
		stopLossPct: stopLossPct,
		stopLoss:    entryPrice.Mul(decimal.NewFromInt(1).Sub(stopLossPct)),
		BaseResolver: stoploss.BaseResolver{
			Active:   true,
			Callback: callback,
		},
	}
	return s, nil
}

func (f *fixedPercent) CalculateStopLoss(entryPrice decimal.Decimal) (decimal.Decimal, error) {
	if !f.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return f.stopLoss, nil
}

func (f *fixedPercent) ShouldTriggerStopLoss(currentPrice decimal.Decimal) (bool, error) {
	if !f.Active {
		return false, stoploss.ErrStatusInvalid
	}
	if currentPrice.LessThanOrEqual(f.stopLoss) {
		err := f.Trigger(stoploss.TRIGGERED_REASON_FIXED_PERCENTCILE_STOPLOSS)
		if err != nil {
			return true, stoploss.ErrCallBackFail
		}
		return true, nil
	}
	return false, nil
}

func (f *fixedPercent) GetStopLoss() (decimal.Decimal, error) {
	if !f.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return f.stopLoss, nil
}

func (f *fixedPercent) ReSet(currentPrice decimal.Decimal) error {
	if !f.Active {
		return stoploss.ErrStatusInvalid
	}
	f.LastPrice = currentPrice
	f.stopLoss = currentPrice.Mul(decimal.NewFromInt(1).Sub(f.stopLossPct))
	f.Active = true
	return nil
}
