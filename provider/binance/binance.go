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

package binance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/model"
	"github.com/wang900115/quant/model/trade"
)

const (
	SPOT_END_POINT    = "https://api.binance.com"
	INVERSE_END_POINT = "https://fapi.binance.com"
	TIMEOUT           = 10 * time.Second
)

var (
	errBinanceNoData  = errors.New("binance: no data returned")
	errResponseFailed = errors.New("binance: response failed")
)

type BinanceClient struct {
	spotEndpoint    string
	inverseEndpoint string
	httpClient      *http.Client
}

func NewClient() *BinanceClient {
	return &BinanceClient{
		spotEndpoint:    SPOT_END_POINT,
		inverseEndpoint: INVERSE_END_POINT,
		httpClient:      &http.Client{Timeout: TIMEOUT},
	}
}

func (bc *BinanceClient) GetPrice(ctx context.Context, pair model.TradingPair) (*model.PricePoint, error) {
	symbol := fmt.Sprintf("%s%s", pair.Base, pair.Quote)
	var url string
	switch pair.Category {
	case trade.SPOT:
		url = fmt.Sprintf("%s/api/v3/ticker/price?symbol=%s", bc.spotEndpoint, symbol)
	case trade.FUTURES:
		url = fmt.Sprintf("%s/fapi/v1/ticker/price?symbol=%s", bc.inverseEndpoint, symbol)
	default:
		url = fmt.Sprintf("%s/api/v3/ticker/price?symbol=%s", bc.spotEndpoint, symbol)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := bc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errResponseFailed
	}
	var raw struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	if len(raw.Price) == 0 {
		return nil, errBinanceNoData
	}
	price, err := decimal.NewFromString(raw.Price)
	if err != nil {
		return nil, err
	}
	data := &model.PricePoint{
		NewPrice:  price,
		UpdatedAt: time.Now(),
	}
	return data, nil
}

func (bc *BinanceClient) GetKlines(ctx context.Context, pair model.TradingPair, interval string, limit int) ([]model.PriceInterval, error) {
	symbol := fmt.Sprintf("%s%s", pair.Base, pair.Quote)
	var url string
	switch pair.Category {
	case trade.SPOT:
		url = fmt.Sprintf("%s/api/v3/klines?symbol=%s&interval=%s&limit=%d", bc.spotEndpoint, symbol, interval, limit)
	case trade.FUTURES:
		url = fmt.Sprintf("%s/fapi/v1/klines?symbol=%s&interval=%s&limit=%d", bc.inverseEndpoint, symbol, interval, limit)
	default:
		url = fmt.Sprintf("%s/api/v3/klines?symbol=%s&interval=%s&limit=%d", bc.spotEndpoint, symbol, interval, limit)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := bc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errResponseFailed
	}
	// Binance response struct: [[openTime, open, high, low, close, volume, closeTime, ...], ...]
	var rawKlines [][]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawKlines); err != nil {
		return nil, fmt.Errorf("failed to decode klines: %w", err)
	}
	intervals := make([]model.PriceInterval, 0, limit)
	for _, kline := range rawKlines {
		if len(kline) < 6 {
			return nil, errBinanceNoData
		}
		openTimestamp := int64(kline[0].(float64))
		openTime := time.UnixMilli(openTimestamp)
		openPrice, err := decimal.NewFromString(kline[1].(string))
		if err != nil {
			return nil, err
		}
		highPrice, err := decimal.NewFromString(kline[2].(string))
		if err != nil {
			return nil, err
		}
		lowPrice, err := decimal.NewFromString(kline[3].(string))
		if err != nil {
			return nil, err
		}
		closePrice, err := decimal.NewFromString(kline[4].(string))
		if err != nil {
			return nil, err
		}
		volume, err := decimal.NewFromString(kline[5].(string))
		if err != nil {
			return nil, err
		}
		closeTimestamp := int64(kline[6].(float64))
		closeTime := time.UnixMilli(closeTimestamp)
		duration := time.Duration(closeTimestamp-openTimestamp) * time.Millisecond
		intervals = append(intervals, model.PriceInterval{
			OpenTime:         openTime.Format(time.RFC3339),
			OpeningPrice:     openPrice,
			HighestPrice:     highPrice,
			LowestPrice:      lowPrice,
			ClosingPrice:     closePrice,
			Volume:           volume,
			CloseTime:        closeTime.Format(time.RFC3339),
			IntervalDuration: duration,
		})
	}
	return intervals, nil
}

func (bc *BinanceClient) GetOrderBook(ctx context.Context, pair model.TradingPair, limit int) (*model.OrderBook, error) {
	symbol := fmt.Sprintf("%s%s", pair.Base, pair.Quote)
	var url string
	switch pair.Category {
	case trade.SPOT:
		url = fmt.Sprintf("%s/api/v3/depth?symbol=%s&limit=%d", bc.spotEndpoint, symbol, limit)
	case trade.FUTURES:
		url = fmt.Sprintf("%s/fapi/v1/depth?symbol=%s&limit=%d", bc.inverseEndpoint, symbol, limit)
	default:
		url = fmt.Sprintf("%s/api/v3/depth?symbol=%s&limit=%d", bc.spotEndpoint, symbol, limit)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := bc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errResponseFailed
	}
	var raw struct {
		LastUpdateID int             `json:"lastUpdateId"`
		Bids         [][]interface{} `json:"bids"`
		Asks         [][]interface{} `json:"asks"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	bids, err := model.ParseOrderEntries[model.OrderBookBid](raw.Bids)
	if err != nil {
		return nil, err
	}
	asks, err := model.ParseOrderEntries[model.OrderBookAsk](raw.Asks)
	if err != nil {
		return nil, err
	}
	orderBook := &model.OrderBook{
		Symbol: pair.Symbol(),
		Time:   time.Now(),
		Bids:   bids,
		Asks:   asks,
	}
	return orderBook, nil
}
