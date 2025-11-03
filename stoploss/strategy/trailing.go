package strategy

import (
	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/stoploss"
)

type trailing struct {
	stoploss.BaseStopLoss
	stopLossRate decimal.Decimal
	lastAdjusted decimal.Decimal
	stopLoss     decimal.Decimal
}

func NewTrailing(entryPrice, stopLossRate decimal.Decimal, callback stoploss.DefaultCallback) (stoploss.FixedStopLoss, error) {
	if stopLossRate.IsNegative() || stopLossRate.GreaterThan(decimal.NewFromInt(1)) {
		return nil, errStopLossRateInvalid
	}
	return &trailing{
		stopLossRate: stopLossRate,
		lastAdjusted: entryPrice,
		stopLoss:     entryPrice.Mul(decimal.NewFromInt(1).Sub(stopLossRate)),
		BaseStopLoss: stoploss.BaseStopLoss{
			Active:   true,
			Callback: callback,
		},
	}, nil
}

func (t *trailing) CalculateStopLoss(currentPrice decimal.Decimal) (decimal.Decimal, error) {
	if !t.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	if currentPrice.GreaterThan(t.lastAdjusted) {
		t.lastAdjusted = currentPrice
		t.stopLoss = currentPrice.Mul(decimal.NewFromInt(1).Sub(t.stopLossRate))
	}
	return t.stopLoss, nil
}

func (t *trailing) ShouldTriggerStopLoss(currentPrice decimal.Decimal) (bool, error) {
	if !t.Active {
		return false, stoploss.ErrStatusInvalid
	}
	if currentPrice.LessThanOrEqual(t.stopLoss) {
		err := t.Trigger("Trailing Stop Loss Triggered")
		if err != nil {
			return true, stoploss.ErrCallBackFail
		}
		return true, nil
	}
	return false, nil
}

func (t *trailing) GetStopLoss() (decimal.Decimal, error) {
	if !t.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return t.stopLoss, nil
}

func (t *trailing) ReSet(currentPrice decimal.Decimal) error {
	if !t.Active {
		return stoploss.ErrStatusInvalid
	}
	t.lastAdjusted = currentPrice
	t.stopLoss = currentPrice.Mul(decimal.NewFromInt(1).Sub(t.stopLossRate))
	t.Active = true
	return nil
}
