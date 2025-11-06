// Copyright 2024 Perry. All rights reserved.

// Licensed MIT License

// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// https://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stoploss

import (
	"errors"

	"github.com/shopspring/decimal"
)

var (
	ErrStatusInvalid = errors.New("stop loss status is invalid")
	ErrCallBackFail  = errors.New("stop loss callback function failed")
)

// Time-based
type TimeBasedStopLoss interface {
	StopLoss
	StopLossCondT
}

// Fixed strategy
type FixedStopLoss interface {
	StopLoss
	StopLossCond
}

// Fixed-ATR
type FixedVolatilityStopLoss interface {
	FixedStopLoss
	UpdateATR(currentATR decimal.Decimal) error
}

// Fixed-Moving Average
type FixedMAStopLoss interface {
	FixedStopLoss
	SetMA(value decimal.Decimal)
}

// general StopLoss interface
type StopLoss interface {
	CalculateStopLoss(currentPrice decimal.Decimal) (decimal.Decimal, error)
	Trigger(reason string) error
	ReSetStopLosser(currentPrice decimal.Decimal) error
	Deactivate() error
}

// StopLoss Condition with timestamp
type StopLossCondT interface {
	ShouldTriggerStopLoss(currentPrice decimal.Decimal, timestamp int64) (bool, error)
	GetStopLoss() (decimal.Decimal, error)
}

// StopLoss Condition
type StopLossCond interface {
	ShouldTriggerStopLoss(currentPrice decimal.Decimal) (bool, error)
	GetStopLoss() (decimal.Decimal, error)
}

type DefaultCallback func(reason string) error

type BaseResolver struct {
	Active   bool
	Callback DefaultCallback
}

func (b *BaseResolver) Deactivate() error {
	if !b.Active {
		return ErrStatusInvalid
	}
	b.Active = false
	return nil
}

func (b *BaseResolver) Trigger(reason string) error {
	if !b.Active {
		return ErrStatusInvalid
	}
	if b.Callback != nil {
		return b.Callback(reason)
	}
	return nil
}
