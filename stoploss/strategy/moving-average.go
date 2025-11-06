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

type movingAverage struct {
	stoploss.BaseResolver
	stopLoss      decimal.Decimal
	lastPrice     decimal.Decimal
	movingAverage decimal.Decimal
	offsetPercent decimal.Decimal
}

func NewMovingAverageStop(entryPrice, initialMA, offsetPercent decimal.Decimal, callback stoploss.DefaultCallback) (stoploss.MAStopLoss, error) {
	if initialMA.LessThanOrEqual(decimal.Zero) {
		return nil, errMAInvalid
	}
	if offsetPercent.IsNegative() || offsetPercent.GreaterThan(decimal.NewFromInt(1)) {
		return nil, errOffsetInvalid
	}
	return &movingAverage{
		lastPrice:     entryPrice,
		movingAverage: initialMA,
		offsetPercent: offsetPercent,
		stopLoss:      initialMA.Mul(decimal.NewFromInt(1).Sub(offsetPercent)),
		BaseResolver: stoploss.BaseResolver{
			Active:   true,
			Callback: callback,
		},
	}, nil
}

func (ma *movingAverage) SetMA(value decimal.Decimal) {
	if !ma.Active {
		return
	}

	ma.movingAverage = value

	stopLossValue := value.Mul(decimal.NewFromInt(1).Sub(ma.offsetPercent))

	// Only update if new stop loss is higher (for long positions)
	// This prevents stop loss from moving down with the MA
	if stopLossValue.GreaterThan(ma.stopLoss) {
		ma.stopLoss = stopLossValue
	}
}

func (ma *movingAverage) CalculateStopLoss(currentPrice decimal.Decimal) (decimal.Decimal, error) {
	if !ma.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return ma.stopLoss, nil
}

func (ma *movingAverage) ShouldTriggerStopLoss(currentPrice decimal.Decimal) (bool, error) {
	if !ma.Active {
		return false, stoploss.ErrStatusInvalid
	}

	if currentPrice.LessThanOrEqual(ma.stopLoss) {
		err := ma.Trigger(stoploss.TRIGGERED_REASON_FIXED_MA_STOPLOSS)
		if err != nil {
			return true, stoploss.ErrCallBackFail
		}
		return true, nil
	}
	return false, nil
}

func (ma *movingAverage) GetStopLoss() (decimal.Decimal, error) {
	if !ma.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return ma.stopLoss, nil
}

func (ma *movingAverage) ReSet(currentPrice decimal.Decimal) error {
	if !ma.Active {
		return stoploss.ErrStatusInvalid
	}

	ma.lastPrice = currentPrice
	ma.stopLoss = ma.movingAverage.Mul(decimal.NewFromInt(1).Sub(ma.offsetPercent))
	ma.Active = true
	return nil
}
