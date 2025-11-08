// Copyright 2025 Quantive. All rights reserved.

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
	"github.com/shopspring/decimal"
)

type TriggerMode int

const (
	TriggerAny TriggerMode = iota
	TriggerAll
)

type Composite struct {
	mode      TriggerMode
	stopCondT []StopLossCondT
	stopCond  []StopLossCond
	takeCondT []TakeProfitCondT
	takeCond  []TakeProfitCond
}

func NewComposite(mode TriggerMode) *Composite {
	return &Composite{
		mode: mode,
	}
}

func (c *Composite) AddTimedCondition(cond1 StopLossCondT, cond2 TakeProfitCondT) {
	c.stopCondT = append(c.stopCondT, cond1)
	c.takeCondT = append(c.takeCondT, cond2)
}

func (c *Composite) AddCondition(cond1 StopLossCond, cond2 TakeProfitCond) {
	c.stopCond = append(c.stopCond, cond1)
	c.takeCond = append(c.takeCond, cond2)
}

func (c *Composite) ShouldTriggerStopLoss(currentPrice decimal.Decimal, timestamp int64) bool {
	count := 0
	total := len(c.stopCond) + len(c.stopCondT)

	for _, cond := range c.stopCond {
		if triggered, _ := cond.ShouldTriggerStopLoss(currentPrice); triggered {
			if c.mode == TriggerAny {
				return true
			}
			count++
		} else if c.mode == TriggerAll {
			return false
		}
	}

	for _, cond := range c.stopCondT {
		if triggered, _ := cond.ShouldTriggerStopLoss(currentPrice, timestamp); triggered {
			if c.mode == TriggerAny {
				return true
			}
			count++
		} else if c.mode == TriggerAll {
			return false
		}
	}
	return c.mode == TriggerAll && count == total
}

func (c *Composite) ShouldTriggerTakeProfit(currentPrice decimal.Decimal, timestamp int64) bool {
	count := 0
	total := len(c.takeCond) + len(c.takeCondT)
	for _, cond := range c.takeCond {
		if triggered, _ := cond.ShouldTriggerTakeProfit(currentPrice); triggered {
			if c.mode == TriggerAny {
				return true
			}
			count++
		} else if c.mode == TriggerAll {
			return false
		}
	}

	for _, cond := range c.takeCondT {
		if triggered, _ := cond.ShouldTriggerTakeProfit(currentPrice, timestamp); triggered {
			if c.mode == TriggerAny {
				return true
			}
			count++
		} else if c.mode == TriggerAll {
			return false
		}
	}
	return c.mode == TriggerAll && count == total
}

func (c *Composite) GetMinTakeProfit() (decimal.Decimal, error) {
	return c.getTakeProfit(func(a, b decimal.Decimal) bool {
		return a.LessThan(b)
	})
}

func (c *Composite) GetMaxTakeProfit() (decimal.Decimal, error) {
	return c.getTakeProfit(func(a, b decimal.Decimal) bool {
		return a.GreaterThan(b)
	})
}

func (c *Composite) GetMinStopLoss() (decimal.Decimal, error) {
	return c.getStopLoss(func(a, b decimal.Decimal) bool {
		return a.LessThan(b)
	})
}

func (c *Composite) GetMaxStopLoss() (decimal.Decimal, error) {
	return c.getStopLoss(func(a, b decimal.Decimal) bool {
		return a.GreaterThan(b)
	})
}

func (c *Composite) ResetMode(mode TriggerMode) {
	c.mode = mode
}

func (c *Composite) getStopLoss(
	cmp func(a, b decimal.Decimal) bool,
) (decimal.Decimal, error) {
	var result decimal.Decimal
	init := false

	for _, cond := range c.stopCond {
		stop, err := cond.GetStopLoss()
		if err != nil {
			continue
		}
		if !init || cmp(stop, result) {
			result = stop
			init = true
		}
	}

	for _, cond := range c.stopCondT {
		stop, err := cond.GetStopLoss()
		if err != nil {
			continue
		}
		if !init || cmp(stop, result) {
			result = stop
			init = true
		}
	}

	if !init {
		return decimal.Zero, ErrStatusInvalid
	}
	return result, nil
}

func (c *Composite) getTakeProfit(
	cmp func(a, b decimal.Decimal) bool,
) (decimal.Decimal, error) {
	var result decimal.Decimal
	init := false

	for _, cond := range c.takeCond {
		take, err := cond.GetTakeProfit()
		if err != nil {
			continue
		}
		if !init || cmp(take, result) {
			result = take
			init = true
		}
	}

	for _, cond := range c.takeCondT {
		take, err := cond.GetTakeProfit()
		if err != nil {
			continue
		}
		if !init || cmp(take, result) {
			result = take
			init = true
		}
	}

	if !init {
		return decimal.Zero, ErrStatusInvalid
	}
	return result, nil
}
