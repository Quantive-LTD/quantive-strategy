package strategy

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/stoploss"
)

var fixedPercentCallback = func(reason string) error {
	fmt.Println("Fixed Percent Callback:", reason)
	return nil
}

func TestFixed(t *testing.T) {
	entryPrice := decimal.NewFromFloat(100.0)
	stopLossPct := decimal.NewFromFloat(0.1) // 10%
	fixedStopLoss, err := NewFixedPercent(entryPrice, stopLossPct, callback)
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
		stopLoss, err := fixedStopLoss.CalculateStopLoss(price)
		if err != nil {
			t.Logf("Error calculating stop loss: %v", err)
			continue
		}
		t.Logf("Current Price: %s, Stop Loss: %s", price.String(), stopLoss.String())
		triggered, err := fixedStopLoss.ShouldTriggerStopLoss(price)
		if err != nil {
			t.Logf("Error checking stop loss trigger: %v", err)
			continue
		}
		if triggered {
			t.Logf("Stop loss triggered at: %s", price.String())
		}
	}
}

func TestFixedPercentStopLoss_NewFixedPercentStopLoss(t *testing.T) {
	tests := []struct {
		name        string
		entryPrice  decimal.Decimal
		stopLossPct decimal.Decimal
		expectError bool
		errorType   error
	}{
		{
			name:        "Valid 5% stop loss",
			entryPrice:  d(100),
			stopLossPct: d(0.05),
			expectError: false,
		},
		{
			name:        "Valid 10% stop loss",
			entryPrice:  d(50000),
			stopLossPct: d(0.10),
			expectError: false,
		},
		{
			name:        "Valid 1% stop loss",
			entryPrice:  d(100),
			stopLossPct: d(0.01),
			expectError: false,
		},
		{
			name:        "Zero percent stop loss",
			entryPrice:  d(100),
			stopLossPct: d(0),
			expectError: false,
		},
		{
			name:        "Negative percent should fail",
			entryPrice:  d(100),
			stopLossPct: d(-0.05),
			expectError: true,
			errorType:   errStopLossRateInvalid,
		},
		{
			name:        "Greater than 1 (100%) should fail",
			entryPrice:  d(100),
			stopLossPct: d(1.1),
			expectError: true,
			errorType:   errStopLossRateInvalid,
		},
		{
			name:        "Exactly 1 (100%) should work",
			entryPrice:  d(100),
			stopLossPct: d(1.0),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewFixedPercent(tt.entryPrice, tt.stopLossPct, fixedPercentCallback)
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

func TestFixedPercentStopLoss_CalculateStopLoss(t *testing.T) {
	tests := []struct {
		name             string
		entryPrice       decimal.Decimal
		stopLossPct      decimal.Decimal
		expectedStopLoss decimal.Decimal
	}{
		{
			name:             "5% stop loss from 100",
			entryPrice:       d(100),
			stopLossPct:      d(0.05),
			expectedStopLoss: d(95), // 100 * (1 - 0.05)
		},
		{
			name:             "10% stop loss from 50000",
			entryPrice:       d(50000),
			stopLossPct:      d(0.10),
			expectedStopLoss: d(45000), // 50000 * (1 - 0.10)
		},
		{
			name:             "2% stop loss from 1000",
			entryPrice:       d(1000),
			stopLossPct:      d(0.02),
			expectedStopLoss: d(980), // 1000 * (1 - 0.02)
		},
		{
			name:             "3.5% stop loss from 200",
			entryPrice:       d(200),
			stopLossPct:      d(0.035),
			expectedStopLoss: d(193), // 200 * (1 - 0.035)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewFixedPercent(tt.entryPrice, tt.stopLossPct, nil)
			if err != nil {
				t.Fatalf("Failed to create stop loss: %v", err)
			}

			stopLoss, err := s.CalculateStopLoss(tt.entryPrice)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !stopLoss.Equal(tt.expectedStopLoss) {
				t.Errorf("Expected stop loss %v but got %v", tt.expectedStopLoss, stopLoss)
			}

			t.Logf("Entry: %v, Stop Loss %%: %.2f%%, Stop Loss Price: %v",
				tt.entryPrice, tt.stopLossPct.Mul(d(100)).InexactFloat64(), stopLoss)
		})
	}
}

func TestFixedPercentStopLoss_ShouldTriggerStopLoss(t *testing.T) {
	entryPrice := d(100)
	stopLossPct := d(0.05) // 5%

	s, err := NewFixedPercent(entryPrice, stopLossPct, nil)
	if err != nil {
		t.Fatalf("Failed to create stop loss: %v", err)
	}

	stopLoss, _ := s.GetStopLoss()
	t.Logf("Entry Price: %v, Stop Loss: %v (5%%)", entryPrice, stopLoss)

	tests := []struct {
		name          string
		currentPrice  decimal.Decimal
		shouldTrigger bool
	}{
		{"Price above stop loss", d(96), false},
		{"Price at stop loss", d(95), true},
		{"Price below stop loss", d(94), true},
		{"Price well above entry", d(105), false},
		{"Price slightly below entry", d(99), false},
		{"Price at entry", d(100), false},
		{"Price far below stop loss", d(90), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重新創建避免已觸發狀態
			s, _ := NewFixedPercent(entryPrice, stopLossPct, nil)

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

func TestFixedPercentStopLoss_GetStopLoss(t *testing.T) {
	entryPrice := d(100)
	stopLossPct := d(0.05)

	s, err := NewFixedPercent(entryPrice, stopLossPct, nil)
	if err != nil {
		t.Fatalf("Failed to create stop loss: %v", err)
	}

	stopLoss, err := s.GetStopLoss()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedStopLoss := d(95)
	if !stopLoss.Equal(expectedStopLoss) {
		t.Errorf("Expected stop loss %v but got %v", expectedStopLoss, stopLoss)
	}
}

func TestFixedPercentStopLoss_ReSet(t *testing.T) {
	entryPrice := d(100)
	stopLossPct := d(0.05)

	s, err := NewFixedPercent(entryPrice, stopLossPct, nil)
	if err != nil {
		t.Fatalf("Failed to create stop loss: %v", err)
	}

	initialStopLoss, _ := s.GetStopLoss()
	t.Logf("Initial Entry: %v, Stop Loss: %v", entryPrice, initialStopLoss)

	// 重置到新價格
	newEntryPrice := d(110)
	err = s.ReSet(newEntryPrice)
	if err != nil {
		t.Errorf("Failed to reset: %v", err)
	}

	// 新止損: 110 * (1 - 0.05) = 104.5
	expectedNewStopLoss := d(104.5)
	newStopLoss, _ := s.GetStopLoss()

	if !newStopLoss.Equal(expectedNewStopLoss) {
		t.Errorf("Expected stop loss %v after reset but got %v",
			expectedNewStopLoss, newStopLoss)
	}

	t.Logf("Reset to new entry %v, New Stop Loss: %v", newEntryPrice, newStopLoss)
}

func TestFixedPercentStopLoss_Deactivate(t *testing.T) {
	entryPrice := d(100)
	stopLossPct := d(0.05)

	s, err := NewFixedPercent(entryPrice, stopLossPct, nil)
	if err != nil {
		t.Fatalf("Failed to create stop loss: %v", err)
	}

	// 停用止損
	err = s.Deactivate()
	if err != nil {
		t.Errorf("Failed to deactivate: %v", err)
	}

	// 停用後應該無法計算止損
	_, err = s.CalculateStopLoss(entryPrice)
	if err != stoploss.ErrStatusInvalid {
		t.Errorf("Expected ErrStatusInvalid after deactivation but got %v", err)
	}

	// 停用後應該無法檢查觸發
	_, err = s.ShouldTriggerStopLoss(d(90))
	if err != stoploss.ErrStatusInvalid {
		t.Errorf("Expected ErrStatusInvalid after deactivation but got %v", err)
	}

	// 停用後應該無法獲取止損
	_, err = s.GetStopLoss()
	if err != stoploss.ErrStatusInvalid {
		t.Errorf("Expected ErrStatusInvalid after deactivation but got %v", err)
	}
}

// 測試使用統一 Mock 數據：下跌趨勢觸發止損
func TestFixedPercentStopLoss_WithHistoricalData_DownTrend(t *testing.T) {
	data := GetMockHistoricalData()

	// 在第一個週期入場
	entryData := data[0]
	entryPrice := entryData.Close
	stopLossPct := d(0.05) // 5% 止損

	triggered := false
	triggerCallback := func(reason string) error {
		triggered = true
		t.Logf("Stop Loss Triggered: %s", reason)
		return nil
	}

	s, err := NewFixedPercent(entryPrice, stopLossPct, triggerCallback)
	if err != nil {
		t.Fatalf("Failed to create stop loss: %v", err)
	}

	stopLoss, _ := s.GetStopLoss()
	t.Logf("Entry Price: %v, Stop Loss Percent: %.2f%%, Stop Loss Price: %v",
		entryPrice, stopLossPct.Mul(d(100)).InexactFloat64(), stopLoss)

	// 模擬每個週期的價格變動
	for i := 1; i < len(data); i++ {
		period := data[i]

		pnlPercent := period.Close.Sub(entryPrice).Div(entryPrice).Mul(d(100))

		t.Logf("Period %d: Price=%v, High=%v, Low=%v, PnL=%.2f%%",
			period.Period, period.Close, period.High, period.Low,
			pnlPercent.InexactFloat64())

		// 檢查是否觸發止損（使用最低價檢查）
		isTriggered, err := s.ShouldTriggerStopLoss(period.Low)
		if err != nil && err != stoploss.ErrStatusInvalid {
			t.Errorf("Period %d: Error checking trigger: %v", period.Period, err)
		}

		if isTriggered {
			t.Logf("Period %d: Stop Loss TRIGGERED at price %v (Stop Loss: %v)",
				period.Period, period.Low, stopLoss)
			break
		}
	}

	if !triggered {
		t.Log("Stop Loss was not triggered during the simulation")
	}
}

// 測試使用統一 Mock 數據：震盪市場不觸發止損
func TestFixedPercentStopLoss_WithHistoricalData_Consolidation(t *testing.T) {
	data := GetMockConsolidationData()

	entryData := data[0]
	entryPrice := entryData.Close
	stopLossPct := d(0.05) // 5% 止損

	s, err := NewFixedPercent(entryPrice, stopLossPct, nil)
	if err != nil {
		t.Fatalf("Failed to create stop loss: %v", err)
	}

	stopLoss, _ := s.GetStopLoss()
	t.Logf("Consolidation Test - Entry: %v, Stop Loss: %v (%.2f%%)",
		entryPrice, stopLoss, stopLossPct.Mul(d(100)).InexactFloat64())

	triggerCount := 0

	for i := 1; i < len(data); i++ {
		period := data[i]

		pnlPercent := period.Close.Sub(entryPrice).Div(entryPrice).Mul(d(100))
		t.Logf("Period %d: Close=%v, PnL=%.2f%%",
			period.Period, period.Close, pnlPercent.InexactFloat64())

		isTriggered, _ := s.ShouldTriggerStopLoss(period.Low)
		if isTriggered {
			triggerCount++
			t.Logf("Period %d: Triggered at %v", period.Period, period.Low)
			break
		}
	}

	if triggerCount > 0 {
		t.Errorf("Stop Loss should not trigger in consolidation with 5%% stop")
	} else {
		t.Log("✓ Stop Loss correctly not triggered during consolidation")
	}
}

// 測試使用統一 Mock 數據：上升趨勢不觸發止損
func TestFixedPercentStopLoss_WithHistoricalData_UpTrend(t *testing.T) {
	data := GetMockTrendingData()

	entryData := data[0]
	entryPrice := entryData.Close
	stopLossPct := d(0.05) // 5% 止損

	s, err := NewFixedPercent(entryPrice, stopLossPct, nil)
	if err != nil {
		t.Fatalf("Failed to create stop loss: %v", err)
	}

	stopLoss, _ := s.GetStopLoss()
	t.Logf("Uptrend Test - Entry: %v, Stop Loss: %v (%.2f%%)",
		entryPrice, stopLoss, stopLossPct.Mul(d(100)).InexactFloat64())

	for i := 1; i < len(data); i++ {
		period := data[i]

		pnlPercent := period.Close.Sub(entryPrice).Div(entryPrice).Mul(d(100))
		t.Logf("Period %d: Close=%v, PnL=%.2f%%",
			period.Period, period.Close, pnlPercent.InexactFloat64())

		isTriggered, _ := s.ShouldTriggerStopLoss(period.Low)
		if isTriggered {
			t.Errorf("Period %d: Stop Loss should not trigger in uptrend", period.Period)
			break
		}
	}

	finalPrice := data[len(data)-1].Close
	profit := finalPrice.Sub(entryPrice)
	profitPercent := profit.Div(entryPrice).Mul(d(100))

	t.Logf("Final Price: %v, Profit: %v (%.2f%%)",
		finalPrice, profit, profitPercent.InexactFloat64())
	t.Log("✓ Stop Loss correctly not triggered during uptrend")
}

// 測試使用統一 Mock 數據：緩慢下跌
func TestFixedPercentStopLoss_WithHistoricalData_GradualDecline(t *testing.T) {
	data := GetMockGradualDeclineData()

	entryData := data[0]
	entryPrice := entryData.Close
	stopLossPct := d(0.05) // 5% 止損

	triggered := false
	triggerPeriod := 0
	triggerPrice := decimal.Zero

	triggerCallback := func(reason string) error {
		triggered = true
		return nil
	}

	s, err := NewFixedPercent(entryPrice, stopLossPct, triggerCallback)
	if err != nil {
		t.Fatalf("Failed to create stop loss: %v", err)
	}

	stopLoss, _ := s.GetStopLoss()
	t.Logf("Gradual Decline Test - Entry: %v, Stop Loss: %v (%.2f%%)",
		entryPrice, stopLoss, stopLossPct.Mul(d(100)).InexactFloat64())

	for i := 1; i < len(data); i++ {
		period := data[i]

		pnlPercent := period.Close.Sub(entryPrice).Div(entryPrice).Mul(d(100))
		t.Logf("Period %d: Close=%v, Low=%v, PnL=%.2f%%",
			period.Period, period.Close, period.Low, pnlPercent.InexactFloat64())

		isTriggered, _ := s.ShouldTriggerStopLoss(period.Low)
		if isTriggered {
			triggerPeriod = period.Period
			triggerPrice = period.Low
			t.Logf("Period %d: TRIGGERED at %v", period.Period, period.Low)
			break
		}
	}

	if triggered {
		expectedTriggerPeriod := 6 // 應該在第6期左右觸發（跌幅達到5%）
		if triggerPeriod < expectedTriggerPeriod-1 || triggerPeriod > expectedTriggerPeriod+1 {
			t.Logf("Warning: Trigger period %d differs from expected ~%d",
				triggerPeriod, expectedTriggerPeriod)
		}
		t.Logf("✓ Stop Loss correctly triggered at period %d, price %v",
			triggerPeriod, triggerPrice)
	} else {
		t.Error("Stop Loss should have been triggered during gradual decline")
	}
}

// 測試使用統一 Mock 數據：急速下跌
func TestFixedPercentStopLoss_WithHistoricalData_SharpDrop(t *testing.T) {
	data := GetMockSharpDropData()

	entryData := data[0]
	entryPrice := entryData.Close
	stopLossPct := d(0.05) // 5% 止損

	triggered := false
	triggerCallback := func(reason string) error {
		triggered = true
		return nil
	}

	s, err := NewFixedPercent(entryPrice, stopLossPct, triggerCallback)
	if err != nil {
		t.Fatalf("Failed to create stop loss: %v", err)
	}

	stopLoss, _ := s.GetStopLoss()
	t.Logf("Sharp Drop Test - Entry: %v, Stop Loss: %v (%.2f%%)",
		entryPrice, stopLoss, stopLossPct.Mul(d(100)).InexactFloat64())

	for i := 1; i < len(data); i++ {
		period := data[i]

		pnlPercent := period.Close.Sub(entryPrice).Div(entryPrice).Mul(d(100))
		t.Logf("Period %d: Close=%v, Low=%v, PnL=%.2f%%",
			period.Period, period.Close, period.Low, pnlPercent.InexactFloat64())

		isTriggered, _ := s.ShouldTriggerStopLoss(period.Low)
		if isTriggered {
			t.Logf("Period %d: TRIGGERED at %v", period.Period, period.Low)
			break
		}
	}

	if !triggered {
		t.Error("Stop Loss should have been triggered during sharp drop")
	} else {
		t.Log("✓ Stop Loss correctly triggered during sharp drop")
	}
}

// 測試比較不同止損百分比
func TestFixedPercentStopLoss_CompareStopLossPercentages(t *testing.T) {
	data := GetMockHistoricalData()
	entryPrice := data[0].Close

	percentages := []decimal.Decimal{d(0.02), d(0.05), d(0.10)}

	for _, pct := range percentages {
		t.Run(fmt.Sprintf("%.0f%% stop loss", pct.Mul(d(100)).InexactFloat64()), func(t *testing.T) {
			s, _ := NewFixedPercent(entryPrice, pct, nil)
			stopLoss, _ := s.GetStopLoss()

			t.Logf("Stop Loss: %v (%.2f%%)", stopLoss, pct.Mul(d(100)).InexactFloat64())

			triggerPeriod := -1
			for i := 1; i < len(data); i++ {
				period := data[i]
				isTriggered, _ := s.ShouldTriggerStopLoss(period.Low)
				if isTriggered {
					triggerPeriod = period.Period
					t.Logf("Triggered at period %d, price %v", period.Period, period.Low)
					break
				}
			}

			if triggerPeriod == -1 {
				t.Log("Not triggered")
			}
		})
	}
}

// 基準測試：創建固定百分比止損
func BenchmarkNewFixedPercentStopLoss(b *testing.B) {
	entryPrice := d(100)
	stopLossPct := d(0.05)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewFixedPercent(entryPrice, stopLossPct, nil)
	}
}

// 基準測試：計算止損
func BenchmarkFixedPercentStopLoss_CalculateStopLoss(b *testing.B) {
	entryPrice := d(100)
	stopLossPct := d(0.05)
	s, _ := NewFixedPercent(entryPrice, stopLossPct, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.CalculateStopLoss(entryPrice)
	}
}

// 基準測試：檢查觸發
func BenchmarkFixedPercentStopLoss_ShouldTriggerStopLoss(b *testing.B) {
	entryPrice := d(100)
	stopLossPct := d(0.05)
	currentPrice := d(94)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, _ := NewFixedPercent(entryPrice, stopLossPct, nil)
		_, _ = s.ShouldTriggerStopLoss(currentPrice)
	}
}

// 基準測試：重置止損
func BenchmarkFixedPercentStopLoss_ReSet(b *testing.B) {
	entryPrice := d(100)
	stopLossPct := d(0.05)
	s, _ := NewFixedPercent(entryPrice, stopLossPct, nil)
	newPrice := d(110)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.ReSet(newPrice)
	}
}
