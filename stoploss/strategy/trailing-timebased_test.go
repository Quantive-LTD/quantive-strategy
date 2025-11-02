package strategy

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestTrailingTimeWithThreshold(t *testing.T) {
	entryPrice := decimal.NewFromFloat(100)
	stopLossRate := decimal.NewFromFloat(0.05)
	timeThreshold := int64(1)

	sl := NewTrailingTimeBased(entryPrice, stopLossRate, timeThreshold, func(reason string) {
		t.Log(reason)
	})

	priceSeries := []float64{
		101,
		102,
		103,
		97,
		95,
		98,
		90,
		85,
		70,
	}

	simTime := int64(0)
	for _, p := range priceSeries {
		simTime++
		currentPrice := decimal.NewFromFloat(p)
		sl.CalculateStopLoss(currentPrice)
		t.Logf("Current Price: %s, Calculated Stop Loss: %s", currentPrice.String(), sl.GetStopLoss().String())
		if sl.ShouldTriggerStopLoss(currentPrice, simTime) {
			t.Logf("StopLoss triggered at price %s at simTime %d", currentPrice.String(), simTime)
			break
		}
	}

}
