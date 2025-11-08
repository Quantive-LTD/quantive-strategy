// Copyright 2025 Perry. All rights reserved.

// Licensed MIT License

// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// https://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package manager

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/model"
	"github.com/wang900115/quant/stoploss/strategy"
)

// ExampleManagerUsage demonstrates how to use the Manager with result handling
func ExampleManagerUsage() {
	fmt.Println("=== Strategy Manager Usage Example ===")

	// Create manager configuration
	config := ManagerConfig{
		BufferSize:    1000,
		ReadTimeout:   time.Second * 5,
		CheckInterval: time.Millisecond * 100,
	}

	// Create the manager
	manager := New(config)

	// Entry price for all strategies
	entryPrice := decimal.NewFromFloat(100.0)

	// Create callback function for all strategies
	callback := func(reason string) error {
		fmt.Printf("Strategy triggered: %s\n", reason)
		return nil
	}

	fmt.Println("Registering strategies...")

	// Register Fixed Stop Loss strategies
	percentStopStrategy, _ := strategy.NewFixedPercentStop(
		entryPrice,
		decimal.NewFromFloat(0.05), // 5% stop loss
		callback,
	)
	manager.RegisterStrategy("Fixed-Percent-Stop-5%", percentStopStrategy)

	atrStopStrategy, _ := strategy.NewFixedATRStop(
		entryPrice,
		decimal.NewFromFloat(2.5), // ATR value
		decimal.NewFromFloat(2.0), // ATR multiplier
		callback,
	)
	manager.RegisterStrategy("Fixed-ATR-Stop-2x", atrStopStrategy)

	// Register Timed Stop Loss strategies
	timedTrailingStrategy, _ := strategy.NewTrailingTimeBasedStop(
		entryPrice,
		decimal.NewFromFloat(0.04), // 4% tolerance
		int64(300),                 // 5 minutes time threshold
		callback,
	)
	manager.RegisterStrategy("Timed-Trailing-4%", timedTrailingStrategy)

	// Register Fixed Take Profit strategies
	percentProfitStrategy, _ := strategy.NewFixedPercentProfit(
		entryPrice,
		decimal.NewFromFloat(0.08), // 8% take profit
		callback,
	)
	manager.RegisterStrategy("Fixed-Percent-Profit-8%", percentProfitStrategy)

	atrProfitStrategy, _ := strategy.NewFixedATRProfit(
		entryPrice,
		decimal.NewFromFloat(2.5), // ATR value
		decimal.NewFromFloat(3.0), // ATR multiplier for profit
		callback,
	)
	manager.RegisterStrategy("Fixed-ATR-Profit-3x", atrProfitStrategy)

	// Register Hybrid strategies
	riskRewardStrategy, _ := strategy.NewRiskRewardRatio(
		entryPrice,
		decimal.NewFromFloat(0.03), // 3% risk
		decimal.NewFromFloat(0.06), // 6% reward (1:2 ratio)
		callback,
	)
	manager.RegisterStrategy("Risk-Reward-1:2", riskRewardStrategy)

	structureSwingStrategy, _ := strategy.NewStructureSwingStop(
		20, // lookback period
		entryPrice,
		decimal.NewFromFloat(0.01), // 1% swing distance
		decimal.NewFromFloat(0.05), // 5% stop loss
		decimal.NewFromFloat(0.10), // 10% take profit
		true,                       // long position
		callback,
	)
	manager.RegisterStrategy("Structure-Swing", structureSwingStrategy)

	// Start the manager (launches 6 goroutines)
	fmt.Println("Starting manager with 6 goroutines...")
	manager.Start()

	// Simulate price movements
	fmt.Println("\nSimulating price movements...")
	prices := []float64{100.0, 102.0, 101.0, 99.0, 98.5, 97.0, 95.0, 93.0, 105.0, 108.0}

	for i, price := range prices {
		fmt.Printf("\n--- Price Update #%d: $%.2f ---\n", i+1, price)
		currentPrice := decimal.NewFromFloat(price)

		// Create price point
		pricePoint := model.PricePoint{
			NewPrice:  currentPrice,
			UpdatedAt: time.Now(),
		}

		// Send price update to all goroutines with error callback
		manager.Collect(pricePoint, func() {
			fmt.Printf("Warning: Channel full for price $%.2f\n", price)
		})

		// Small delay to see the processing
		time.Sleep(200 * time.Millisecond)
	}

	// Let the system process for a bit more
	fmt.Println("\nLetting system process for 2 more seconds...")
	time.Sleep(2 * time.Second)

	// Stop the manager
	fmt.Println("Stopping manager...")
	manager.Stop()

	fmt.Println("Strategy Manager example completed!")
}

// AdvancedResultProcessingExample shows advanced result processing
func AdvancedResultProcessingExample() {
	fmt.Println("\n=== Advanced Result Processing Example ===")

	config := ManagerConfig{
		BufferSize:    500,
		ReadTimeout:   time.Second * 3,
		CheckInterval: time.Millisecond * 50,
	}

	manager := New(config)

	// Register some strategies
	entryPrice := decimal.NewFromFloat(100.0)
	callback := func(reason string) error { return nil }

	// Register multiple strategies for better demo
	for i := 0; i < 3; i++ {
		stopLoss := decimal.NewFromFloat(0.03 + float64(i)*0.01)   // 3%, 4%, 5%
		takeProfit := decimal.NewFromFloat(0.06 + float64(i)*0.02) // 6%, 8%, 10%

		percentStop, _ := strategy.NewFixedPercentStop(entryPrice, stopLoss, callback)
		manager.RegisterStrategy(fmt.Sprintf("Stop-%d", i+1), percentStop)

		percentProfit, _ := strategy.NewFixedPercentProfit(entryPrice, takeProfit, callback)
		manager.RegisterStrategy(fmt.Sprintf("Profit-%d", i+1), percentProfit)

		riskReward, _ := strategy.NewRiskRewardRatio(entryPrice, stopLoss, takeProfit, callback)
		manager.RegisterStrategy(fmt.Sprintf("RiskReward-%d", i+1), riskReward)
	}

	manager.Start()

	// Send rapid price updates
	fmt.Println("Sending rapid price updates...")
	for i := 0; i < 50; i++ {
		// Create volatile price movement
		price := 100.0 + 20.0*(float64(i%10)/10.0) - 10.0 // Oscillate between 90-110

		pricePoint := model.PricePoint{
			NewPrice:  decimal.NewFromFloat(price),
			UpdatedAt: time.Now(),
		}

		manager.Collect(pricePoint, func() {
			fmt.Printf("Channel full at update %d\n", i)
		})

		if i%10 == 0 {
			fmt.Printf("Sent %d price updates...\n", i)
		}

		time.Sleep(50 * time.Millisecond)
	}

	// Wait for processing
	time.Sleep(1 * time.Second)

	// Print statistics
	stats := manager.Snapshot()
	fmt.Println("\n=== Processing Statistics ===")
	for key, value := range stats {
		fmt.Printf("%s: %d\n", key, value)
	}
	manager.Stop()
	fmt.Println("Advanced result processing example completed!")
}
