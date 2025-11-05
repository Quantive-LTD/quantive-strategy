package stoploss

import "github.com/shopspring/decimal"

type HybridWithoutTime interface {
	TakeProfit
	StopLoss
	TakeProfitCond
	StopLossCond
}

type HybridWithTime interface {
	TakeProfit
	StopLoss
	TakeProfitCondT
	StopLossCondT
}

// Fixed-Percentage | Fixed-Trailing
type FixedTakeProfit struct {
	TakeProfit
	TakeProfitCond
}

// Time-based + Trailing
type PeriodicTakeProfit struct {
	TakeProfit
	TakeProfitCondT
}

// Fixed-ATR
type VolatilityTakeProfit interface {
	TakeProfit
	TakeProfitCond
	UpdateATR(currentATR decimal.Decimal) error
}

// Fixed-Moving Average
type MATakeProfit interface {
	TakeProfit
	TakeProfitCond
	SetMA(value decimal.Decimal)
}

type TakeProfit interface {
	CalculateTakeProfit(currentPrice decimal.Decimal) (decimal.Decimal, error)
	Trigger(reason string) error
	Deactivate() error
}

type TakeProfitCond interface {
	ShouldTriggerTakeProfit(currentPrice decimal.Decimal) (bool, error)
	GetTakeProfit() (decimal.Decimal, error)
	ReSet(currentPrice decimal.Decimal) error
}

type TakeProfitCondT interface {
	ShouldTriggerTakeProfit(currentPrice decimal.Decimal, timestamp int64) (bool, error)
	GetTakeProfit() (decimal.Decimal, error)
	ReSet(currentPrice decimal.Decimal) error
}
