package strategy

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestNewTrailing(t *testing.T) {
	entryPrice := decimal.NewFromFloat(100)
	stopLossRate := decimal.NewFromFloat(0.1)
	trailing := NewTrailing(entryPrice, stopLossRate, func(reason string) {
		t.Log(reason)
	})
	if trailing == nil {
		t.Fatal("expected non-nil trailing stop loss")
	}
	t.Logf("Initial Stop Loss: %s", trailing.GetStopLoss().String())
	priceSeries := []float64{
		105,
		110,
		108,
		95,
		90,
		92,
		85,
	}
	for _, p := range priceSeries {
		currentPrice := decimal.NewFromFloat(p)
		trailing.CalculateStopLoss(currentPrice)
		t.Logf("Current Price: %s, Calculated Stop Loss: %s", currentPrice.String(), trailing.GetStopLoss().String())
		if trailing.ShouldTriggerStopLoss(currentPrice) {
			t.Logf("StopLoss triggered at price %s", currentPrice.String())
			break
		}
	}
}
