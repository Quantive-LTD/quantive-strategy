package stoploss

import (
	"errors"

	"github.com/shopspring/decimal"
)

var (
	ErrStatusInvalid = errors.New("stop loss status is invalid")
	ErrCallBackFail  = errors.New("stop loss callback function failed")
)

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
	UpdateATR(currentATR decimal.Decimal) error
}

// Fixed-Moving Average
type MAStopLoss interface {
	StopLoss
	StopLossCond
	SetMA(value decimal.Decimal)
}

// general StopLoss interface
type StopLoss interface {
	CalculateStopLoss(currentPrice decimal.Decimal) (decimal.Decimal, error)
	Trigger(reason string) error
	Deactivate() error
}

// StopLoss Condition with timestamp
type StopLossCondT interface {
	ShouldTriggerStopLoss(currentPrice decimal.Decimal, timestamp int64) (bool, error)
	GetStopLoss() (decimal.Decimal, error)
	ReSet(currentPrice decimal.Decimal) error
}

// StopLoss Condition
type StopLossCond interface {
	ShouldTriggerStopLoss(currentPrice decimal.Decimal) (bool, error)
	GetStopLoss() (decimal.Decimal, error)
	ReSet(currentPrice decimal.Decimal) error
}

type DefaultCallback func(reason string) error

type BaseStopLoss struct {
	Active   bool
	Callback DefaultCallback
}

func (b *BaseStopLoss) Deactivate() error {
	if !b.Active {
		return ErrStatusInvalid
	}
	b.Active = false
	return nil
}

func (b *BaseStopLoss) Trigger(reason string) error {
	if !b.Active {
		return ErrStatusInvalid
	}
	if b.Callback != nil {
		return b.Callback(reason)
	}
	return nil
}
