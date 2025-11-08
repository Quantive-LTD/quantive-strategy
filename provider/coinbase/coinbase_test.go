//go:build integration
// +build integration

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

package coinbase

import (
	"context"
	"testing"

	"github.com/wang900115/quant/model"
	"github.com/wang900115/quant/model/currency"
	"github.com/wang900115/quant/model/trade"
)

func TestGetPrice(t *testing.T) {
	cb := NewClient()
	t.Run("SPOT", func(t *testing.T) {
		pair := model.TradingPair{
			Base:     currency.BTCSymbol,
			Quote:    currency.USDTSymbol,
			Category: trade.SPOT,
		}
		price, err := cb.GetPrice(context.Background(), pair)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		t.Logf("SPOT Price: %s", price.NewPrice.String())
	})
}

func TestGetKlines(t *testing.T) {
	cb := NewClient()

	t.Run("SPOT", func(t *testing.T) {
		pair := model.TradingPair{
			Base:     currency.BTCSymbol,
			Quote:    currency.USDTSymbol,
			Category: trade.SPOT,
		}
		intervals, err := cb.GetKlines(context.Background(), pair, "1h", 10)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(intervals) == 0 {
			t.Fatal("expected intervals, got empty slice")
		}
	})

}

func TestGetOrderBook(t *testing.T) {
	cb := NewClient()
	t.Run("SPOT", func(t *testing.T) {
		pair := model.TradingPair{
			Base:     currency.BTCSymbol,
			Quote:    currency.USDTSymbol,
			Category: trade.SPOT,
		}
		orderBook, err := cb.GetOrderBook(context.Background(), pair, 5)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		t.Logf("Order Book: %+v", orderBook)
	})
}
