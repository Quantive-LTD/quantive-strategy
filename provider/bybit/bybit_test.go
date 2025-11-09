//go:build integration
// +build integration

// Copyright (C) 2025 Quantive
//
// SPDX-License-Identifier: MIT OR AGPL-3.0-or-later
//
// This file is part of the Decision Engine project.
// You may choose to use this file under the terms of either
// the MIT License or the GNU Affero General Public License v3.0 or later.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the LICENSE files for more details.

package bybit

import (
	"context"
	"testing"

	"github.com/wang900115/quant/model"
	"github.com/wang900115/quant/model/currency"
	"github.com/wang900115/quant/model/trade"
)

func TestGetPrice(t *testing.T) {
	by := NewClient()
	t.Run("SPOT", func(t *testing.T) {
		pair := model.TradingPair{
			Base:     currency.BTCSymbol,
			Quote:    currency.USDTSymbol,
			Category: trade.SPOT,
		}
		price, err := by.GetPrice(context.Background(), pair)
		if err != nil {
			t.Fatalf("failed to get spot price: %v", err)
		}
		if price == nil {
			t.Fatal("expected price point, got nil")
		}
		t.Logf("Spot Price for %s-%s: %s", pair.Base, pair.Quote, price.NewPrice.String())
	})
	t.Run("FUTURES", func(t *testing.T) {
		pair := model.TradingPair{
			Base:     currency.BTCSymbol,
			Quote:    currency.USDTSymbol,
			Category: trade.FUTURES,
		}
		price, err := by.GetPrice(context.Background(), pair)
		if err != nil {
			t.Fatalf("failed to get futures price: %v", err)
		}
		if price == nil {
			t.Fatal("expected price point, got nil")
		}
		t.Logf("Futures Price for %s-%s: %s", pair.Base, pair.Quote, price.NewPrice.String())
	})
}

func TestGetKlines(t *testing.T) {
	by := NewClient()

	t.Run("SPOT", func(t *testing.T) {
		pair := model.TradingPair{
			Base:     currency.BTCSymbol,
			Quote:    currency.USDTSymbol,
			Category: trade.SPOT,
		}
		klines, err := by.GetKlines(context.Background(), pair, "1", 10)
		if err != nil {
			t.Fatalf("failed to get spot klines: %v", err)
		}
		if len(klines) == 0 {
			t.Fatal("expected klines, got none")
		}
		t.Logf("Spot Klines for %s-%s: %v", pair.Base, pair.Quote, klines)
	})
	t.Run("FUTURES", func(t *testing.T) {
		pair := model.TradingPair{
			Base:     currency.BTCSymbol,
			Quote:    currency.USDTSymbol,
			Category: trade.FUTURES,
		}
		klines, err := by.GetKlines(context.Background(), pair, "1", 10)
		if err != nil {
			t.Fatalf("failed to get futures klines: %v", err)
		}
		if len(klines) == 0 {
			t.Fatal("expected klines, got none")
		}
		t.Logf("Futures Klines for %s-%s: %v", pair.Base, pair.Quote, klines)
	})
}

func TestGetOrderBook(t *testing.T) {
	by := NewClient()
	t.Run("SPOT", func(t *testing.T) {
		pair := model.TradingPair{
			Base:     currency.BTCSymbol,
			Quote:    currency.USDTSymbol,
			Category: trade.SPOT,
		}
		orderBook, err := by.GetOrderBook(context.Background(), pair, 20)
		if err != nil {
			t.Fatalf("failed to get spot order book: %v", err)
		}
		if orderBook == nil {
			t.Fatal("expected order book, got nil")
		}
		t.Logf("Spot Order Book for %s-%s: %v", pair.Base, pair.Quote, orderBook)
	})
	t.Run("FUTURES", func(t *testing.T) {
		pair := model.TradingPair{
			Base:     currency.BTCSymbol,
			Quote:    currency.USDTSymbol,
			Category: trade.FUTURES,
		}
		orderBook, err := by.GetOrderBook(context.Background(), pair, 20)
		if err != nil {
			t.Fatalf("failed to get futures order book: %v", err)
		}
		if orderBook == nil {
			t.Fatal("expected order book, got nil")
		}
		t.Logf("Futures Order Book for %s-%s: %v", pair.Base, pair.Quote, orderBook)
	})
}
