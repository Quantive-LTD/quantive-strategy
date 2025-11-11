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

package coinbase

import (
	"context"
	"testing"

	"github.com/wang900115/quant/model"
	"github.com/wang900115/quant/model/currency"
	"github.com/wang900115/quant/model/trade"
)

func TestGetPrice(t *testing.T) {
	cb := New()
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
	cb := New()
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
	cb := New()
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

func TestWebSocketConnect(t *testing.T) {
	bc := NewStreamClient()
	err := bc.Connect()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer bc.Close()

	if bc.client == nil {
		t.Fatal("expected client to be initialized, got nil")
	}
	t.Logf("Successfully connected to Spot WebSocket")
}

func TestWebSocketSubscribe(t *testing.T) {
	bc := NewStreamClient()
	err := bc.Connect()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer bc.Close()
	err = bc.Subscribe(context.Background(), model.TradingPair{
		Base:     currency.BTCSymbol,
		Quote:    currency.USDTSymbol,
		Category: trade.SPOT,
	}, []string{"ticker"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
