package strategy

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
)

var maCallback = func(reason string) error {
	fmt.Println("Callback:", reason)
	return nil
}

func TestMovingAverageStopLoss(t *testing.T) {
	entry := decimal.NewFromFloat(100)
	initialMA := decimal.NewFromFloat(95)
	offsetPercent := decimal.NewFromFloat(0.02) // 2% below MA

	s, err := NewMovingAverageStop(entry, initialMA, offsetPercent, maCallback)
	if err != nil {
		t.Fatal("Failed to create Moving Average Stop Loss:", err)
	}

	// Test initial stop loss calculation
	stop, err := s.CalculateStopLoss(decimal.NewFromFloat(100))
	if err != nil {
		t.Fatal("Failed to calculate stop loss:", err)
	}
	expected := decimal.NewFromFloat(93.1) // 95 * (1 - 0.02) = 93.1
	if !stop.Equal(expected) {
		t.Errorf("Expected stop loss %v, got %v", expected, stop)
	}
	fmt.Println("Initial Stop Loss:", stop)

	// Test MA update
	s.SetMA(decimal.NewFromFloat(98))
	stop, err = s.GetStopLoss()
	if err != nil {
		t.Fatal("Failed to get stop loss:", err)
	}
	newExpected := decimal.NewFromFloat(96.04) // 98 * (1 - 0.02) = 96.04
	if !stop.Equal(newExpected) {
		t.Errorf("Expected stop loss after MA update %v, got %v", newExpected, stop)
	}
	fmt.Println("Stop Loss after MA update:", stop)

	// Test trigger condition
	triggered, err := s.ShouldTriggerStopLoss(decimal.NewFromFloat(96))
	if err != nil {
		t.Fatal("Failed to check if stop loss triggered:", err)
	}
	if !triggered {
		t.Error("Expected stop loss to be triggered")
	}

	// Test reset
	err = s.ReSet(decimal.NewFromFloat(105))
	if err != nil {
		t.Fatal("Failed to reset Moving Average Stop Loss:", err)
	}
}

func TestNewMovingAverageStop(t *testing.T) {
	tests := []struct {
		name          string
		entryPrice    decimal.Decimal
		initialMA     decimal.Decimal
		offsetPercent decimal.Decimal
		expectError   bool
		errorMsg      string
	}{
		{
			name:          "Valid parameters",
			entryPrice:    decimal.NewFromFloat(100),
			initialMA:     decimal.NewFromFloat(95),
			offsetPercent: decimal.NewFromFloat(0.02),
			expectError:   false,
		},
		{
			name:          "Zero moving average",
			entryPrice:    decimal.NewFromFloat(100),
			initialMA:     decimal.Zero,
			offsetPercent: decimal.NewFromFloat(0.02),
			expectError:   true,
			errorMsg:      "moving average value must be greater than 0",
		},
		{
			name:          "Negative moving average",
			entryPrice:    decimal.NewFromFloat(100),
			initialMA:     decimal.NewFromFloat(-5),
			offsetPercent: decimal.NewFromFloat(0.02),
			expectError:   true,
			errorMsg:      "moving average value must be greater than 0",
		},
		{
			name:          "Invalid offset percentage - negative",
			entryPrice:    decimal.NewFromFloat(100),
			initialMA:     decimal.NewFromFloat(95),
			offsetPercent: decimal.NewFromFloat(-0.1),
			expectError:   true,
			errorMsg:      "offset percentage must be between 0 and 1",
		},
		{
			name:          "Invalid offset percentage - greater than 1",
			entryPrice:    decimal.NewFromFloat(100),
			initialMA:     decimal.NewFromFloat(95),
			offsetPercent: decimal.NewFromFloat(1.5),
			expectError:   true,
			errorMsg:      "offset percentage must be between 0 and 1",
		},
		{
			name:          "Zero offset percentage",
			entryPrice:    decimal.NewFromFloat(100),
			initialMA:     decimal.NewFromFloat(95),
			offsetPercent: decimal.Zero,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewMovingAverageStop(tt.entryPrice, tt.initialMA, tt.offsetPercent, maCallback)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if s == nil {
				t.Error("Expected valid stop loss instance, got nil")
			}
		})
	}
}

func TestMovingAverageSetMA(t *testing.T) {
	entry := decimal.NewFromFloat(100)
	initialMA := decimal.NewFromFloat(95)
	offsetPercent := decimal.NewFromFloat(0.02) // 2% below MA

	s, err := NewMovingAverageStop(entry, initialMA, offsetPercent, maCallback)
	if err != nil {
		t.Fatal("Failed to create Moving Average Stop Loss:", err)
	}

	// Test that stop loss only moves up, not down
	s.SetMA(decimal.NewFromFloat(90)) // Lower MA
	stop, err := s.GetStopLoss()
	if err != nil {
		t.Fatal("Failed to get stop loss:", err)
	}

	// Stop loss should remain at previous higher level
	expected := decimal.NewFromFloat(93.1) // Original: 95 * (1 - 0.02) = 93.1
	if !stop.Equal(expected) {
		t.Errorf("Expected stop loss to remain at %v, got %v", expected, stop)
	}

	// Test that stop loss moves up with higher MA
	s.SetMA(decimal.NewFromFloat(100)) // Higher MA
	stop, err = s.GetStopLoss()
	if err != nil {
		t.Fatal("Failed to get stop loss:", err)
	}

	newExpected := decimal.NewFromFloat(98) // 100 * (1 - 0.02) = 98
	if !stop.Equal(newExpected) {
		t.Errorf("Expected stop loss to move up to %v, got %v", newExpected, stop)
	}
}

func TestMovingAverageWithoutOffset(t *testing.T) {
	entry := decimal.NewFromFloat(100)
	initialMA := decimal.NewFromFloat(95)

	s, err := NewMovingAverageStop(entry, initialMA, decimal.Zero, maCallback)
	if err != nil {
		t.Fatal("Failed to create Moving Average Stop Loss:", err)
	}

	// Without offset, stop loss should equal MA
	stop, err := s.GetStopLoss()
	if err != nil {
		t.Fatal("Failed to get stop loss:", err)
	}

	if !stop.Equal(initialMA) {
		t.Errorf("Expected stop loss to equal MA %v, got %v", initialMA, stop)
	}
}
