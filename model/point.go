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

import (
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/model/currency"
	"github.com/wang900115/quant/model/trade"
)

type PriceTick struct {
	TradingPair
	PriceInterval
}

type TradingPair struct {
	ExchangeID ExchangeId
	Base       currency.CurrencySymbol
	Quote      currency.CurrencySymbol
	Category   trade.Category
	MinLot     decimal.Decimal
	MinTick    decimal.Decimal
}

func (tp TradingPair) Symbol() string {
	return fmt.Sprintf("%s/%s", tp.Base, tp.Quote)
}

type PriceInterval struct {
	OpenTime         string
	CloseTime        string
	OpeningPrice     decimal.Decimal
	HighestPrice     decimal.Decimal
	LowestPrice      decimal.Decimal
	ClosingPrice     decimal.Decimal
	Volume           decimal.Decimal
	IntervalDuration time.Duration
}

func (pi PriceInterval) String() string {
	return fmt.Sprintf("OpenTime: %s, CloseTime: %s, Duration: %s, Open: %s, Close: %s, High: %s, Low: %s, Volume: %s",
		pi.OpenTime, pi.CloseTime, pi.IntervalDuration.String(),
		pi.OpeningPrice.String(), pi.ClosingPrice.String(),
		pi.HighestPrice.String(), pi.LowestPrice.String(),
		pi.Volume.String())
}

func (pi *PriceInterval) CalculateIntervalDuration() error {
	openTime, err := time.Parse(time.RFC3339, pi.OpenTime)
	if err != nil {
		return err
	}
	closeTime, err := time.Parse(time.RFC3339, pi.CloseTime)
	if err != nil {
		return err
	}
	pi.IntervalDuration = closeTime.Sub(openTime)
	return nil
}

type PricePoint struct {
	NewPrice  decimal.Decimal
	UpdatedAt time.Time
}

var (
	errInValidOrderBookType = errors.New("invalid order book type")
	errNoOrderBookEntries   = errors.New("no order book entries")
)

type OrderBook struct {
	Symbol string
	Time   time.Time
	Bids   []OrderBookBid
	Asks   []OrderBookAsk
}

type OrderBookBase struct {
	Price    decimal.Decimal
	Quantity decimal.Decimal
}

type OrderBookBid OrderBookBase
type OrderBookAsk OrderBookBase

type OrderBookEntry interface {
	OrderBookBid | OrderBookAsk
}

func ParseOrderEntries[T OrderBookEntry](entries [][]interface{}) ([]T, error) {
	result := make([]T, 0, len(entries))
	for _, entry := range entries {
		if len(entry) < 2 {
			return nil, errNoOrderBookEntries
		}
		priceStr, ok1 := entry[0].(string)
		qtyStr, ok2 := entry[1].(string)
		if !ok1 || !ok2 {
			return nil, errInValidOrderBookType
		}
		price, err := decimal.NewFromString(priceStr)
		if err != nil {
			return nil, err
		}
		quantity, err := decimal.NewFromString(qtyStr)
		if err != nil {
			return nil, err
		}
		entry := OrderBookBase{
			Price:    price,
			Quantity: quantity,
		}
		result = append(result, T(entry))
	}
	return result, nil
}
