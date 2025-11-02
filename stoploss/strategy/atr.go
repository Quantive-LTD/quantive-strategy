package strategy

import (
	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/stoploss"
)

type atr struct {
	stop       decimal.Decimal
	entryPrice decimal.Decimal
	multiplier decimal.Decimal
	currentATR decimal.Decimal
	active     bool
	callback   DefaultCallback
}

func NewATRStop(entryPrice, ATR, k decimal.Decimal, callback DefaultCallback) stoploss.VolatilityStopLoss {
	return &atr{
		entryPrice: entryPrice,
		currentATR: ATR,
		multiplier: k,
		stop:       entryPrice.Sub(ATR.Mul(k)),
		active:     true,
		callback:   callback,
	}
}

func (a *atr) UpdateATR(currentATR decimal.Decimal) {
	a.currentATR = currentATR
	a.stop = a.entryPrice.Sub(a.currentATR.Mul(a.multiplier))
}

func (a *atr) CalculateStopLoss(currentPrice decimal.Decimal) decimal.Decimal {
	return a.stop
}

func (a *atr) ShouldTriggerStopLoss(currentPrice decimal.Decimal) bool {
	if currentPrice.LessThanOrEqual(a.stop) {
		a.OnTriggered("ATR Stop Loss Triggered")
		return true
	}
	return false
}

func (a *atr) GetStopLoss() decimal.Decimal {
	return a.stop
}

func (a *atr) ReSet(currentPrice decimal.Decimal) {
	a.entryPrice = currentPrice
	a.stop = currentPrice.Sub(a.currentATR.Mul(a.multiplier))
}

func (a *atr) OnTriggered(reason string) {
	if a.callback != nil {
		a.callback(reason)
	}
}
