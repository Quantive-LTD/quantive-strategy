package strategy

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/stoploss"
)

var (
	errATRStopLossKInvalid = errors.New("stop loss k must be greater than 0")
)

type atr struct {
	stoploss.BaseStopLoss
	stopLoss   decimal.Decimal
	entryPrice decimal.Decimal
	multiplier decimal.Decimal
	currentATR decimal.Decimal
}

func NewATRStop(entryPrice, ATR, k decimal.Decimal, callback stoploss.DefaultCallback) (stoploss.VolatilityStopLoss, error) {
	if k.LessThanOrEqual(decimal.Zero) {
		return nil, errATRStopLossKInvalid
	}
	return &atr{
		entryPrice: entryPrice,
		currentATR: ATR,
		multiplier: k,
		stopLoss:   entryPrice.Sub(ATR.Mul(k)),
		BaseStopLoss: stoploss.BaseStopLoss{
			Active:   true,
			Callback: callback,
		},
	}, nil
}

func (a *atr) UpdateATR(currentATR decimal.Decimal) error {
	if !a.Active {
		return stoploss.ErrStatusInvalid
	}
	a.currentATR = currentATR
	a.stopLoss = a.entryPrice.Sub(a.currentATR.Mul(a.multiplier))
	return nil
}

func (a *atr) CalculateStopLoss(currentPrice decimal.Decimal) (decimal.Decimal, error) {
	if !a.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return a.stopLoss, nil
}

func (a *atr) ShouldTriggerStopLoss(currentPrice decimal.Decimal) (bool, error) {
	if !a.Active {
		return false, stoploss.ErrStatusInvalid
	}
	if currentPrice.LessThanOrEqual(a.stopLoss) {
		err := a.Trigger("ATR Stop Loss Triggered")
		if err != nil {
			return true, stoploss.ErrCallBackFail
		}
		return true, nil
	}
	return false, nil
}

func (a *atr) GetStopLoss() (decimal.Decimal, error) {
	if !a.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return a.stopLoss, nil
}

func (a *atr) ReSet(currentPrice decimal.Decimal) error {
	if !a.Active {
		return stoploss.ErrStatusInvalid
	}
	a.entryPrice = currentPrice
	a.stopLoss = currentPrice.Sub(a.currentATR.Mul(a.multiplier))
	a.Active = true
	return nil
}
