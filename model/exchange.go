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

package model

import "github.com/shopspring/decimal"

type ExchangeId byte

const (
	BINANCE ExchangeId = iota
	COINBASE
	OKX
	BYBIT
)

type Exchange struct {
	ID              ExchangeId
	Name            string
	FeeRate         decimal.Decimal
	DefaultCurrency string
}

var ExchangeMap = map[ExchangeId]Exchange{
	BINANCE:  {ID: BINANCE, Name: "Binance", FeeRate: decimal.NewFromFloat(0.001), DefaultCurrency: "USDT"},
	COINBASE: {ID: COINBASE, Name: "Coinbase", FeeRate: decimal.NewFromFloat(0.005), DefaultCurrency: "USD"},
	OKX:      {ID: OKX, Name: "OKX", FeeRate: decimal.NewFromFloat(0.0015), DefaultCurrency: "USDT"},
	BYBIT:    {ID: BYBIT, Name: "Bybit", FeeRate: decimal.NewFromFloat(0.00075), DefaultCurrency: "USDT"},
}

func GetExchange(id ExchangeId) Exchange {
	ex, _ := ExchangeMap[id]
	return ex
}
