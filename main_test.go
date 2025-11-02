package main

import (
	"context"
	"testing"

	"github.com/wang900115/quant/model"
	"github.com/wang900115/quant/model/currency"
	"github.com/wang900115/quant/model/trade"
	"github.com/wang900115/quant/provider"
	"github.com/wang900115/quant/provider/binance"
	"github.com/wang900115/quant/provider/bybit"
	"github.com/wang900115/quant/provider/coinbase"
	"github.com/wang900115/quant/provider/okx"
)

func TestMain(t *testing.T) {
	// Initial Provider
	ps := provider.New()
	ps.Register(model.BINANCE, binance.NewClient())
	ps.Register(model.BYBIT, bybit.NewClient())
	ps.Register(model.COINBASE, coinbase.NewClient())
	ps.Register(model.OKX, okx.NewClient())

	t.Logf("Providers registered: %+v", ps.ListProviders())

	tradingPair := model.TradingPair{
		ExchangeID: model.BINANCE,
		Base:       currency.BTCSymbol,
		Quote:      currency.USDTSymbol,
		Category:   trade.SPOT,
	}
	pricePoint, err := ps.GetPrice(context.Background(), tradingPair)
	if err != nil {
		t.Fatalf("Failed to get price: %v", err)
	}
	t.Logf("Price for %s: %+v", tradingPair.Symbol(), *pricePoint)

	tradingPair = model.TradingPair{
		ExchangeID: model.COINBASE,
		Base:       currency.BTCSymbol,
		Quote:      currency.USDTSymbol,
		Category:   trade.SPOT,
	}
	pricePoint, err = ps.GetPrice(context.Background(), tradingPair)
	if err != nil {
		t.Fatalf("Failed to get price: %v", err)
	}
	t.Logf("Price for %s: %+v \n", tradingPair.Symbol(), *pricePoint)

	klines, err := ps.GetKlines(context.Background(), tradingPair, "1h", 10)
	if err != nil {
		t.Fatalf("Failed to get klines: %v", err)
	}
	t.Logf("Klines for %s: %+v \n", tradingPair.Symbol(), klines)

	orderBook, err := ps.GetOrderBook(context.Background(), tradingPair, 5)
	if err != nil {
		t.Fatalf("Failed to get order book: %v", err)
	}
	t.Logf("Order book for %s: %+v \n", tradingPair.Symbol(), orderBook)
}
