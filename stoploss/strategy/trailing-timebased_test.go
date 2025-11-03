package strategy

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/stoploss"
)

var timeBasedCallback = func(reason string) error {
	fmt.Println("Time-Based Trailing Callback:", reason)
	return nil
}

func TestTrailingTimeWithThreshold(t *testing.T) {
	entryPrice := decimal.NewFromFloat(100)
	stopLossRate := decimal.NewFromFloat(0.05)
	timeThreshold := int64(300) // 5 minutes
	trailingTimeBased, err := NewTrailingTimeBased(entryPrice, stopLossRate, timeThreshold, callback)
	if err != nil {
		t.Fatalf("Failed to create trailing time-based stop loss: %v", err)
	}
	if trailingTimeBased == nil {
		t.Fatal("Trailing time-based stop loss is nil")
	}
	_, err = trailingTimeBased.GetStopLoss()
	if err != nil {
		t.Fatalf("Failed to get stop loss: %v", err)
	}
	_, err = trailingTimeBased.ShouldTriggerStopLoss(decimal.NewFromFloat(95), 1000)
	if err != nil {
		t.Fatalf("Failed to check if stop loss should trigger: %v", err)
	}
	_, err = trailingTimeBased.CalculateStopLoss(decimal.NewFromFloat(110))
	if err != nil {
		t.Fatalf("Failed to calculate stop loss: %v", err)
	}
	err = trailingTimeBased.ReSet(decimal.NewFromFloat(105))
	if err != nil {
		t.Fatalf("Failed to reset trailing time-based stop loss: %v", err)
	}
}

func TestTrailingTimeBased_NewTrailingTimeBased(t *testing.T) {
	tests := []struct {
		name          string
		entryPrice    decimal.Decimal
		stopLossRate  decimal.Decimal
		timeThreshold int64
		expectError   bool
		errorType     error
	}{
		{
			name:          "Valid 5% trailing with 5 min threshold",
			entryPrice:    d(100),
			stopLossRate:  d(0.05),
			timeThreshold: 300,
			expectError:   false,
		},
		{
			name:          "Valid 10% trailing with 10 min threshold",
			entryPrice:    d(50000),
			stopLossRate:  d(0.10),
			timeThreshold: 600,
			expectError:   false,
		},
		{
			name:          "Valid with 1 second threshold",
			entryPrice:    d(100),
			stopLossRate:  d(0.05),
			timeThreshold: 1,
			expectError:   false,
		},
		{
			name:          "Zero time threshold",
			entryPrice:    d(100),
			stopLossRate:  d(0.05),
			timeThreshold: 0,
			expectError:   true,
			errorType:     errTimeThresholdInvalid,
		},
		{
			name:          "Negative percent should fail",
			entryPrice:    d(100),
			stopLossRate:  d(-0.05),
			timeThreshold: 300,
			expectError:   true,
			errorType:     errStopLossRateInvalid,
		},
		{
			name:          "Greater than 100% should fail",
			entryPrice:    d(100),
			stopLossRate:  d(1.1),
			timeThreshold: 300,
			expectError:   true,
			errorType:     errStopLossRateInvalid,
		},
		{
			name:          "Negative time threshold should fail",
			entryPrice:    d(100),
			stopLossRate:  d(0.05),
			timeThreshold: -100,
			expectError:   true,
			errorType:     errTimeThresholdInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewTrailingTimeBased(tt.entryPrice, tt.stopLossRate, tt.timeThreshold, timeBasedCallback)
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

func TestTrailingTimeBased_CalculateStopLoss(t *testing.T) {
	entryPrice := d(100)
	stopLossRate := d(0.05)     // 5%
	timeThreshold := int64(300) // 5 minutes

	s, err := NewTrailingTimeBased(entryPrice, stopLossRate, timeThreshold, nil)
	if err != nil {
		t.Fatalf("Failed to create trailing time-based stop loss: %v", err)
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

func TestTrailingTimeBased_TimeThresholdBehavior(t *testing.T) {
	entryPrice := d(100)
	stopLossRate := d(0.05)     // 5%
	timeThreshold := int64(300) // 5 minutes

	s, err := NewTrailingTimeBased(entryPrice, stopLossRate, timeThreshold, nil)
	if err != nil {
		t.Fatalf("Failed to create trailing time-based stop loss: %v", err)
	}

	stopLoss, _ := s.GetStopLoss()
	t.Logf("Entry: %v, Stop Loss: %v, Time Threshold: %d seconds",
		entryPrice, stopLoss, timeThreshold)

	// Simulate time progression
	baseTime := int64(1000000) // Starting timestamp
	pricesBelowStop := d(94)   // Below stop loss of 95

	tests := []struct {
		name          string
		timestamp     int64
		price         decimal.Decimal
		shouldTrigger bool
		description   string
	}{
		{
			name:          "Initial - below stop but not enough time",
			timestamp:     baseTime,
			price:         pricesBelowStop,
			shouldTrigger: false,
			description:   "0 seconds elapsed, needs 300s",
		},
		{
			name:          "100 seconds - still not enough",
			timestamp:     baseTime + 100,
			price:         pricesBelowStop,
			shouldTrigger: false,
			description:   "100 seconds elapsed, needs 300s",
		},
		{
			name:          "299 seconds - just before threshold",
			timestamp:     baseTime + 299,
			price:         pricesBelowStop,
			shouldTrigger: false,
			description:   "299 seconds elapsed, needs 300s",
		},
		{
			name:          "300 seconds - threshold reached",
			timestamp:     baseTime + 300,
			price:         pricesBelowStop,
			shouldTrigger: true,
			description:   "300 seconds elapsed, threshold met",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Recreate for each test to avoid triggered state
			s, _ := NewTrailingTimeBased(entryPrice, stopLossRate, timeThreshold, nil)

			// Simulate continuous low prices until test timestamp
			for ts := baseTime; ts <= tt.timestamp; ts += 10 {
				triggered, _ := s.ShouldTriggerStopLoss(pricesBelowStop, ts)
				if ts == tt.timestamp {
					if triggered != tt.shouldTrigger {
						t.Errorf("At timestamp %d (%s): expected trigger=%v but got %v",
							tt.timestamp, tt.description, tt.shouldTrigger, triggered)
					} else {
						t.Logf("✓ %s: trigger=%v (correct)", tt.description, triggered)
					}
				}
			}
		})
	}
}

func TestTrailingTimeBased_PriceRecoveryResetsTimer(t *testing.T) {
	entryPrice := d(100)
	stopLossRate := d(0.05)     // 5% - stop loss at 95
	timeThreshold := int64(300) // 5 minutes

	s, err := NewTrailingTimeBased(entryPrice, stopLossRate, timeThreshold, nil)
	if err != nil {
		t.Fatalf("Failed to create trailing time-based stop loss: %v", err)
	}

	baseTime := int64(1000000)

	// Price drops below stop loss
	triggered, _ := s.ShouldTriggerStopLoss(d(94), baseTime)
	if triggered {
		t.Error("Should not trigger immediately")
	}
	t.Logf("T+0s: Price 94 (below 95 stop)")

	// 200 seconds pass
	triggered, _ = s.ShouldTriggerStopLoss(d(94), baseTime+200)
	if triggered {
		t.Error("Should not trigger at 200s (needs 300s)")
	}
	t.Logf("T+200s: Price 94 (still below stop)")

	// Price recovers above stop loss - should reset timer
	triggered, _ = s.ShouldTriggerStopLoss(d(96), baseTime+250)
	if triggered {
		t.Error("Should not trigger when price recovers")
	}
	t.Logf("T+250s: Price 96 (recovered above 95 stop) - timer should reset")

	// Price drops again below stop loss
	triggered, _ = s.ShouldTriggerStopLoss(d(94), baseTime+260)
	if triggered {
		t.Error("Should not trigger immediately after recovery")
	}
	t.Logf("T+260s: Price 94 (below stop again) - timer restarted")

	// 100 seconds after the new drop (not enough time)
	triggered, _ = s.ShouldTriggerStopLoss(d(94), baseTime+360)
	if triggered {
		t.Error("Should not trigger - only 100s since last drop below stop")
	}
	t.Logf("T+360s: Price 94 (only 100s since timer restart)")

	// 300 seconds after the recovery drop - should trigger now
	triggered, _ = s.ShouldTriggerStopLoss(d(94), baseTime+560)
	if !triggered {
		t.Error("Should trigger - 300s since timer restart at T+260")
	}
	t.Logf("T+560s: Price 94 - TRIGGERED (300s elapsed since T+260)")
}

func TestTrailingTimeBased_ShouldTriggerStopLoss(t *testing.T) {
	entryPrice := d(100)
	stopLossRate := d(0.05)    // 5%
	timeThreshold := int64(60) // 1 minute

	s, err := NewTrailingTimeBased(entryPrice, stopLossRate, timeThreshold, nil)
	if err != nil {
		t.Fatalf("Failed to create trailing time-based stop loss: %v", err)
	}

	// Move price up to 110
	s.CalculateStopLoss(d(110))
	stopLoss, _ := s.GetStopLoss()
	t.Logf("Price moved to 110, Stop Loss: %v", stopLoss)

	baseTime := int64(1000000)

	tests := []struct {
		name          string
		price         decimal.Decimal
		timestamp     int64
		shouldTrigger bool
		description   string
	}{
		{"Price above stop at T+0", d(106), baseTime, false, "Above stop loss"},
		{"Price at stop at T+0", d(104.5), baseTime, false, "At stop but no time elapsed"},
		{"Price below stop at T+0", d(104), baseTime, false, "Below stop but no time elapsed"},
		{"Price below stop at T+30s", d(104), baseTime + 30, false, "Below stop for 30s (needs 60s)"},
		{"Price below stop at T+60s", d(104), baseTime + 60, true, "Below stop for 60s - should trigger"},
		{"Price well above stop", d(115), baseTime + 100, false, "Well above stop loss"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Recreate to avoid triggered state
			s, _ := NewTrailingTimeBased(entryPrice, stopLossRate, timeThreshold, nil)
			s.CalculateStopLoss(d(110))

			// Simulate continuous prices
			if tt.price.LessThanOrEqual(d(104.5)) {
				for ts := baseTime; ts <= tt.timestamp; ts += 10 {
					triggered, _ := s.ShouldTriggerStopLoss(tt.price, ts)
					if ts == tt.timestamp && triggered != tt.shouldTrigger {
						t.Errorf("%s: expected trigger=%v but got %v",
							tt.description, tt.shouldTrigger, triggered)
					}
				}
			} else {
				triggered, _ := s.ShouldTriggerStopLoss(tt.price, tt.timestamp)
				if triggered != tt.shouldTrigger {
					t.Errorf("%s: expected trigger=%v but got %v",
						tt.description, tt.shouldTrigger, triggered)
				}
			}
		})
	}
}

func TestTrailingTimeBased_ReSet(t *testing.T) {
	entryPrice := d(100)
	stopLossRate := d(0.05)
	timeThreshold := int64(300)

	s, err := NewTrailingTimeBased(entryPrice, stopLossRate, timeThreshold, nil)
	if err != nil {
		t.Fatalf("Failed to create trailing time-based stop loss: %v", err)
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

	expectedNewStopLoss := d(104.5) // 110 * 0.95
	newStopLoss, _ := s.GetStopLoss()

	if !newStopLoss.Equal(expectedNewStopLoss) {
		t.Errorf("Expected stop loss %v after reset but got %v",
			expectedNewStopLoss, newStopLoss)
	}

	t.Logf("Reset to %v, New Stop Loss: %v", newEntryPrice, newStopLoss)
}

func TestTrailingTimeBased_Deactivate(t *testing.T) {
	entryPrice := d(100)
	stopLossRate := d(0.05)
	timeThreshold := int64(300)

	s, err := NewTrailingTimeBased(entryPrice, stopLossRate, timeThreshold, nil)
	if err != nil {
		t.Fatalf("Failed to create trailing time-based stop loss: %v", err)
	}

	err = s.Deactivate()
	if err != nil {
		t.Errorf("Failed to deactivate: %v", err)
	}

	_, err = s.CalculateStopLoss(d(110))
	if err != stoploss.ErrStatusInvalid {
		t.Errorf("Expected ErrStatusInvalid after deactivation but got %v", err)
	}

	_, err = s.ShouldTriggerStopLoss(d(90), 1000)
	if err != stoploss.ErrStatusInvalid {
		t.Errorf("Expected ErrStatusInvalid after deactivation but got %v", err)
	}
}

// Test with unified Mock data: Uptrend with time simulation
func TestTrailingTimeBased_WithHistoricalData_UpTrend(t *testing.T) {
	data := GetMockTrendingData()

	entryData := data[0]
	entryPrice := entryData.Close
	stopLossRate := d(0.05)     // 5% trailing stop
	timeThreshold := int64(600) // 10 minutes

	triggered := false
	triggerCallback := func(reason string) error {
		triggered = true
		t.Logf("Stop Loss Triggered: %s", reason)
		return nil
	}

	s, err := NewTrailingTimeBased(entryPrice, stopLossRate, timeThreshold, triggerCallback)
	if err != nil {
		t.Fatalf("Failed to create trailing time-based stop loss: %v", err)
	}

	initialStopLoss, _ := s.GetStopLoss()
	t.Logf("Uptrend Test - Entry: %v, Initial Stop Loss: %v, Time Threshold: %ds",
		entryPrice, initialStopLoss, timeThreshold)

	baseTime := int64(1000000)
	highestPrice := entryPrice

	for i := 1; i < len(data); i++ {
		period := data[i]
		timestamp := baseTime + int64(i*900) // 15 minutes per period

		stopLoss, _ := s.CalculateStopLoss(period.Close)

		if period.Close.GreaterThan(highestPrice) {
			highestPrice = period.Close
		}

		pnlPercent := period.Close.Sub(entryPrice).Div(entryPrice).Mul(d(100))

		t.Logf("Period %d (T+%dm): Close=%v, High=%v, StopLoss=%v, PnL=%.2f%%",
			period.Period, (timestamp-baseTime)/60, period.Close, highestPrice, stopLoss,
			pnlPercent.InexactFloat64())

		isTriggered, _ := s.ShouldTriggerStopLoss(period.Low, timestamp)
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

// Test with unified Mock data: Gradual decline with time
func TestTrailingTimeBased_WithHistoricalData_GradualDecline(t *testing.T) {
	data := GetMockGradualDeclineData()

	entryData := data[0]
	entryPrice := entryData.Close
	stopLossRate := d(0.03)     // 3% trailing stop
	timeThreshold := int64(300) // 5 minutes

	triggered := false
	triggerPeriod := 0
	triggerCallback := func(reason string) error {
		triggered = true
		return nil
	}

	s, err := NewTrailingTimeBased(entryPrice, stopLossRate, timeThreshold, triggerCallback)
	if err != nil {
		t.Fatalf("Failed to create trailing time-based stop loss: %v", err)
	}

	initialStopLoss, _ := s.GetStopLoss()
	t.Logf("Gradual Decline Test - Entry: %v, Stop Loss: %v (3%%), Threshold: 5min",
		entryPrice, initialStopLoss)

	baseTime := int64(1000000)

	for i := 1; i < len(data); i++ {
		period := data[i]
		timestamp := baseTime + int64(i*300) // 5 minutes per period

		stopLoss, _ := s.CalculateStopLoss(period.Close)
		pnlPercent := period.Close.Sub(entryPrice).Div(entryPrice).Mul(d(100))

		t.Logf("Period %d (T+%dm): Close=%v, Low=%v, StopLoss=%v, PnL=%.2f%%",
			period.Period, (timestamp-baseTime)/60, period.Close, period.Low, stopLoss,
			pnlPercent.InexactFloat64())

		isTriggered, _ := s.ShouldTriggerStopLoss(period.Low, timestamp)
		if isTriggered {
			triggerPeriod = period.Period
			t.Logf("Period %d (T+%dm): TRIGGERED at %v (Stop Loss: %v)",
				period.Period, (timestamp-baseTime)/60, period.Low, stopLoss)
			break
		}
	}

	if triggered {
		t.Logf("✓ Time-based trailing stop correctly triggered at period %d", triggerPeriod)
	} else {
		t.Log("Stop loss was not triggered")
	}
}

// Test with unified Mock data: Sharp drop with time confirmation
func TestTrailingTimeBased_WithHistoricalData_SharpDrop(t *testing.T) {
	data := GetMockSharpDropData()

	entryData := data[0]
	entryPrice := entryData.Close
	stopLossRate := d(0.05)     // 5% trailing stop
	timeThreshold := int64(180) // 3 minutes

	s, err := NewTrailingTimeBased(entryPrice, stopLossRate, timeThreshold, nil)
	if err != nil {
		t.Fatalf("Failed to create trailing time-based stop loss: %v", err)
	}

	initialStopLoss, _ := s.GetStopLoss()
	t.Logf("Sharp Drop Test - Entry: %v, Stop Loss: %v, Threshold: 3min",
		entryPrice, initialStopLoss)

	baseTime := int64(1000000)
	triggerPeriod := -1

	for i := 1; i < len(data); i++ {
		period := data[i]
		timestamp := baseTime + int64(i*180) // 3 minutes per period

		stopLoss, _ := s.CalculateStopLoss(period.Close)
		pnlPercent := period.Close.Sub(entryPrice).Div(entryPrice).Mul(d(100))

		t.Logf("Period %d (T+%dm): Close=%v, Low=%v, StopLoss=%v, PnL=%.2f%%",
			period.Period, (timestamp-baseTime)/60, period.Close, period.Low, stopLoss,
			pnlPercent.InexactFloat64())

		isTriggered, _ := s.ShouldTriggerStopLoss(period.Low, timestamp)
		if isTriggered {
			triggerPeriod = period.Period
			t.Logf("Period %d: TRIGGERED after sharp drop", period.Period)
			break
		}
	}

	if triggerPeriod > 0 {
		t.Logf("✓ Time-based stop correctly triggered at period %d", triggerPeriod)
	}
}

// Test with unified Mock data: Recovery scenario
func TestTrailingTimeBased_WithHistoricalData_Recovery(t *testing.T) {
	data := GetMockRecoveryData()

	entryData := data[0]
	entryPrice := entryData.Close
	stopLossRate := d(0.05)     // 5% trailing stop
	timeThreshold := int64(600) // 10 minutes

	s, err := NewTrailingTimeBased(entryPrice, stopLossRate, timeThreshold, nil)
	if err != nil {
		t.Fatalf("Failed to create trailing time-based stop loss: %v", err)
	}

	initialStopLoss, _ := s.GetStopLoss()
	t.Logf("Recovery Test - Entry: %v, Stop Loss: %v, Threshold: 10min",
		entryPrice, initialStopLoss)

	baseTime := int64(1000000)

	for i := 1; i < len(data); i++ {
		period := data[i]
		timestamp := baseTime + int64(i*300) // 5 minutes per period

		stopLoss, _ := s.CalculateStopLoss(period.Close)
		pnlPercent := period.Close.Sub(entryPrice).Div(entryPrice).Mul(d(100))

		t.Logf("Period %d (T+%dm): Close=%v, Low=%v, StopLoss=%v, PnL=%.2f%%",
			period.Period, (timestamp-baseTime)/60, period.Close, period.Low, stopLoss,
			pnlPercent.InexactFloat64())

		isTriggered, _ := s.ShouldTriggerStopLoss(period.Low, timestamp)
		if isTriggered {
			t.Logf("Period %d: TRIGGERED at %v", period.Period, period.Low)
			break
		}
	}

	t.Log("✓ Recovery scenario completed - stop loss trails up during recovery")
}

// Compare different time thresholds
func TestTrailingTimeBased_CompareTimeThresholds(t *testing.T) {
	data := GetMockHistoricalData()
	entryPrice := data[0].Close
	stopLossRate := d(0.05)

	thresholds := []int64{60, 300, 600} // 1min, 5min, 10min

	for _, threshold := range thresholds {
		t.Run(fmt.Sprintf("%ds threshold", threshold), func(t *testing.T) {
			s, _ := NewTrailingTimeBased(entryPrice, stopLossRate, threshold, nil)
			initialStopLoss, _ := s.GetStopLoss()

			t.Logf("Initial Stop Loss: %v, Threshold: %ds", initialStopLoss, threshold)

			baseTime := int64(1000000)
			triggerPeriod := -1

			for i := 1; i < len(data); i++ {
				period := data[i]
				timestamp := baseTime + int64(i*300) // 5 minutes per period

				s.CalculateStopLoss(period.Close)
				isTriggered, _ := s.ShouldTriggerStopLoss(period.Low, timestamp)
				if isTriggered {
					triggerPeriod = period.Period
					stopLoss, _ := s.GetStopLoss()
					t.Logf("Triggered at period %d (T+%dm), price %v (Stop Loss: %v)",
						period.Period, (timestamp-baseTime)/60, period.Low, stopLoss)
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
func BenchmarkNewTrailingTimeBased(b *testing.B) {
	entryPrice := d(100)
	stopLossRate := d(0.05)
	timeThreshold := int64(300)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewTrailingTimeBased(entryPrice, stopLossRate, timeThreshold, nil)
	}
}

func BenchmarkTrailingTimeBased_CalculateStopLoss(b *testing.B) {
	entryPrice := d(100)
	stopLossRate := d(0.05)
	timeThreshold := int64(300)
	s, _ := NewTrailingTimeBased(entryPrice, stopLossRate, timeThreshold, nil)
	currentPrice := d(110)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.CalculateStopLoss(currentPrice)
	}
}

func BenchmarkTrailingTimeBased_ShouldTriggerStopLoss(b *testing.B) {
	entryPrice := d(100)
	stopLossRate := d(0.05)
	timeThreshold := int64(300)
	currentPrice := d(94)
	timestamp := int64(1000000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, _ := NewTrailingTimeBased(entryPrice, stopLossRate, timeThreshold, nil)
		_, _ = s.ShouldTriggerStopLoss(currentPrice, timestamp)
	}
}

func BenchmarkTrailingTimeBased_ReSet(b *testing.B) {
	entryPrice := d(100)
	stopLossRate := d(0.05)
	timeThreshold := int64(300)
	s, _ := NewTrailingTimeBased(entryPrice, stopLossRate, timeThreshold, nil)
	newPrice := d(110)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.ReSet(newPrice)
	}
}
