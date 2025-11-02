package strategy

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
)

func TestATRStopLoss(t *testing.T) {
	callback := func(reason string) {
		fmt.Println("Callback:", reason)
	}

	entry := decimal.NewFromFloat(100)
	atr := decimal.NewFromFloat(2.5)
	mult := decimal.NewFromFloat(3)

	s := NewATRStop(entry, atr, mult, callback)

	fmt.Println("Initial Stop:", s.GetStopLoss()) // = 100 - 2.5*3 = 92.5

	triggered := s.ShouldTriggerStopLoss(decimal.NewFromFloat(91))
	fmt.Println("Triggered:", triggered)
}
