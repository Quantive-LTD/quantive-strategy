package strategy

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/stoploss"
)

var callback = func(reason string) error {
	fmt.Println("Callback:", reason)
	return nil
}

func TestATRStopLoss(t *testing.T) {
	entry := decimal.NewFromFloat(100)
	atr := decimal.NewFromFloat(2.5)
	mult := decimal.NewFromFloat(3)
	s, err := NewATRStop(entry, atr, mult, callback)
	if err != nil {
		t.Fatal("Failed to create ATR Stop Loss:", err)
	}
	stop, err := s.CalculateStopLoss(decimal.NewFromFloat(100))
	if err != nil {
		t.Fatal("Failed to calculate stop loss:", err)
	}
	fmt.Println("Initial Stop Loss:", stop)
	err = s.UpdateATR(decimal.NewFromFloat(3.0))
	if err != nil {
		t.Fatal("Failed to update ATR:", err)
	}
	_, err = s.ShouldTriggerStopLoss(decimal.NewFromFloat(91))
	if err != nil {
		t.Fatal("Failed to check if stop loss triggered:", err)
	}
	err = s.ReSet(decimal.NewFromFloat(105))
	if err != nil {
		t.Fatal("Failed to reset ATR Stop Loss:", err)
	}
}

func TestATRStopLoss_NewATRStop(t *testing.T) {
	tests := []struct {
		name        string
		entryPrice  decimal.Decimal
		atr         decimal.Decimal
		multiplier  decimal.Decimal
		expectError bool
		errorType   error
	}{
		{
			name:        "Valid ATR Stop Loss",
			entryPrice:  d(100),
			atr:         d(2.5),
			multiplier:  d(2),
			expectError: false,
		},
		{
			name:        "Zero multiplier should fail",
			entryPrice:  d(100),
			atr:         d(2.5),
			multiplier:  d(0),
			expectError: true,
			errorType:   errATRStopLossKInvalid,
		},
		{
			name:        "Negative multiplier should fail",
			entryPrice:  d(100),
			atr:         d(2.5),
			multiplier:  d(-1),
			expectError: true,
			errorType:   errATRStopLossKInvalid,
		},
		{
			name:        "High multiplier",
			entryPrice:  d(50000),
			atr:         d(500),
			multiplier:  d(3),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewATRStop(tt.entryPrice, tt.atr, tt.multiplier, callback)
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

func TestATRStopLoss_CalculateStopLoss(t *testing.T) {
	entryPrice := d(100)
	atr := d(2.5)
	multiplier := d(2)

	s, err := NewATRStop(entryPrice, atr, multiplier, callback)
	if err != nil {
		t.Fatalf("Failed to create ATR Stop Loss: %v", err)
	}

	expectedStopLoss := d(95)

	stopLoss, err := s.CalculateStopLoss(entryPrice)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !stopLoss.Equal(expectedStopLoss) {
		t.Errorf("Expected stop loss %v but got %v", expectedStopLoss, stopLoss)
	}

	t.Logf("Entry: %v, ATR: %v, Multiplier: %v, Stop Loss: %v",
		entryPrice, atr, multiplier, stopLoss)
}

func TestATRStopLoss_UpdateATR(t *testing.T) {
	entryPrice := d(100)
	initialATR := d(2.5)
	multiplier := d(2)

	s, err := NewATRStop(entryPrice, initialATR, multiplier, callback)
	if err != nil {
		t.Fatalf("Failed to create ATR Stop Loss: %v", err)
	}

	// 初始止損: 100 - (2.5 * 2) = 95
	initialStopLoss, _ := s.GetStopLoss()
	t.Logf("Initial Stop Loss: %v", initialStopLoss)

	// 更新 ATR 到 3.0
	newATR := d(3.0)
	err = s.UpdateATR(newATR)
	if err != nil {
		t.Errorf("Failed to update ATR: %v", err)
	}

	// 新止損: 100 - (3.0 * 2) = 94
	expectedNewStopLoss := d(94)
	newStopLoss, _ := s.GetStopLoss()

	if !newStopLoss.Equal(expectedNewStopLoss) {
		t.Errorf("Expected stop loss %v after ATR update but got %v",
			expectedNewStopLoss, newStopLoss)
	}

	t.Logf("Updated ATR: %v, New Stop Loss: %v", newATR, newStopLoss)
}

func TestATRStopLoss_ShouldTriggerStopLoss(t *testing.T) {
	entryPrice := d(100)
	atr := d(2.5)
	multiplier := d(2)

	s, err := NewATRStop(entryPrice, atr, multiplier, callback)
	if err != nil {
		t.Fatalf("Failed to create ATR Stop Loss: %v", err)
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
		{"Price well above entry", d(105), false},
		{"Price slightly below entry", d(99), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重新創建以避免已觸發狀態
			s, _ := NewATRStop(entryPrice, atr, multiplier, nil)

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

func TestATRStopLoss_ReSet(t *testing.T) {
	entryPrice := d(100)
	atr := d(2.5)
	multiplier := d(2)

	s, err := NewATRStop(entryPrice, atr, multiplier, callback)
	if err != nil {
		t.Fatalf("Failed to create ATR Stop Loss: %v", err)
	}

	initialStopLoss, _ := s.GetStopLoss()
	t.Logf("Initial Stop Loss: %v", initialStopLoss)

	// 重置到新價格
	newEntryPrice := d(110)
	err = s.ReSet(newEntryPrice)
	if err != nil {
		t.Errorf("Failed to reset: %v", err)
	}

	// 新止損: 110 - (2.5 * 2) = 105
	expectedNewStopLoss := d(105)
	newStopLoss, _ := s.GetStopLoss()

	if !newStopLoss.Equal(expectedNewStopLoss) {
		t.Errorf("Expected stop loss %v after reset but got %v",
			expectedNewStopLoss, newStopLoss)
	}

	t.Logf("Reset to new entry %v, New Stop Loss: %v", newEntryPrice, newStopLoss)
}

func TestATRStopLoss_Deactivate(t *testing.T) {
	entryPrice := d(100)
	atr := d(2.5)
	multiplier := d(2)

	s, err := NewATRStop(entryPrice, atr, multiplier, callback)
	if err != nil {
		t.Fatalf("Failed to create ATR Stop Loss: %v", err)
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
}

// 測試使用 Mock 歷史數據：下跌趨勢觸發止損
func TestATRStopLoss_WithHistoricalData_DownTrend(t *testing.T) {
	data := GetMockHistoricalData()

	// 在第一個週期入場
	entryData := data[0]
	entryPrice := entryData.Close
	multiplier := d(2)

	triggered := false
	triggerCallback := func(reason string) error {
		triggered = true
		t.Logf("Stop Loss Triggered: %s", reason)
		return nil
	}

	s, err := NewATRStop(entryPrice, entryData.ATR, multiplier, triggerCallback)
	if err != nil {
		t.Fatalf("Failed to create ATR Stop Loss: %v", err)
	}

	t.Logf("Entry Price: %v, Initial ATR: %v, Multiplier: %v",
		entryPrice, entryData.ATR, multiplier)

	initialStopLoss, _ := s.GetStopLoss()
	t.Logf("Initial Stop Loss: %v", initialStopLoss)

	// 模擬每個週期的價格變動
	for i := 1; i < len(data); i++ {
		period := data[i]

		// 更新 ATR
		err = s.UpdateATR(period.ATR)
		if err != nil {
			t.Logf("Period %d: Cannot update ATR (stop loss may be inactive)", period.Period)
			break
		}

		currentStopLoss, _ := s.GetStopLoss()

		t.Logf("Period %d: Price=%v, High=%v, Low=%v, ATR=%v, StopLoss=%v",
			period.Period, period.Close, period.High, period.Low,
			period.ATR, currentStopLoss)

		// 檢查是否觸發止損（使用最低價檢查）
		isTriggered, err := s.ShouldTriggerStopLoss(period.Low)
		if err != nil && err != stoploss.ErrStatusInvalid {
			t.Errorf("Period %d: Error checking trigger: %v", period.Period, err)
		}

		if isTriggered {
			t.Logf("Period %d: Stop Loss TRIGGERED at price %v (Stop Loss: %v)",
				period.Period, period.Low, currentStopLoss)
			break
		}
	}

	if !triggered {
		t.Log("Stop Loss was not triggered during the simulation")
	}
}

// 測試使用 Mock 歷史數據：震盪市場不觸發止損
func TestATRStopLoss_WithHistoricalData_Consolidation(t *testing.T) {
	data := GetMockConsolidationData()

	entryData := data[0]
	entryPrice := entryData.Close
	multiplier := d(3) // 較大的倍數避免震盪觸發

	s, err := NewATRStop(entryPrice, entryData.ATR, multiplier, nil)
	if err != nil {
		t.Fatalf("Failed to create ATR Stop Loss: %v", err)
	}

	t.Logf("Consolidation Test - Entry Price: %v, Initial ATR: %v, Multiplier: %v",
		entryPrice, entryData.ATR, multiplier)

	triggerCount := 0

	for i := 1; i < len(data); i++ {
		period := data[i]

		s.UpdateATR(period.ATR)
		currentStopLoss, _ := s.GetStopLoss()

		t.Logf("Period %d: Close=%v, ATR=%v, StopLoss=%v",
			period.Period, period.Close, period.ATR, currentStopLoss)

		isTriggered, _ := s.ShouldTriggerStopLoss(period.Low)
		if isTriggered {
			triggerCount++
			t.Logf("Period %d: Triggered at %v", period.Period, period.Low)
			break
		}
	}

	if triggerCount > 0 {
		t.Errorf("Stop Loss should not trigger in consolidation with multiplier=%v", multiplier)
	} else {
		t.Log("✓ Stop Loss correctly not triggered during consolidation")
	}
}

// 測試使用 Mock 歷史數據：上升趨勢不觸發止損
func TestATRStopLoss_WithHistoricalData_UpTrend(t *testing.T) {
	data := GetMockTrendingData()

	entryData := data[0]
	entryPrice := entryData.Close
	multiplier := d(2)

	s, err := NewATRStop(entryPrice, entryData.ATR, multiplier, nil)
	if err != nil {
		t.Fatalf("Failed to create ATR Stop Loss: %v", err)
	}

	t.Logf("Uptrend Test - Entry Price: %v, Initial ATR: %v, Multiplier: %v",
		entryPrice, entryData.ATR, multiplier)

	initialStopLoss, _ := s.GetStopLoss()
	t.Logf("Initial Stop Loss: %v", initialStopLoss)

	for i := 1; i < len(data); i++ {
		period := data[i]

		s.UpdateATR(period.ATR)
		currentStopLoss, _ := s.GetStopLoss()

		t.Logf("Period %d: Close=%v, ATR=%v, StopLoss=%v, PnL=%.2f%%",
			period.Period, period.Close, period.ATR, currentStopLoss,
			period.Close.Sub(entryPrice).Div(entryPrice).Mul(d(100)).InexactFloat64())

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
}

// 基準測試：創建 ATR 止損
func BenchmarkNewATRStop(b *testing.B) {
	entryPrice := d(100)
	atr := d(2.5)
	multiplier := d(2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewATRStop(entryPrice, atr, multiplier, nil)
	}
}

// 基準測試：計算止損
func BenchmarkCalculateStopLoss(b *testing.B) {
	entryPrice := d(100)
	atr := d(2.5)
	multiplier := d(2)
	s, _ := NewATRStop(entryPrice, atr, multiplier, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.CalculateStopLoss(entryPrice)
	}
}

// 基準測試：檢查觸發
func BenchmarkShouldTriggerStopLoss(b *testing.B) {
	entryPrice := d(100)
	atr := d(2.5)
	multiplier := d(2)
	currentPrice := d(94)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, _ := NewATRStop(entryPrice, atr, multiplier, nil)
		_, _ = s.ShouldTriggerStopLoss(currentPrice)
	}
}

// 基準測試：更新 ATR
func BenchmarkUpdateATR(b *testing.B) {
	entryPrice := d(100)
	atr := d(2.5)
	multiplier := d(2)
	s, _ := NewATRStop(entryPrice, atr, multiplier, nil)
	newATR := d(3.0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.UpdateATR(newATR)
	}
}
