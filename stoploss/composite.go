package stoploss

import (
	"github.com/shopspring/decimal"
)

type TriggerMode int

const (
	TriggerAny TriggerMode = iota
	TriggerAll
)

type CompositeStopLoss struct {
	mode        TriggerMode
	conditionsT []StopLossCondT
	conditions  []StopLossCond
}

func NewCompositeStopLoss(mode TriggerMode) *CompositeStopLoss {
	return &CompositeStopLoss{
		mode: mode,
	}
}

func (c *CompositeStopLoss) AddTimedCondition(cond StopLossCondT) {
	c.conditionsT = append(c.conditionsT, cond)
}

func (c *CompositeStopLoss) AddCondition(cond StopLossCond) {
	c.conditions = append(c.conditions, cond)
}

func (c *CompositeStopLoss) ShouldTriggerStopLoss(currentPrice decimal.Decimal, timestamp int64) bool {
	count := 0
	total := len(c.conditions) + len(c.conditionsT)

	for _, cond := range c.conditions {
		if triggered, _ := cond.ShouldTriggerStopLoss(currentPrice); triggered {
			if c.mode == TriggerAny {
				return true
			}
			count++
		} else if c.mode == TriggerAll {
			return false
		}
	}

	for _, cond := range c.conditionsT {
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

func (c *CompositeStopLoss) GetMinStopLoss() (decimal.Decimal, error) {
	return c.getStopLoss(func(a, b decimal.Decimal) bool {
		return a.LessThan(b)
	})
}

func (c *CompositeStopLoss) GetMaxStopLoss() (decimal.Decimal, error) {
	return c.getStopLoss(func(a, b decimal.Decimal) bool {
		return a.GreaterThan(b)
	})
}
func (c *CompositeStopLoss) ReSet(currentPrice decimal.Decimal) {
	for _, cond := range c.conditions {
		err := cond.ReSet(currentPrice)
		if err != nil {
			// Handle error
		}
	}
	for _, cond := range c.conditionsT {
		err := cond.ReSet(currentPrice)
		if err != nil {
			// Handle error
		}
	}
}

func (c *CompositeStopLoss) getStopLoss(
	cmp func(a, b decimal.Decimal) bool,
) (decimal.Decimal, error) {
	var result decimal.Decimal
	init := false

	for _, cond := range c.conditions {
		stop, err := cond.GetStopLoss()
		if err != nil {
			continue
		}
		if !init || cmp(stop, result) {
			result = stop
			init = true
		}
	}

	for _, cond := range c.conditionsT {
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
