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

package provider

import (
	"context"
	"testing"

	"github.com/wang900115/quant/model"
	"github.com/wang900115/quant/model/currency"
	"github.com/wang900115/quant/model/trade"
	"github.com/wang900115/quant/provider/binance"
	"github.com/wang900115/quant/provider/bybit"
	"github.com/wang900115/quant/provider/coinbase"
	"github.com/wang900115/quant/provider/okx"
)

func TestProvider(t *testing.T) {
	// This is a placeholder test to ensure the provider package builds correctly.
	ps := New()
	ps.Register(model.BINANCE, binance.NewClient())
	ps.Register(model.COINBASE, coinbase.NewClient())
	ps.Register(model.BYBIT, bybit.NewClient())
	ps.Register(model.OKX, okx.NewClient())

	pair := model.TradingPair{
		ExchangeID: model.BINANCE,
		Base:       currency.BTCSymbol,
		Quote:      currency.USDTSymbol,
	}

	p, err := ps.GetPrice(context.Background(), pair)
	if err != nil {
		t.Fatalf("Failed to get price: %v", err)
	}

	t.Logf("Price for %s: %+v", pair.Symbol(), *p)

	pair = model.TradingPair{
		ExchangeID: model.COINBASE,
		Base:       currency.BTCSymbol,
		Quote:      currency.USDTSymbol,
	}

	p, err = ps.GetPrice(context.Background(), pair)
	if err != nil {
		t.Fatalf("Failed to get price: %v", err)
	}
	t.Logf("Price for %s: %+v", pair.Symbol(), *p)

	pair = model.TradingPair{
		ExchangeID: model.BYBIT,
		Base:       currency.BTCSymbol,
		Quote:      currency.USDTSymbol,
	}

	p, err = ps.GetPrice(context.Background(), pair)
	if err != nil {
		t.Fatalf("Failed to get price: %v", err)
	}
	t.Logf("Price for %s: %+v", pair.Symbol(), *p)

	pair = model.TradingPair{
		ExchangeID: model.OKX,
		Base:       currency.BTCSymbol,
		Quote:      currency.USDTSymbol,
	}
	p, err = ps.GetPrice(context.Background(), pair)
	if err != nil {
		t.Fatalf("Failed to get price: %v", err)
	}
	t.Logf("Price for %s: %+v", pair.Symbol(), *p)
}

func TestGetKlines(t *testing.T) {
	ps := New()
	ps.Register(model.BINANCE, binance.NewClient())
	ps.Register(model.COINBASE, coinbase.NewClient())
	ps.Register(model.OKX, okx.NewClient())
	ps.Register(model.BYBIT, bybit.NewClient())

	t.Run("SPOT", func(t *testing.T) {

		pair := model.TradingPair{
			ExchangeID: model.BINANCE,
			Base:       currency.BTCSymbol,
			Quote:      currency.USDTSymbol,
			Category:   trade.SPOT,
		}
		klines, err := ps.GetKlines(context.Background(), pair, "1m", 5)
		if err != nil {
			t.Fatalf("Failed to get klines: %v", err)
		}
		t.Logf("Klines for %s: %+v", pair.Symbol(), klines)
		pair = model.TradingPair{
			ExchangeID: model.OKX,
			Base:       currency.ETHSymbol,
			Quote:      currency.USDTSymbol,
			Category:   trade.SPOT,
		}
		klines, err = ps.GetKlines(context.Background(), pair, "5m", 10)
		if err != nil {
			t.Fatalf("Failed to get klines: %v", err)
		}
		t.Logf("Klines for %s: %+v", pair.Symbol(), klines)
		pair = model.TradingPair{
			ExchangeID: model.COINBASE,
			Base:       currency.BTCSymbol,
			Quote:      currency.USDTSymbol,
			Category:   trade.SPOT,
		}
		klines, err = ps.GetKlines(context.Background(), pair, "1h", 10)
		if err != nil {
			t.Fatalf("Failed to get klines: %v", err)
		}
		t.Logf("Klines for %s: %+v", pair.Symbol(), klines)
		pair = model.TradingPair{
			ExchangeID: model.BYBIT,
			Base:       currency.BTCSymbol,
			Quote:      currency.USDTSymbol,
			Category:   trade.SPOT,
		}
		klines, err = ps.GetKlines(context.Background(), pair, "15m", 20)
		if err != nil {
			t.Fatalf("Failed to get klines: %v", err)
		}
		t.Logf("Klines for %s: %+v", pair.Symbol(), klines)
	})

	t.Run("FUTURES", func(t *testing.T) {

		pair := model.TradingPair{
			ExchangeID: model.BINANCE,
			Base:       currency.BTCSymbol,
			Quote:      currency.USDTSymbol,
			Category:   trade.FUTURES,
		}
		klines, err := ps.GetKlines(context.Background(), pair, "1m", 5)
		if err != nil {
			t.Fatalf("Failed to get klines: %v", err)
		}
		t.Logf("Klines for %s: %+v", pair.Symbol(), klines)
		pair = model.TradingPair{
			ExchangeID: model.OKX,
			Base:       currency.ETHSymbol,
			Quote:      currency.USDTSymbol,
			Category:   trade.FUTURES,
		}
		klines, err = ps.GetKlines(context.Background(), pair, "5m", 10)
		if err != nil {
			t.Fatalf("Failed to get klines: %v", err)
		}
		t.Logf("Klines for %s: %+v", pair.Symbol(), klines)
		pair = model.TradingPair{
			ExchangeID: model.COINBASE,
			Base:       currency.BTCSymbol,
			Quote:      currency.USDTSymbol,
			Category:   trade.FUTURES,
		}
		klines, err = ps.GetKlines(context.Background(), pair, "1h", 10)
		if err != nil {
			t.Fatalf("Failed to get klines: %v", err)
		}
		t.Logf("Klines for %s: %+v", pair.Symbol(), klines)
		pair = model.TradingPair{
			ExchangeID: model.BYBIT,
			Base:       currency.BTCSymbol,
			Quote:      currency.USDTSymbol,
			Category:   trade.FUTURES,
		}
		klines, err = ps.GetKlines(context.Background(), pair, "15m", 20)
		if err != nil {
			t.Fatalf("Failed to get klines: %v", err)
		}
		t.Logf("Klines for %s: %+v", pair.Symbol(), klines)
	})
}

func TestOrderBook(t *testing.T) {
	ps := New()
	ps.Register(model.COINBASE, coinbase.NewClient())
	ps.Register(model.BINANCE, binance.NewClient())
	ps.Register(model.OKX, okx.NewClient())
	ps.Register(model.BYBIT, bybit.NewClient())

	t.Run("SPOT", func(t *testing.T) {
		pair := model.TradingPair{
			ExchangeID: model.BINANCE,
			Base:       currency.BTCSymbol,
			Quote:      currency.USDTSymbol,
			Category:   trade.SPOT,
		}
		orderBook, err := ps.GetOrderBook(context.Background(), pair, 5)
		if err != nil {
			t.Fatalf("Failed to get order book: %v", err)
		}
		t.Logf("Order book for %s: %+v", pair.Symbol(), orderBook)
		pair = model.TradingPair{
			ExchangeID: model.COINBASE,
			Base:       currency.BTCSymbol,
			Quote:      currency.USDTSymbol,
			Category:   trade.SPOT,
		}
		orderBook, err = ps.GetOrderBook(context.Background(), pair, 5)
		if err != nil {
			t.Fatalf("Failed to get order book: %v", err)
		}
		t.Logf("Order book for %s: %+v", pair.Symbol(), orderBook)
		pair = model.TradingPair{
			ExchangeID: model.BYBIT,
			Base:       currency.BTCSymbol,
			Quote:      currency.USDTSymbol,
			Category:   trade.SPOT,
		}
		orderBook, err = ps.GetOrderBook(context.Background(), pair, 5)
		if err != nil {
			t.Fatalf("Failed to get order book: %v", err)
		}
		t.Logf("Order book for %s: %+v", pair.Symbol(), orderBook)
		pair = model.TradingPair{
			ExchangeID: model.OKX,
			Base:       currency.ETHSymbol,
			Quote:      currency.USDTSymbol,
			Category:   trade.SPOT,
		}
		orderBook, err = ps.GetOrderBook(context.Background(), pair, 5)
		if err != nil {
			t.Fatalf("Failed to get order book: %v", err)
		}
		t.Logf("Order book for %s: %+v", pair.Symbol(), orderBook)
	})

	t.Run("FUTURES", func(t *testing.T) {

		pair := model.TradingPair{
			ExchangeID: model.BINANCE,
			Base:       currency.BTCSymbol,
			Quote:      currency.USDTSymbol,
			Category:   trade.FUTURES,
		}
		orderBook, err := ps.GetOrderBook(context.Background(), pair, 5)
		if err != nil {
			t.Fatalf("Failed to get order book: %v", err)
		}
		t.Logf("Order book for %s: %+v", pair.Symbol(), orderBook)
		pair = model.TradingPair{
			ExchangeID: model.COINBASE,
			Base:       currency.BTCSymbol,
			Quote:      currency.USDTSymbol,
			Category:   trade.FUTURES,
		}
		orderBook, err = ps.GetOrderBook(context.Background(), pair, 5)
		if err != nil {
			t.Fatalf("Failed to get order book: %v", err)
		}
		t.Logf("Order book for %s: %+v", pair.Symbol(), orderBook)
		pair = model.TradingPair{
			ExchangeID: model.BYBIT,
			Base:       currency.BTCSymbol,
			Quote:      currency.USDTSymbol,
			Category:   trade.FUTURES,
		}
		orderBook, err = ps.GetOrderBook(context.Background(), pair, 5)
		if err != nil {
			t.Fatalf("Failed to get order book: %v", err)
		}
		t.Logf("Order book for %s: %+v", pair.Symbol(), orderBook)
		pair = model.TradingPair{
			ExchangeID: model.OKX,
			Base:       currency.ETHSymbol,
			Quote:      currency.USDTSymbol,
			Category:   trade.FUTURES,
		}
		orderBook, err = ps.GetOrderBook(context.Background(), pair, 5)
		if err != nil {
			t.Fatalf("Failed to get order book: %v", err)
		}
		t.Logf("Order book for %s: %+v", pair.Symbol(), orderBook)

	})
}
