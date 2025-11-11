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

package binance

import (
	"context"
	"testing"

	"github.com/wang900115/quant/model"
	"github.com/wang900115/quant/model/currency"
	"github.com/wang900115/quant/model/trade"
)

func TestGetPrice(t *testing.T) {
	bc := New()
	t.Run("SPOT", func(t *testing.T) {
		pair := model.TradingPair{
			Base:     currency.BTCSymbol,
			Quote:    currency.USDTSymbol,
			Category: trade.SPOT,
		}
		price, err := bc.GetPrice(context.Background(), pair)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		t.Logf("SPOT Price: %s", price.NewPrice.String())
	})
	t.Run("FUTURES", func(t *testing.T) {
		pair := model.TradingPair{
			Base:     currency.BTCSymbol,
			Quote:    currency.USDTSymbol,
			Category: trade.FUTURES,
		}
		price, err := bc.GetPrice(context.Background(), pair)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		t.Logf("FUTURES Price: %s", price.NewPrice.String())
	})
	t.Run("INVERSE", func(t *testing.T) {
		pair := model.TradingPair{
			Base:     currency.BTCSymbol,
			Quote:    currency.USDTSymbol,
			Category: trade.INVERSE,
		}
		price, err := bc.GetPrice(context.Background(), pair)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		t.Logf("INVERSE Price: %s", price.NewPrice.String())
	})
}

func TestGetKlines(t *testing.T) {
	bc := New()
	t.Run("SPOT", func(t *testing.T) {
		pair := model.TradingPair{
			Base:     currency.BTCSymbol,
			Quote:    currency.USDTSymbol,
			Category: trade.SPOT,
		}
		intervals, err := bc.GetKlines(context.Background(), pair, "1h", 10)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(intervals) == 0 {
			t.Fatal("expected intervals, got empty slice")
		}
	})
	t.Run("FUTURES", func(t *testing.T) {
		pair := model.TradingPair{
			Base:     currency.BTCSymbol,
			Quote:    currency.USDTSymbol,
			Category: trade.FUTURES,
		}
		intervals, err := bc.GetKlines(context.Background(), pair, "15m", 5)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(intervals) == 0 {
			t.Fatal("expected intervals, got empty slice")
		}
	})
	t.Run("INVERSE", func(t *testing.T) {
		pair := model.TradingPair{
			Base:     currency.BTCSymbol,
			Quote:    currency.USDTSymbol,
			Category: trade.INVERSE,
		}
		intervals, err := bc.GetKlines(context.Background(), pair, "15m", 5)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(intervals) == 0 {
			t.Fatal("expected intervals, got empty slice")
		}
	})
}

func TestGetOrderBook(t *testing.T) {
	bc := New()
	t.Run("SPOT", func(t *testing.T) {
		pair := model.TradingPair{
			Base:     currency.BTCSymbol,
			Quote:    currency.USDTSymbol,
			Category: trade.SPOT,
		}
		orderBook, err := bc.GetOrderBook(context.Background(), pair, 5)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		t.Logf("Order Book: %+v", orderBook)
	})
	t.Run("FUTURES", func(t *testing.T) {
		pair := model.TradingPair{
			Base:     currency.BTCSymbol,
			Quote:    currency.USDTSymbol,
			Category: trade.FUTURES,
		}
		orderBook, err := bc.GetOrderBook(context.Background(), pair, 5)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		t.Logf("Order Book: %+v", orderBook)
	})
	t.Run("INVERSE", func(t *testing.T) {
		pair := model.TradingPair{
			Base:     currency.BTCSymbol,
			Quote:    currency.USDTSymbol,
			Category: trade.INVERSE,
		}
		orderBook, err := bc.GetOrderBook(context.Background(), pair, 5)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		t.Logf("Order Book: %+v", orderBook)
	})
}

func TestWebSocketConnection(t *testing.T) {
	testCases := []struct {
		name     string
		category trade.Category
	}{
		{name: "SPOT WebSocket", category: trade.SPOT},
		{name: "FUTURES WebSocket", category: trade.FUTURES},
		{name: "INVERSE WebSocket", category: trade.INVERSE},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := NewStreamClient()
			err := client.ConnectByCategory(tc.category)
			if err != nil {
				t.Fatalf("failed to connect: %v", err)
			}
			defer client.Close()

			if client.client == nil {
				t.Fatal("connection is nil")
			}
			t.Logf("Successfully connected to %s WebSocket", tc.name)
		})
	}
}

func TestStreamClientConnect(t *testing.T) {
	testCases := []struct {
		name     string
		endpoint string
	}{
		{name: "SPOT Endpoint", endpoint: spotWsEndpoint},
		{name: "FUTURES Endpoint", endpoint: futuresWsEndpoint},
		{name: "INVERSE Endpoint", endpoint: inverseWsEndpoint},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := NewStreamClient()
			err := client.connect(tc.endpoint)
			if err != nil {
				t.Fatalf("failed to connect to %s: %v", tc.endpoint, err)
			}
			defer client.Close()

			t.Logf("Successfully connected to %s", tc.endpoint)
		})
	}
}

func TestWebSocketSubscribe(t *testing.T) {
	pair := model.TradingPair{
		Base:     currency.BTCSymbol,
		Quote:    currency.USDTSymbol,
		Category: trade.SPOT,
	}

	testCases := []struct {
		name       string
		streamType string
	}{
		{name: "Ticker Stream", streamType: "ticker"},
		{name: "Mini Ticker Stream", streamType: "miniTicker"},
		{name: "Trade Stream", streamType: "trade"},
		{name: "Aggregate Trade Stream", streamType: "aggTrade"},
		{name: "Kline 1m Stream", streamType: "kline_1m"},
		{name: "Depth Stream", streamType: "depth"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := NewStreamClient()
			err := client.ConnectByCategory(pair.Category)
			if err != nil {
				t.Fatalf("failed to connect: %v", err)
			}
			defer client.Close()

			// Subscribe using the client's method
			ctx := context.Background()
			if err := client.Subscribe(ctx, pair, tc.streamType); err != nil {
				t.Fatalf("failed to subscribe to %s: %v", tc.streamType, err)
			}

			t.Logf("Successfully subscribed to %s", tc.streamType)
		})
	}
}
