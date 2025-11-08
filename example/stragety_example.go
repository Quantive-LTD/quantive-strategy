// Copyright 2025 Quantive. All rights reserved.

// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// https://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package example

import (
	"fmt"
	"log"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/model"
	"github.com/wang900115/quant/stoploss/engine"
	"github.com/wang900115/quant/stoploss/strategy"
)

// demonstrates usage with various strategies
func StragetyUsage() {
	fmt.Println("=== Strategy Usage Example ===")
	// Create  configuration
	config := engine.Config{
		BufferSize:    1000,
		ReadTimeout:   time.Second * 5,
		CheckInterval: time.Second * 3,
	}

	// Create the manager instance
	manager := engine.New(config)

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

	log.Printf("Total strategies registered: %+v\n", manager)
	// Start the manager (launches 6 goroutines)
	log.Println("Starting manager with 6 goroutines...")
	manager.Start()

	// Simulate price movements
	log.Println("\nSimulating price movements...")
	prices := []float64{100.0, 102.0, 101.0, 99.0, 98.5, 97.0, 95.0, 93.0, 105.0, 108.0}

	for i, price := range prices {
		log.Printf("\n--- Price Update #%d: $%.2f ---\n", i+1, price)
		currentPrice := decimal.NewFromFloat(price)

		// Create price point
		pricePoint := model.PricePoint{
			NewPrice:  currentPrice,
			UpdatedAt: time.Now(),
		}

		// Send price update to all goroutines with error callback
		manager.Collect(pricePoint, func() {
			log.Printf("Warning: Channel full for price $%.2f\n", price)
		})

		// Small delay to see the processing
		time.Sleep(200 * time.Millisecond)
	}

	// Let the system process for a bit more
	log.Println("\nLetting system process for 2 more seconds...")
	time.Sleep(30 * time.Second)

	// Stop the manager
	log.Println("Stopping manager...")
	manager.Stop()

	log.Println("Strategy Manager example completed!")
}
