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
	"context"
	"log"

	"github.com/wang900115/quant/model"
	"github.com/wang900115/quant/model/currency"
	"github.com/wang900115/quant/model/trade"
	"github.com/wang900115/quant/provider"
	"github.com/wang900115/quant/provider/binance"
	"github.com/wang900115/quant/provider/bybit"
	"github.com/wang900115/quant/provider/coinbase"
	"github.com/wang900115/quant/provider/okx"
)

func ProviderExample() {
	ps := provider.New()
	ps.Register(model.BINANCE, binance.NewClient())
	ps.Register(model.BYBIT, bybit.NewClient())
	ps.Register(model.COINBASE, coinbase.NewClient())
	ps.Register(model.OKX, okx.NewClient())

	log.Printf("Providers registered: %+v \n", ps.ListProviders())

	tradingPair := model.TradingPair{
		ExchangeID: model.BINANCE,
		Base:       currency.BTCSymbol,
		Quote:      currency.USDTSymbol,
		Category:   trade.SPOT,
	}
	pricePoint, err := ps.GetPrice(context.Background(), tradingPair)
	if err != nil {
		log.Fatalf("Failed to get price: %v \n", err)
	}
	log.Printf("Price for %s: %+v", tradingPair.Symbol(), *pricePoint)

	tradingPair = model.TradingPair{
		ExchangeID: model.COINBASE,
		Base:       currency.BTCSymbol,
		Quote:      currency.USDTSymbol,
		Category:   trade.SPOT,
	}
	pricePoint, err = ps.GetPrice(context.Background(), tradingPair)
	if err != nil {
		log.Fatalf("Failed to get price: %v \n", err)
	}
	log.Printf("Price for %s: %+v \n", tradingPair.Symbol(), *pricePoint)

	klines, err := ps.GetKlines(context.Background(), tradingPair, "1h", 10)
	if err != nil {
		log.Fatalf("Failed to get klines: %v \n", err)
	}
	log.Printf("Klines for %s: %+v \n", tradingPair.Symbol(), klines)

	orderBook, err := ps.GetOrderBook(context.Background(), tradingPair, 5)
	if err != nil {
		log.Fatalf("Failed to get order book: %v \n", err)
	}
	log.Printf("Order book for %s: %+v \n", tradingPair.Symbol(), orderBook)
}
