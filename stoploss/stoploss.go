package stoploss

import "github.com/shopspring/decimal"

// Time-based + Trailing
type PeriodicStopLoss interface {
	StopLoss
	StopLossCondT
}

// Fixed-Percentage | Fixed-Trailing
type FixedStopLoss interface {
	StopLoss
	StopLossCond
}

// Fixed-ATR
type VolatilityStopLoss interface {
	StopLoss
	StopLossCond
	UpdateATR(currentATR decimal.Decimal)
}

// Fixed-Moving Average
type MAStopLoss interface {
	StopLoss
	StopLossCond
	SetMA(value decimal.Decimal)
}

// general StopLoss interface
type StopLoss interface {
	CalculateStopLoss(currentPrice decimal.Decimal) decimal.Decimal
	ReSet(currentPrice decimal.Decimal)
	OnTriggered(reason string)
}

// StopLoss Condition with timestamp
type StopLossCondT interface {
	ShouldTriggerStopLoss(currentPrice decimal.Decimal, timestamp int64) bool
	GetStopLoss() decimal.Decimal
}

// StopLoss Condition
type StopLossCond interface {
	ShouldTriggerStopLoss(currentPrice decimal.Decimal) bool
	GetStopLoss() decimal.Decimal
}
