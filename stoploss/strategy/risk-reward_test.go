package strategy

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/stoploss"
)

var riskRewardCallback = func(reason string) error {
	fmt.Println("Risk-Reward Callback:", reason)
	return nil
}

func TestRiskRewardRatio_NewRiskRewardRatio(t *testing.T) {
	tests := []struct {
		name        string
		entryPrice  decimal.Decimal
		riskRatio   decimal.Decimal
		rewardRatio decimal.Decimal
		expectError bool
		errorType   error
	}{
		{
			name:        "Valid 1:2 risk-reward",
			entryPrice:  d(100),
			riskRatio:   d(0.05),
			rewardRatio: d(0.10),
			expectError: false,
		},
		{
			name:        "Valid 1:3 risk-reward",
			entryPrice:  d(50000),
			riskRatio:   d(0.03),
			rewardRatio: d(0.09),
			expectError: false,
		},
		{
			name:        "Valid 1:1 risk-reward",
			entryPrice:  d(100),
			riskRatio:   d(0.05),
			rewardRatio: d(0.05),
			expectError: false,
		},
		{
			name:        "Negative risk ratio should fail",
			entryPrice:  d(100),
			riskRatio:   d(-0.05),
			rewardRatio: d(0.10),
			expectError: true,
			errorType:   errStopLossRateInvalid,
		},
		{
			name:        "Negative reward ratio should fail",
			entryPrice:  d(100),
			riskRatio:   d(0.05),
			rewardRatio: d(-0.10),
			expectError: true,
			errorType:   errStopLossRateInvalid,
		},
		{
			name:        "Risk ratio > 1 should fail",
			entryPrice:  d(100),
			riskRatio:   d(1.1),
			rewardRatio: d(0.10),
			expectError: true,
			errorType:   errStopLossRateInvalid,
		},
		{
			name:        "Reward ratio > 1 should fail",
			entryPrice:  d(100),
			riskRatio:   d(0.05),
			rewardRatio: d(1.1),
			expectError: true,
			errorType:   errStopLossRateInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewRiskRewardRatio(tt.entryPrice, tt.riskRatio, tt.rewardRatio, riskRewardCallback)
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
					t.Error("Expected non-nil risk-reward strategy")
				}
			}
		})
	}
}

func TestRiskRewardRatio_CalculateStopLoss(t *testing.T) {
	entryPrice := d(100)
	riskRatio := d(0.05)   // 5% risk
	rewardRatio := d(0.10) // 10% reward

	s, err := NewRiskRewardRatio(entryPrice, riskRatio, rewardRatio, nil)
	if err != nil {
		t.Fatalf("Failed to create strategy: %v", err)
	}

	// Stop loss = 100 * (1 - 0.05) = 95
	expectedStopLoss := d(95)
	stopLoss, err := s.CalculateStopLoss(entryPrice)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !stopLoss.Equal(expectedStopLoss) {
		t.Errorf("Expected stop loss %v but got %v", expectedStopLoss, stopLoss)
	}

	t.Logf("Entry: %v, Risk: %.0f%%, Stop Loss: %v",
		entryPrice, riskRatio.Mul(d(100)).InexactFloat64(), stopLoss)
}

func TestRiskRewardRatio_CalculateTakeProfit(t *testing.T) {
	entryPrice := d(100)
	riskRatio := d(0.05)   // 5% risk
	rewardRatio := d(0.10) // 10% reward

	s, err := NewRiskRewardRatio(entryPrice, riskRatio, rewardRatio, nil)
	if err != nil {
		t.Fatalf("Failed to create strategy: %v", err)
	}

	// Take profit = 100 * (1 + 0.10) = 110
	expectedTakeProfit := d(110)
	takeProfit, err := s.CalculateTakeProfit(entryPrice)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !takeProfit.Equal(expectedTakeProfit) {
		t.Errorf("Expected take profit %v but got %v", expectedTakeProfit, takeProfit)
	}

	t.Logf("Entry: %v, Reward: %.0f%%, Take Profit: %v",
		entryPrice, rewardRatio.Mul(d(100)).InexactFloat64(), takeProfit)
}

func TestRiskRewardRatio_RiskRewardCalculation(t *testing.T) {
	tests := []struct {
		name               string
		entryPrice         decimal.Decimal
		riskRatio          decimal.Decimal
		rewardRatio        decimal.Decimal
		expectedStopLoss   decimal.Decimal
		expectedTakeProfit decimal.Decimal
		ratioDescription   string
	}{
		{
			name:               "1:2 ratio (5% risk, 10% reward)",
			entryPrice:         d(100),
			riskRatio:          d(0.05),
			rewardRatio:        d(0.10),
			expectedStopLoss:   d(95),
			expectedTakeProfit: d(110),
			ratioDescription:   "1:2",
		},
		{
			name:               "1:3 ratio (3% risk, 9% reward)",
			entryPrice:         d(100),
			riskRatio:          d(0.03),
			rewardRatio:        d(0.09),
			expectedStopLoss:   d(97),
			expectedTakeProfit: d(109),
			ratioDescription:   "1:3",
		},
		{
			name:               "1:1 ratio (5% risk, 5% reward)",
			entryPrice:         d(100),
			riskRatio:          d(0.05),
			rewardRatio:        d(0.05),
			expectedStopLoss:   d(95),
			expectedTakeProfit: d(105),
			ratioDescription:   "1:1",
		},
		{
			name:               "1:4 ratio (2% risk, 8% reward)",
			entryPrice:         d(50000),
			riskRatio:          d(0.02),
			rewardRatio:        d(0.08),
			expectedStopLoss:   d(49000),
			expectedTakeProfit: d(54000),
			ratioDescription:   "1:4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewRiskRewardRatio(tt.entryPrice, tt.riskRatio, tt.rewardRatio, nil)
			if err != nil {
				t.Fatalf("Failed to create strategy: %v", err)
			}

			stopLoss, _ := s.GetStopLoss()
			takeProfit, _ := s.GetTakeProfit()

			if !stopLoss.Equal(tt.expectedStopLoss) {
				t.Errorf("Expected stop loss %v but got %v", tt.expectedStopLoss, stopLoss)
			}

			if !takeProfit.Equal(tt.expectedTakeProfit) {
				t.Errorf("Expected take profit %v but got %v", tt.expectedTakeProfit, takeProfit)
			}

			t.Logf("Ratio %s: Entry=%v, SL=%v (%.0f%%), TP=%v (%.0f%%)",
				tt.ratioDescription, tt.entryPrice, stopLoss,
				tt.riskRatio.Mul(d(100)).InexactFloat64(),
				takeProfit, tt.rewardRatio.Mul(d(100)).InexactFloat64())
		})
	}
}

func TestRiskRewardRatio_ShouldTriggerStopLoss(t *testing.T) {
	entryPrice := d(100)
	riskRatio := d(0.05)
	rewardRatio := d(0.10)

	s, err := NewRiskRewardRatio(entryPrice, riskRatio, rewardRatio, nil)
	if err != nil {
		t.Fatalf("Failed to create strategy: %v", err)
	}

	stopLoss, _ := s.GetStopLoss()
	t.Logf("Stop Loss Level: %v", stopLoss)

	tests := []struct {
		name          string
		currentPrice  decimal.Decimal
		shouldTrigger bool
	}{
		{"Price above stop loss", d(96), false},
		{"Price at stop loss", d(95), true},
		{"Price below stop loss", d(94), true},
		{"Price at entry", d(100), false},
		{"Price near take profit", d(109), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewRiskRewardRatio(entryPrice, riskRatio, rewardRatio, nil)

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

func TestRiskRewardRatio_ShouldTriggerTakeProfit(t *testing.T) {
	entryPrice := d(100)
	riskRatio := d(0.05)
	rewardRatio := d(0.10)

	s, err := NewRiskRewardRatio(entryPrice, riskRatio, rewardRatio, nil)
	if err != nil {
		t.Fatalf("Failed to create strategy: %v", err)
	}

	takeProfit, _ := s.GetTakeProfit()
	t.Logf("Take Profit Level: %v", takeProfit)

	tests := []struct {
		name          string
		currentPrice  decimal.Decimal
		shouldTrigger bool
	}{
		{"Price below take profit", d(109), false},
		{"Price at take profit", d(110), true},
		{"Price above take profit", d(111), true},
		{"Price at entry", d(100), false},
		{"Price near stop loss", d(96), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewRiskRewardRatio(entryPrice, riskRatio, rewardRatio, nil)

			triggered, err := s.ShouldTriggerTakeProfit(tt.currentPrice)
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

func TestRiskRewardRatio_ReSet(t *testing.T) {
	entryPrice := d(100)
	riskRatio := d(0.05)
	rewardRatio := d(0.10)

	s, err := NewRiskRewardRatio(entryPrice, riskRatio, rewardRatio, nil)
	if err != nil {
		t.Fatalf("Failed to create strategy: %v", err)
	}

	initialStopLoss, _ := s.GetStopLoss()
	initialTakeProfit, _ := s.GetTakeProfit()
	t.Logf("Initial - SL: %v, TP: %v", initialStopLoss, initialTakeProfit)

	// Reset to new price
	newEntryPrice := d(110)
	err = s.ReSet(newEntryPrice)
	if err != nil {
		t.Errorf("Failed to reset: %v", err)
	}

	// New stop loss: 110 * (1 - 0.05) = 104.5
	expectedNewStopLoss := d(104.5)
	// New take profit: 110 * (1 + 0.10) = 121
	expectedNewTakeProfit := d(121)

	newStopLoss, _ := s.GetStopLoss()
	newTakeProfit, _ := s.GetTakeProfit()

	if !newStopLoss.Equal(expectedNewStopLoss) {
		t.Errorf("Expected stop loss %v after reset but got %v",
			expectedNewStopLoss, newStopLoss)
	}

	if !newTakeProfit.Equal(expectedNewTakeProfit) {
		t.Errorf("Expected take profit %v after reset but got %v",
			expectedNewTakeProfit, newTakeProfit)
	}

	t.Logf("After Reset - Entry: %v, SL: %v, TP: %v",
		newEntryPrice, newStopLoss, newTakeProfit)
}

func TestRiskRewardRatio_Deactivate(t *testing.T) {
	entryPrice := d(100)
	riskRatio := d(0.05)
	rewardRatio := d(0.10)

	s, err := NewRiskRewardRatio(entryPrice, riskRatio, rewardRatio, nil)
	if err != nil {
		t.Fatalf("Failed to create strategy: %v", err)
	}

	err = s.Deactivate()
	if err != nil {
		t.Errorf("Failed to deactivate: %v", err)
	}

	// After deactivation, should return errors
	_, err = s.CalculateStopLoss(entryPrice)
	if err != stoploss.ErrStatusInvalid {
		t.Errorf("Expected ErrStatusInvalid after deactivation but got %v", err)
	}

	_, err = s.CalculateTakeProfit(entryPrice)
	if err != stoploss.ErrStatusInvalid {
		t.Errorf("Expected ErrStatusInvalid after deactivation but got %v", err)
	}

	_, err = s.ShouldTriggerStopLoss(d(90))
	if err != stoploss.ErrStatusInvalid {
		t.Errorf("Expected ErrStatusInvalid after deactivation but got %v", err)
	}

	_, err = s.ShouldTriggerTakeProfit(d(110))
	if err != stoploss.ErrStatusInvalid {
		t.Errorf("Expected ErrStatusInvalid after deactivation but got %v", err)
	}
}

// Test with historical data: Uptrend reaching take profit
func TestRiskRewardRatio_WithHistoricalData_UpTrend(t *testing.T) {
	data := GetMockTrendingData()

	entryData := data[0]
	entryPrice := entryData.Close
	riskRatio := d(0.05)   // 5% risk
	rewardRatio := d(0.10) // 10% reward (1:2 ratio)

	stopLossTriggered := false
	takeProfitTriggered := false

	triggerCallback := func(reason string) error {
		if reason == "Risk-Reward Stop Loss Triggered" {
			stopLossTriggered = true
		} else if reason == "Risk-Reward Take Profit Triggered" {
			takeProfitTriggered = true
		}
		t.Logf("Triggered: %s", reason)
		return nil
	}

	s, err := NewRiskRewardRatio(entryPrice, riskRatio, rewardRatio, triggerCallback)
	if err != nil {
		t.Fatalf("Failed to create strategy: %v", err)
	}

	stopLoss, _ := s.GetStopLoss()
	takeProfit, _ := s.GetTakeProfit()
	t.Logf("Entry: %v, Stop Loss: %v, Take Profit: %v (1:2 ratio)",
		entryPrice, stopLoss, takeProfit)

	for i := 1; i < len(data); i++ {
		period := data[i]

		pnlPercent := period.Close.Sub(entryPrice).Div(entryPrice).Mul(d(100))

		t.Logf("Period %d: Price=%v, PnL=%.2f%%",
			period.Period, period.Close, pnlPercent.InexactFloat64())

		// Check stop loss
		slTriggered, _ := s.ShouldTriggerStopLoss(period.Low)
		if slTriggered {
			t.Logf("Period %d: STOP LOSS triggered at %v", period.Period, period.Low)
			break
		}

		// Check take profit
		tpTriggered, _ := s.ShouldTriggerTakeProfit(period.High)
		if tpTriggered {
			t.Logf("Period %d: TAKE PROFIT triggered at %v", period.Period, period.High)
			break
		}
	}

	if !takeProfitTriggered && !stopLossTriggered {
		t.Log("Neither take profit nor stop loss triggered during uptrend")
	} else if takeProfitTriggered {
		t.Log("✓ Take profit correctly triggered during uptrend")
	} else if stopLossTriggered {
		t.Error("Stop loss should not trigger during uptrend")
	}
}

// Test with historical data: Downtrend hitting stop loss
func TestRiskRewardRatio_WithHistoricalData_DownTrend(t *testing.T) {
	data := GetMockHistoricalData()

	entryData := data[0]
	entryPrice := entryData.Close
	riskRatio := d(0.05)   // 5% risk
	rewardRatio := d(0.15) // 15% reward (1:3 ratio)

	stopLossTriggered := false
	takeProfitTriggered := false

	triggerCallback := func(reason string) error {
		if reason == "Risk-Reward Stop Loss Triggered" {
			stopLossTriggered = true
		} else if reason == "Risk-Reward Take Profit Triggered" {
			takeProfitTriggered = true
		}
		t.Logf("Triggered: %s", reason)
		return nil
	}

	s, err := NewRiskRewardRatio(entryPrice, riskRatio, rewardRatio, triggerCallback)
	if err != nil {
		t.Fatalf("Failed to create strategy: %v", err)
	}

	stopLoss, _ := s.GetStopLoss()
	takeProfit, _ := s.GetTakeProfit()
	t.Logf("Entry: %v, Stop Loss: %v, Take Profit: %v (1:3 ratio)",
		entryPrice, stopLoss, takeProfit)

	for i := 1; i < len(data); i++ {
		period := data[i]

		pnlPercent := period.Close.Sub(entryPrice).Div(entryPrice).Mul(d(100))

		t.Logf("Period %d: Price=%v, High=%v, Low=%v, PnL=%.2f%%",
			period.Period, period.Close, period.High, period.Low,
			pnlPercent.InexactFloat64())

		// Check stop loss first (priority in downtrend)
		slTriggered, _ := s.ShouldTriggerStopLoss(period.Low)
		if slTriggered {
			t.Logf("Period %d: STOP LOSS triggered at %v", period.Period, period.Low)
			break
		}

		// Check take profit
		tpTriggered, _ := s.ShouldTriggerTakeProfit(period.High)
		if tpTriggered {
			t.Logf("Period %d: TAKE PROFIT triggered at %v", period.Period, period.High)
			break
		}
	}

	if stopLossTriggered {
		t.Log("✓ Stop loss correctly triggered during downtrend")
	} else if takeProfitTriggered {
		t.Log("Take profit triggered (less likely in downtrend)")
	} else {
		t.Log("Neither trigger hit during simulation")
	}
}

// Test with consolidation: Neither trigger should hit
func TestRiskRewardRatio_WithHistoricalData_Consolidation(t *testing.T) {
	data := GetMockConsolidationData()

	entryData := data[0]
	entryPrice := entryData.Close
	riskRatio := d(0.05)  // 5% risk
	rewardRatio := d(0.1) // 10% reward

	s, err := NewRiskRewardRatio(entryPrice, riskRatio, rewardRatio, nil)
	if err != nil {
		t.Fatalf("Failed to create strategy: %v", err)
	}

	stopLoss, _ := s.GetStopLoss()
	takeProfit, _ := s.GetTakeProfit()
	t.Logf("Consolidation Test - Entry: %v, SL: %v, TP: %v",
		entryPrice, stopLoss, takeProfit)

	slTriggered := false
	tpTriggered := false

	for i := 1; i < len(data); i++ {
		period := data[i]

		pnlPercent := period.Close.Sub(entryPrice).Div(entryPrice).Mul(d(100))
		t.Logf("Period %d: Close=%v, PnL=%.2f%%",
			period.Period, period.Close, pnlPercent.InexactFloat64())

		sl, _ := s.ShouldTriggerStopLoss(period.Low)
		tp, _ := s.ShouldTriggerTakeProfit(period.High)

		if sl {
			slTriggered = true
			t.Logf("Period %d: Stop loss triggered", period.Period)
			break
		}
		if tp {
			tpTriggered = true
			t.Logf("Period %d: Take profit triggered", period.Period)
			break
		}
	}

	if !slTriggered && !tpTriggered {
		t.Log("✓ Neither trigger hit during consolidation (expected)")
	} else {
		t.Errorf("Should not trigger during normal consolidation")
	}
}

// Test different risk-reward ratios
func TestRiskRewardRatio_CompareRatios(t *testing.T) {
	data := GetMockHistoricalData()
	entryPrice := data[0].Close

	ratios := []struct {
		name        string
		riskRatio   decimal.Decimal
		rewardRatio decimal.Decimal
	}{
		{"1:1", d(0.05), d(0.05)},
		{"1:2", d(0.05), d(0.10)},
		{"1:3", d(0.05), d(0.15)},
		{"1:4", d(0.05), d(0.20)},
	}

	for _, ratio := range ratios {
		t.Run(ratio.name, func(t *testing.T) {
			s, _ := NewRiskRewardRatio(entryPrice, ratio.riskRatio, ratio.rewardRatio, nil)
			stopLoss, _ := s.GetStopLoss()
			takeProfit, _ := s.GetTakeProfit()

			t.Logf("Ratio %s: SL=%v (%.0f%%), TP=%v (%.0f%%)",
				ratio.name, stopLoss,
				ratio.riskRatio.Mul(d(100)).InexactFloat64(),
				takeProfit,
				ratio.rewardRatio.Mul(d(100)).InexactFloat64())

			slTriggered := false
			tpTriggered := false
			// triggerPeriod := -1

			for i := 1; i < len(data); i++ {
				period := data[i]
				sl, _ := s.ShouldTriggerStopLoss(period.Low)
				tp, _ := s.ShouldTriggerTakeProfit(period.High)

				if sl {
					slTriggered = true
					// triggerPeriod = period.Period
					t.Logf("  SL triggered at period %d", period.Period)
					break
				}
				if tp {
					tpTriggered = true
					// triggerPeriod = period.Period
					t.Logf("  TP triggered at period %d", period.Period)
					break
				}
			}

			if !slTriggered && !tpTriggered {
				t.Log("  No trigger")
			}
		})
	}
}

// Benchmarks
func BenchmarkNewRiskRewardRatio(b *testing.B) {
	entryPrice := d(100)
	riskRatio := d(0.05)
	rewardRatio := d(0.10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewRiskRewardRatio(entryPrice, riskRatio, rewardRatio, nil)
	}
}

func BenchmarkRiskRewardRatio_CalculateStopLoss(b *testing.B) {
	entryPrice := d(100)
	riskRatio := d(0.05)
	rewardRatio := d(0.10)
	s, _ := NewRiskRewardRatio(entryPrice, riskRatio, rewardRatio, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.CalculateStopLoss(entryPrice)
	}
}

func BenchmarkRiskRewardRatio_CalculateTakeProfit(b *testing.B) {
	entryPrice := d(100)
	riskRatio := d(0.05)
	rewardRatio := d(0.10)
	s, _ := NewRiskRewardRatio(entryPrice, riskRatio, rewardRatio, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.CalculateTakeProfit(entryPrice)
	}
}

func BenchmarkRiskRewardRatio_ShouldTriggerStopLoss(b *testing.B) {
	entryPrice := d(100)
	riskRatio := d(0.05)
	rewardRatio := d(0.10)
	currentPrice := d(94)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, _ := NewRiskRewardRatio(entryPrice, riskRatio, rewardRatio, nil)
		_, _ = s.ShouldTriggerStopLoss(currentPrice)
	}
}

func BenchmarkRiskRewardRatio_ShouldTriggerTakeProfit(b *testing.B) {
	entryPrice := d(100)
	riskRatio := d(0.05)
	rewardRatio := d(0.10)
	currentPrice := d(111)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, _ := NewRiskRewardRatio(entryPrice, riskRatio, rewardRatio, nil)
		_, _ = s.ShouldTriggerTakeProfit(currentPrice)
	}
}

func BenchmarkRiskRewardRatio_ReSet(b *testing.B) {
	entryPrice := d(100)
	riskRatio := d(0.05)
	rewardRatio := d(0.10)
	s, _ := NewRiskRewardRatio(entryPrice, riskRatio, rewardRatio, nil)
	newPrice := d(110)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.ReSet(newPrice)
	}
}
