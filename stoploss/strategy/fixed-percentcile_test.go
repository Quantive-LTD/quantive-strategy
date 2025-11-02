package strategy

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestFixed(t *testing.T) {
	entryPrice := decimal.NewFromFloat(100.0)
	stopLossPct := decimal.NewFromFloat(0.1) // 10%
	fixedStopLoss, err := NewFixedPercentStopLoss(entryPrice, stopLossPct, func(reason string) {
		t.Log(reason)
	})
	if err != nil {
		t.Fatalf("Failed to create FixedPercentStopLoss: %v", err)
	}
	// Simulate price movements
	prices := []decimal.Decimal{
		decimal.NewFromFloat(95.0),
		decimal.NewFromFloat(105.0),
		decimal.NewFromFloat(90.0),
	}
	for _, price := range prices {
		stopLoss := fixedStopLoss.CalculateStopLoss(price)
		t.Logf("Current Price: %s, Calculated Stop Loss: %s", price.String(), stopLoss.String())
		if fixedStopLoss.ShouldTriggerStopLoss(price) {
			t.Logf("Stop loss triggered at: %s", price.String())
		}
	}
}
