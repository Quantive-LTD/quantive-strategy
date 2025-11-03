package strategy

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/stoploss"
)

var (
	errStopLossRateInvalid = errors.New("stop loss rate must be between 0 and 1")
)

type fixedPercentStopLoss struct {
	stoploss.BaseStopLoss
	stopLoss    decimal.Decimal
	stopLossPct decimal.Decimal
	LastPrice   decimal.Decimal
}

func NewFixedPercentStopLoss(entryPrice, stopLossPct decimal.Decimal, callback stoploss.DefaultCallback) (stoploss.FixedStopLoss, error) {
	if stopLossPct.IsNegative() || stopLossPct.GreaterThan(decimal.NewFromInt(1)) {
		return nil, errStopLossRateInvalid
	}
	s := &fixedPercentStopLoss{
		LastPrice:   entryPrice,
		stopLossPct: stopLossPct,
		stopLoss:    entryPrice.Mul(decimal.NewFromInt(1).Sub(stopLossPct)),
		BaseStopLoss: stoploss.BaseStopLoss{
			Active:   true,
			Callback: callback,
		},
	}
	return s, nil
}

func (f *fixedPercentStopLoss) CalculateStopLoss(entryPrice decimal.Decimal) (decimal.Decimal, error) {
	if !f.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return f.stopLoss, nil
}

func (f *fixedPercentStopLoss) ShouldTriggerStopLoss(currentPrice decimal.Decimal) (bool, error) {
	if !f.Active {
		return false, stoploss.ErrStatusInvalid
	}
	if currentPrice.LessThanOrEqual(f.stopLoss) {
		err := f.Trigger("Fixed Percent Stop Loss Triggered")
		if err != nil {
			return true, stoploss.ErrCallBackFail
		}
		return true, nil
	}
	return false, nil
}

func (f *fixedPercentStopLoss) GetStopLoss() (decimal.Decimal, error) {
	if !f.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return f.stopLoss, nil
}

func (f *fixedPercentStopLoss) ReSet(currentPrice decimal.Decimal) error {
	if !f.Active {
		return stoploss.ErrStatusInvalid
	}
	f.LastPrice = currentPrice
	f.stopLoss = currentPrice.Mul(decimal.NewFromInt(1).Sub(f.stopLossPct))
	f.Active = true
	return nil
}
