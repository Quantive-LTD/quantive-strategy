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

package currency

import "github.com/shopspring/decimal"

type CurrencyId byte

const (
	USD CurrencyId = iota
	USDT
	TWD
	BTC
	ETH
)

type CurrencyName string

const (
	USDName  CurrencyName = "US Dollar"
	USDTName CurrencyName = "Tether"
	TWDName  CurrencyName = "New Taiwan Dollar"
	BTCName  CurrencyName = "Bitcoin"
	ETHName  CurrencyName = "Ethereum"
)

type CurrencySymbol string

const (
	USDSymbol  CurrencySymbol = "USD"
	SOLSymbol  CurrencySymbol = "SOL"
	USDTSymbol CurrencySymbol = "USDT"
	TWDSymbol  CurrencySymbol = "TWD"
	BTCSymbol  CurrencySymbol = "BTC"
	ETHSymbol  CurrencySymbol = "ETH"
)

type Currency struct {
	ID       CurrencyId
	Name     CurrencyName
	Symbol   CurrencySymbol
	Decimals int
	IsFiat   bool
}

var CurrencyMap = map[CurrencyId]Currency{
	USD:  {ID: USD, Name: USDName, Symbol: "USD", Decimals: 2, IsFiat: true},
	USDT: {ID: USDT, Name: USDTName, Symbol: "USDT", Decimals: 2, IsFiat: false},
	TWD:  {ID: TWD, Name: TWDName, Symbol: "TWD", Decimals: 2, IsFiat: true},
	BTC:  {ID: BTC, Name: BTCName, Symbol: "BTC", Decimals: 8, IsFiat: false},
	ETH:  {ID: ETH, Name: ETHName, Symbol: "ETH", Decimals: 8, IsFiat: false},
}

func GetCurrency(id CurrencyId) (Currency, bool) {
	c, ok := CurrencyMap[id]
	return c, ok
}

type CurrencyRate struct {
	Rate         decimal.Decimal
	BaseCurrency CurrencyId
}

func (r CurrencyRate) Convert(amount decimal.Decimal) decimal.Decimal {
	return amount.Mul(r.Rate)
}
