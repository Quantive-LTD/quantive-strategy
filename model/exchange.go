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
