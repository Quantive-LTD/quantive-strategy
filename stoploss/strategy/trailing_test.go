package strategy

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/stoploss"
)

var trailingCallback = func(reason string) error {
	fmt.Println("Trailing Callback:", reason)
	return nil
}

func TestNewTrailing(t *testing.T) {
	entryPrice := decimal.NewFromFloat(100)
	stopLossRate := decimal.NewFromFloat(0.1)
	trailing, err := NewTrailing(entryPrice, stopLossRate, callback)
	if err != nil {
		t.Fatalf("Failed to create trailing stop loss: %v", err)
	}
	if trailing == nil {
		t.Fatal("Trailing stop loss is nil")
	}
	_, err = trailing.GetStopLoss()
	if err != nil {
		t.Fatalf("Failed to get stop loss: %v", err)
	}
	_, err = trailing.ShouldTriggerStopLoss(decimal.NewFromFloat(95))
	if err != nil {
		t.Fatalf("Failed to check if stop loss should trigger: %v", err)
	}
	_, err = trailing.CalculateStopLoss(decimal.NewFromFloat(110))
	if err != nil {
		t.Fatalf("Failed to calculate stop loss: %v", err)
	}
	err = trailing.ReSet(decimal.NewFromFloat(105))
	if err != nil {
		t.Fatalf("Failed to reset trailing stop loss: %v", err)
	}
}

func TestTrailingStopLoss_NewTrailing(t *testing.T) {
	tests := []struct {
		name         string
		entryPrice   decimal.Decimal
		stopLossRate decimal.Decimal
		expectError  bool
		errorType    error
	}{
		{
			name:         "Valid 5% trailing stop",
			entryPrice:   d(100),
			stopLossRate: d(0.05),
			expectError:  false,
		},
		{
			name:         "Valid 10% trailing stop",
			entryPrice:   d(50000),
			stopLossRate: d(0.10),
			expectError:  false,
		},
		{
			name:         "Valid 1% trailing stop",
			entryPrice:   d(100),
			stopLossRate: d(0.01),
			expectError:  false,
		},
		{
			name:         "Zero percent stop loss",
			entryPrice:   d(100),
			stopLossRate: d(0),
			expectError:  false,
		},
		{
			name:         "Negative percent should fail",
			entryPrice:   d(100),
			stopLossRate: d(-0.05),
			expectError:  true,
			errorType:    errStopLossRateInvalid,
		},
		{
			name:         "Greater than 1 (100%) should fail",
			entryPrice:   d(100),
			stopLossRate: d(1.1),
			expectError:  true,
			errorType:    errStopLossRateInvalid,
		},
		{
			name:         "Exactly 1 (100%) should work",
			entryPrice:   d(100),
			stopLossRate: d(1.0),
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewTrailing(tt.entryPrice, tt.stopLossRate, trailingCallback)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
				if tt.errorType != nil && err != tt.errorType {
					t.Errorf("Expected error %v but got %v", tt.errorType, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if s == nil {
					t.Error("Expected non-nil stop loss")
				}
			}
		})
	}
}

func TestTrailingStopLoss_CalculateStopLoss(t *testing.T) {
	entryPrice := d(100)
	stopLossRate := d(0.05) // 5%

	s, err := NewTrailing(entryPrice, stopLossRate, nil)
	if err != nil {
		t.Fatalf("Failed to create trailing stop loss: %v", err)
	}

	tests := []struct {
		name             string
		currentPrice     decimal.Decimal
		expectedStopLoss decimal.Decimal
		description      string
	}{
		{
			name:             "Initial stop loss",
			currentPrice:     d(100),
			expectedStopLoss: d(95), // 100 * (1 - 0.05)
			description:      "Stop loss should be 5% below entry",
		},
		{
			name:             "Price increases to 110",
			currentPrice:     d(110),
			expectedStopLoss: d(104.5), // 110 * (1 - 0.05)
			description:      "Stop loss should trail up to 5% below 110",
		},
		{
			name:             "Price increases to 120",
			currentPrice:     d(120),
			expectedStopLoss: d(114), // 120 * (1 - 0.05)
			description:      "Stop loss should trail up to 5% below 120",
		},
		{
			name:             "Price decreases to 115 (stop loss stays)",
			currentPrice:     d(115),
			expectedStopLoss: d(114), // Stays at previous high
			description:      "Stop loss should not decrease",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stopLoss, err := s.CalculateStopLoss(tt.currentPrice)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !stopLoss.Equal(tt.expectedStopLoss) {
				t.Errorf("Expected stop loss %v but got %v - %s",
					tt.expectedStopLoss, stopLoss, tt.description)
			}

			t.Logf("Price: %v, Stop Loss: %v - %s",
				tt.currentPrice, stopLoss, tt.description)
		})
	}
}

func TestTrailingStopLoss_TrailingBehavior(t *testing.T) {
	entryPrice := d(100)
	stopLossRate := d(0.05) // 5%

	s, err := NewTrailing(entryPrice, stopLossRate, nil)
	if err != nil {
		t.Fatalf("Failed to create trailing stop loss: %v", err)
	}

	initialStopLoss, _ := s.GetStopLoss()
	t.Logf("Initial Entry: %v, Initial Stop Loss: %v", entryPrice, initialStopLoss)

	// Price moves up - stop loss should trail
	prices := []struct {
		price            decimal.Decimal
		expectedStopLoss decimal.Decimal
	}{
		{d(105), d(99.75)},  // 105 * 0.95
		{d(110), d(104.5)},  // 110 * 0.95
		{d(115), d(109.25)}, // 115 * 0.95
		{d(108), d(109.25)}, // Price drops but stop loss stays at previous high
		{d(120), d(114)},    // New high, stop loss trails up
	}

	for _, p := range prices {
		stopLoss, _ := s.CalculateStopLoss(p.price)
		if !stopLoss.Equal(p.expectedStopLoss) {
			t.Errorf("At price %v: expected stop loss %v but got %v",
				p.price, p.expectedStopLoss, stopLoss)
		}
		t.Logf("Price: %v, Stop Loss: %v", p.price, stopLoss)
	}
}

func TestTrailingStopLoss_ShouldTriggerStopLoss(t *testing.T) {
	entryPrice := d(100)
	stopLossRate := d(0.05) // 5%

	s, err := NewTrailing(entryPrice, stopLossRate, nil)
	if err != nil {
		t.Fatalf("Failed to create trailing stop loss: %v", err)
	}

	// Move price up to 110
	s.CalculateStopLoss(d(110))
	stopLoss, _ := s.GetStopLoss()
	t.Logf("Price moved to 110, Stop Loss: %v", stopLoss)

	tests := []struct {
		name          string
		currentPrice  decimal.Decimal
		shouldTrigger bool
	}{
		{"Price above stop loss", d(106), false},
		{"Price at stop loss", d(104.5), true},
		{"Price below stop loss", d(104), true},
		{"Price well above stop loss", d(115), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Recreate to avoid triggered state
			s, _ := NewTrailing(entryPrice, stopLossRate, nil)
			s.CalculateStopLoss(d(110))

			triggered, err := s.ShouldTriggerStopLoss(tt.currentPrice)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if triggered != tt.shouldTrigger {
				t.Errorf("Price %v: expected trigger=%v but got %v",
					tt.currentPrice, tt.shouldTrigger, triggered)
			}
		})
	}
}

func TestTrailingStopLoss_ReSet(t *testing.T) {
	entryPrice := d(100)
	stopLossRate := d(0.05)

	s, err := NewTrailing(entryPrice, stopLossRate, nil)
	if err != nil {
		t.Fatalf("Failed to create trailing stop loss: %v", err)
	}

	// Move price up
	s.CalculateStopLoss(d(120))
	stopLossBeforeReset, _ := s.GetStopLoss()
	t.Logf("Stop Loss before reset: %v", stopLossBeforeReset)

	// Reset to new price
	newEntryPrice := d(110)
	err = s.ReSet(newEntryPrice)
	if err != nil {
		t.Errorf("Failed to reset: %v", err)
	}

	// New stop loss: 110 * 0.95 = 104.5
	expectedNewStopLoss := d(104.5)
	newStopLoss, _ := s.GetStopLoss()

	if !newStopLoss.Equal(expectedNewStopLoss) {
		t.Errorf("Expected stop loss %v after reset but got %v",
			expectedNewStopLoss, newStopLoss)
	}

	t.Logf("Reset to %v, New Stop Loss: %v", newEntryPrice, newStopLoss)
}

func TestTrailingStopLoss_Deactivate(t *testing.T) {
	entryPrice := d(100)
	stopLossRate := d(0.05)

	s, err := NewTrailing(entryPrice, stopLossRate, nil)
	if err != nil {
		t.Fatalf("Failed to create trailing stop loss: %v", err)
	}

	// Deactivate stop loss
	err = s.Deactivate()
	if err != nil {
		t.Errorf("Failed to deactivate: %v", err)
	}

	// Should not be able to calculate after deactivation
	_, err = s.CalculateStopLoss(d(110))
	if err != stoploss.ErrStatusInvalid {
		t.Errorf("Expected ErrStatusInvalid after deactivation but got %v", err)
	}

	// Should not be able to check trigger after deactivation
	_, err = s.ShouldTriggerStopLoss(d(90))
	if err != stoploss.ErrStatusInvalid {
		t.Errorf("Expected ErrStatusInvalid after deactivation but got %v", err)
	}
}

// Test with unified Mock data: Uptrend - stop loss should trail up
func TestTrailingStopLoss_WithHistoricalData_UpTrend(t *testing.T) {
	data := GetMockTrendingData()

	entryData := data[0]
	entryPrice := entryData.Close
	stopLossRate := d(0.05) // 5% trailing stop

	triggered := false
	triggerCallback := func(reason string) error {
		triggered = true
		t.Logf("Stop Loss Triggered: %s", reason)
		return nil
	}

	s, err := NewTrailing(entryPrice, stopLossRate, triggerCallback)
	if err != nil {
		t.Fatalf("Failed to create trailing stop loss: %v", err)
	}

	initialStopLoss, _ := s.GetStopLoss()
	t.Logf("Uptrend Test - Entry: %v, Initial Stop Loss: %v (%.2f%%)",
		entryPrice, initialStopLoss, stopLossRate.Mul(d(100)).InexactFloat64())

	highestPrice := entryPrice

	for i := 1; i < len(data); i++ {
		period := data[i]

		// Calculate new stop loss based on current price
		stopLoss, _ := s.CalculateStopLoss(period.Close)

		if period.Close.GreaterThan(highestPrice) {
			highestPrice = period.Close
		}

		pnlPercent := period.Close.Sub(entryPrice).Div(entryPrice).Mul(d(100))

		t.Logf("Period %d: Close=%v, High=%v, StopLoss=%v, PnL=%.2f%%",
			period.Period, period.Close, highestPrice, stopLoss,
			pnlPercent.InexactFloat64())

		// Check if triggered using low price
		isTriggered, _ := s.ShouldTriggerStopLoss(period.Low)
		if isTriggered {
			t.Logf("Period %d: TRIGGERED at %v", period.Period, period.Low)
			break
		}
	}

	if triggered {
		t.Error("Trailing stop should not trigger during uptrend")
	} else {
		finalStopLoss, _ := s.GetStopLoss()
		t.Logf("✓ Stop loss correctly trailed up to %v without triggering", finalStopLoss)
	}
}

// Test with unified Mock data: Downtrend after uptrend
func TestTrailingStopLoss_WithHistoricalData_DownTrend(t *testing.T) {
	data := GetMockHistoricalData()

	entryData := data[0]
	entryPrice := entryData.Close
	stopLossRate := d(0.05) // 5% trailing stop

	triggered := false
	triggerPeriod := 0
	triggerCallback := func(reason string) error {
		triggered = true
		return nil
	}

	s, err := NewTrailing(entryPrice, stopLossRate, triggerCallback)
	if err != nil {
		t.Fatalf("Failed to create trailing stop loss: %v", err)
	}

	initialStopLoss, _ := s.GetStopLoss()
	t.Logf("Downtrend Test - Entry: %v, Initial Stop Loss: %v",
		entryPrice, initialStopLoss)

	for i := 1; i < len(data); i++ {
		period := data[i]

		stopLoss, _ := s.CalculateStopLoss(period.Close)
		pnlPercent := period.Close.Sub(entryPrice).Div(entryPrice).Mul(d(100))

		t.Logf("Period %d: Close=%v, Low=%v, StopLoss=%v, PnL=%.2f%%",
			period.Period, period.Close, period.Low, stopLoss,
			pnlPercent.InexactFloat64())

		isTriggered, _ := s.ShouldTriggerStopLoss(period.Low)
		if isTriggered {
			triggerPeriod = period.Period
			t.Logf("Period %d: TRIGGERED at %v (Stop Loss: %v)",
				period.Period, period.Low, stopLoss)
			break
		}
	}

	if triggered {
		t.Logf("✓ Trailing stop correctly triggered at period %d", triggerPeriod)
	} else {
		t.Log("Stop loss was not triggered")
	}
}

// Test with unified Mock data: Consolidation
func TestTrailingStopLoss_WithHistoricalData_Consolidation(t *testing.T) {
	data := GetMockConsolidationData()

	entryData := data[0]
	entryPrice := entryData.Close
	stopLossRate := d(0.05) // 5% trailing stop

	s, err := NewTrailing(entryPrice, stopLossRate, nil)
	if err != nil {
		t.Fatalf("Failed to create trailing stop loss: %v", err)
	}

	initialStopLoss, _ := s.GetStopLoss()
	t.Logf("Consolidation Test - Entry: %v, Stop Loss: %v",
		entryPrice, initialStopLoss)

	triggerCount := 0

	for i := 1; i < len(data); i++ {
		period := data[i]

		stopLoss, _ := s.CalculateStopLoss(period.Close)
		pnlPercent := period.Close.Sub(entryPrice).Div(entryPrice).Mul(d(100))

		t.Logf("Period %d: Close=%v, StopLoss=%v, PnL=%.2f%%",
			period.Period, period.Close, stopLoss, pnlPercent.InexactFloat64())

		isTriggered, _ := s.ShouldTriggerStopLoss(period.Low)
		if isTriggered {
			triggerCount++
			t.Logf("Period %d: Triggered at %v", period.Period, period.Low)
			break
		}
	}

	if triggerCount > 0 {
		t.Errorf("Trailing stop should not trigger during consolidation with 5%% stop")
	} else {
		t.Log("✓ Trailing stop correctly not triggered during consolidation")
	}
}

// Test with unified Mock data: Volatile market
func TestTrailingStopLoss_WithHistoricalData_Volatile(t *testing.T) {
	data := GetMockVolatileData()

	entryData := data[0]
	entryPrice := entryData.Close
	stopLossRate := d(0.08) // 8% trailing stop for volatile market

	s, err := NewTrailing(entryPrice, stopLossRate, nil)
	if err != nil {
		t.Fatalf("Failed to create trailing stop loss: %v", err)
	}

	initialStopLoss, _ := s.GetStopLoss()
	t.Logf("Volatile Market Test - Entry: %v, Stop Loss: %v (8%%)",
		entryPrice, initialStopLoss)

	for i := 1; i < len(data); i++ {
		period := data[i]

		stopLoss, _ := s.CalculateStopLoss(period.Close)
		pnlPercent := period.Close.Sub(entryPrice).Div(entryPrice).Mul(d(100))

		t.Logf("Period %d: Close=%v, High=%v, Low=%v, StopLoss=%v, PnL=%.2f%%",
			period.Period, period.Close, period.High, period.Low, stopLoss,
			pnlPercent.InexactFloat64())

		isTriggered, _ := s.ShouldTriggerStopLoss(period.Low)
		if isTriggered {
			t.Logf("Period %d: TRIGGERED at %v", period.Period, period.Low)
			break
		}
	}
}

// Compare different trailing percentages
func TestTrailingStopLoss_CompareTrailingPercentages(t *testing.T) {
	data := GetMockHistoricalData()
	entryPrice := data[0].Close

	percentages := []decimal.Decimal{d(0.03), d(0.05), d(0.10)}

	for _, pct := range percentages {
		t.Run(fmt.Sprintf("%.0f%% trailing stop", pct.Mul(d(100)).InexactFloat64()), func(t *testing.T) {
			s, _ := NewTrailing(entryPrice, pct, nil)
			initialStopLoss, _ := s.GetStopLoss()

			t.Logf("Initial Stop Loss: %v (%.2f%%)", initialStopLoss, pct.Mul(d(100)).InexactFloat64())

			triggerPeriod := -1
			for i := 1; i < len(data); i++ {
				period := data[i]
				s.CalculateStopLoss(period.Close)
				isTriggered, _ := s.ShouldTriggerStopLoss(period.Low)
				if isTriggered {
					triggerPeriod = period.Period
					stopLoss, _ := s.GetStopLoss()
					t.Logf("Triggered at period %d, price %v (Stop Loss: %v)",
						period.Period, period.Low, stopLoss)
					break
				}
			}

			if triggerPeriod == -1 {
				t.Log("Not triggered")
			}
		})
	}
}

// Benchmark tests
func BenchmarkNewTrailing(b *testing.B) {
	entryPrice := d(100)
	stopLossRate := d(0.05)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewTrailing(entryPrice, stopLossRate, nil)
	}
}

func BenchmarkTrailingStopLoss_CalculateStopLoss(b *testing.B) {
	entryPrice := d(100)
	stopLossRate := d(0.05)
	s, _ := NewTrailing(entryPrice, stopLossRate, nil)
	currentPrice := d(110)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.CalculateStopLoss(currentPrice)
	}
}

func BenchmarkTrailingStopLoss_ShouldTriggerStopLoss(b *testing.B) {
	entryPrice := d(100)
	stopLossRate := d(0.05)
	currentPrice := d(94)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, _ := NewTrailing(entryPrice, stopLossRate, nil)
		_, _ = s.ShouldTriggerStopLoss(currentPrice)
	}
}

func BenchmarkTrailingStopLoss_ReSet(b *testing.B) {
	entryPrice := d(100)
	stopLossRate := d(0.05)
	s, _ := NewTrailing(entryPrice, stopLossRate, nil)
	newPrice := d(110)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.ReSet(newPrice)
	}
}
