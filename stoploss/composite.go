package stoploss

import "github.com/shopspring/decimal"

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
	triggered := 0
	total := len(c.conditions) + len(c.conditionsT)

	for _, cond := range c.conditions {
		if cond.ShouldTriggerStopLoss(currentPrice) {
			if c.mode == TriggerAny {
				return true
			}
			triggered++
		} else if c.mode == TriggerAll {
			return false
		}
	}

	for _, cond := range c.conditionsT {
		if cond.ShouldTriggerStopLoss(currentPrice, timestamp) {
			if c.mode == TriggerAny {
				return true
			}
			triggered++
		} else if c.mode == TriggerAll {
			return false
		}
	}

	return c.mode == TriggerAll && triggered == total
}

func (c *CompositeStopLoss) GetStopLoss() decimal.Decimal {
	var min decimal.Decimal
	init := false

	for _, cond := range c.conditions {
		stop := cond.GetStopLoss()
		if !init || stop.LessThan(min) {
			min = stop
			init = true
		}
	}

	for _, cond := range c.conditionsT {
		stop := cond.GetStopLoss()
		if !init || stop.LessThan(min) {
			min = stop
			init = true
		}
	}

	if !init {
		return decimal.Zero
	}
	return min
}
