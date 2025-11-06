package strategy

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/stoploss"
)

var swingCallback stoploss.DefaultCallback

func TestNewStructureSwingStop(t *testing.T) {
	tests := []struct {
		name           string
		entryPrice     decimal.Decimal
		lookbackPeriod int
		swingDistance  decimal.Decimal
		isLong         bool
		expectError    bool
	}{
		{
			name:           "Valid long position",
			entryPrice:     decimal.NewFromInt(100),
			lookbackPeriod: 5,
			swingDistance:  decimal.NewFromFloat(0.02),
			isLong:         true,
			expectError:    false,
		},
		{
			name:           "Valid short position",
			entryPrice:     decimal.NewFromInt(100),
			lookbackPeriod: 10,
			swingDistance:  decimal.NewFromFloat(0.01),
			isLong:         false,
			expectError:    false,
		},
		{
			name:           "Invalid lookback period",
			entryPrice:     decimal.NewFromInt(100),
			lookbackPeriod: 0,
			swingDistance:  decimal.NewFromFloat(0.02),
			isLong:         true,
			expectError:    true,
		},
		{
			name:           "Invalid swing distance",
			entryPrice:     decimal.NewFromInt(100),
			lookbackPeriod: 5,
			swingDistance:  decimal.NewFromInt(0),
			isLong:         true,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy, err := NewStructureSwingStop(tt.entryPrice, tt.lookbackPeriod, tt.swingDistance, tt.isLong, swingCallback)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if strategy == nil {
				t.Errorf("Expected strategy but got nil")
			}
		})
	}
}

func TestStructureSwingStopLoss_LongPosition(t *testing.T) {
	entryPrice := decimal.NewFromInt(100)
	strategy, err := NewStructureSwingStop(entryPrice, 3, decimal.NewFromFloat(0.02), true, swingCallback)
	if err != nil {
		t.Fatalf("Failed to create strategy: %v", err)
	}

	// Test initial stop loss calculation
	stopLoss, err := strategy.CalculateStopLoss(entryPrice)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if stopLoss.LessThanOrEqual(decimal.Zero) {
		t.Errorf("Expected stop loss to be greater than zero, got %v", stopLoss)
	}

	// Simulate price movements to create swing points
	prices := []decimal.Decimal{
		decimal.NewFromInt(102), // up
		decimal.NewFromInt(104), // up
		decimal.NewFromInt(103), // down
		decimal.NewFromInt(101), // down (swing low)
		decimal.NewFromInt(105), // up
		decimal.NewFromInt(107), // up (swing high)
		decimal.NewFromInt(106), // down
	}

	for _, price := range prices {
		_, err := strategy.CalculateStopLoss(price)
		if err != nil {
			t.Errorf("Unexpected error calculating stop loss at price %v: %v", price, err)
		}
	}

	// Test stop loss trigger
	triggered, err := strategy.ShouldTriggerStopLoss(decimal.NewFromInt(90)) // Below stop loss
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !triggered {
		t.Errorf("Expected stop loss to be triggered")
	}
}

func TestStructureSwingStopLoss_ShortPosition(t *testing.T) {
	entryPrice := decimal.NewFromInt(100)
	strategy, err := NewStructureSwingStop(entryPrice, 3, decimal.NewFromFloat(0.02), false, swingCallback)
	if err != nil {
		t.Fatalf("Failed to create strategy: %v", err)
	}

	// Test initial stop loss calculation
	stopLoss, err := strategy.CalculateStopLoss(entryPrice)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if stopLoss.LessThanOrEqual(decimal.Zero) {
		t.Errorf("Expected stop loss to be greater than zero, got %v", stopLoss)
	}

	// Simulate price movements to create swing points
	prices := []decimal.Decimal{
		decimal.NewFromInt(98), // down
		decimal.NewFromInt(96), // down
		decimal.NewFromInt(97), // up
		decimal.NewFromInt(99), // up (swing high)
		decimal.NewFromInt(95), // down
		decimal.NewFromInt(93), // down (swing low)
		decimal.NewFromInt(94), // up
	}

	for _, price := range prices {
		_, err := strategy.CalculateStopLoss(price)
		if err != nil {
			t.Errorf("Unexpected error calculating stop loss at price %v: %v", price, err)
		}
	}

	// Test stop loss trigger
	triggered, err := strategy.ShouldTriggerStopLoss(decimal.NewFromInt(90)) // Above stop loss for short
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !triggered {
		t.Errorf("Expected stop loss to be triggered")
	}
}

func TestStructureSwingTakeProfit(t *testing.T) {
	entryPrice := decimal.NewFromInt(100)
	strategy, err := NewStructureSwingStop(entryPrice, 3, decimal.NewFromFloat(0.02), true, swingCallback)
	if err != nil {
		t.Fatalf("Failed to create strategy: %v", err)
	}

	// Test take profit calculation
	takeProfit, err := strategy.CalculateTakeProfit(entryPrice)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if takeProfit.LessThanOrEqual(decimal.Zero) {
		t.Errorf("Expected take profit to be greater than zero, got %v", takeProfit)
	}

	// Test take profit trigger
	triggered, err := strategy.ShouldTriggerTakeProfit(decimal.NewFromInt(120)) // Above take profit
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !triggered {
		t.Errorf("Expected take profit to be triggered")
	}
}

func TestStructureSwingGetMethods(t *testing.T) {
	entryPrice := decimal.NewFromInt(100)
	strategy, err := NewStructureSwingStop(entryPrice, 3, decimal.NewFromFloat(0.02), true, swingCallback)
	if err != nil {
		t.Fatalf("Failed to create strategy: %v", err)
	}

	// Test GetStopLoss
	stopLoss, err := strategy.GetStopLoss()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if stopLoss.LessThanOrEqual(decimal.Zero) {
		t.Errorf("Expected stop loss to be greater than zero, got %v", stopLoss)
	}

	// Test GetTakeProfit
	takeProfit, err := strategy.GetTakeProfit()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if takeProfit.LessThanOrEqual(decimal.Zero) {
		t.Errorf("Expected take profit to be greater than zero, got %v", takeProfit)
	}
}

func TestStructureSwingReSet(t *testing.T) {
	entryPrice := decimal.NewFromInt(100)
	strategy, err := NewStructureSwingStop(entryPrice, 3, decimal.NewFromFloat(0.02), true, swingCallback)
	if err != nil {
		t.Fatalf("Failed to create strategy: %v", err)
	}

	// Reset with new price
	newPrice := decimal.NewFromInt(120)
	err = strategy.ReSet(newPrice)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify that stop loss and take profit are updated
	stopLoss, err := strategy.GetStopLoss()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	takeProfit, err := strategy.GetTakeProfit()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// For long position, stop loss should be below new entry, take profit above
	if stopLoss.GreaterThanOrEqual(newPrice) {
		t.Errorf("Expected stop loss %v to be below new entry price %v", stopLoss, newPrice)
	}
	if takeProfit.LessThanOrEqual(newPrice) {
		t.Errorf("Expected take profit %v to be above new entry price %v", takeProfit, newPrice)
	}
}

func TestStructureSwingDeactivate(t *testing.T) {
	entryPrice := decimal.NewFromInt(100)
	strategy, err := NewStructureSwingStop(entryPrice, 3, decimal.NewFromFloat(0.02), true, swingCallback)
	if err != nil {
		t.Fatalf("Failed to create strategy: %v", err)
	}

	// Deactivate the strategy
	err = strategy.Deactivate()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test that methods return errors when inactive
	_, err = strategy.CalculateStopLoss(entryPrice)
	if err != stoploss.ErrStatusInvalid {
		t.Errorf("Expected ErrStatusInvalid, got %v", err)
	}

	_, err = strategy.GetStopLoss()
	if err != stoploss.ErrStatusInvalid {
		t.Errorf("Expected ErrStatusInvalid, got %v", err)
	}
}
