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

package bybit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/model"
	"github.com/wang900115/quant/model/trade"
)

const (
	END_POINT = "https://api.bybit.com"
	TIMEOUT   = 10 * time.Second
)

var (
	errBybitNoData    = fmt.Errorf("bybit: no data returned")
	errResponseFailed = fmt.Errorf("bybit: response failed")
)

type BybitClient struct {
	endpoint   string
	httpClient *http.Client
}

func NewClient() *BybitClient {
	return &BybitClient{
		endpoint:   END_POINT,
		httpClient: &http.Client{Timeout: TIMEOUT},
	}
}

func getCategoryParam(category trade.Category) string {
	switch category {
	case trade.SPOT:
		return "spot"
	case trade.FUTURES:
		return "linear"
	default:
		return "spot"
	}
}

func (bc *BybitClient) GetPrice(ctx context.Context, pair model.TradingPair) (*model.PricePoint, error) {
	symbol := fmt.Sprintf("%s%s", pair.Base, pair.Quote)
	category := getCategoryParam(pair.Category)
	url := fmt.Sprintf("%s/v5/market/tickers?category=%s&symbol=%s", bc.endpoint, category, symbol)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := bc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errResponseFailed
	}

	var raw struct {
		RetCode int    `json:"retCode"`
		RetMsg  string `json:"retMsg"`
		Result  struct {
			Category string `json:"category"`
			List     []struct {
				Symbol    string `json:"symbol"`
				LastPrice string `json:"lastPrice"`
			} `json:"list"`
		} `json:"result"`
		Time int64 `json:"time"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	if len(raw.Result.List) == 0 {
		return nil, errBybitNoData
	}
	price, err := decimal.NewFromString(raw.Result.List[0].LastPrice)
	if err != nil {
		return nil, err
	}
	return &model.PricePoint{
		NewPrice:  price,
		UpdatedAt: time.Now(),
	}, nil
}

// interval: 1, 3, 5, 15, 30, 60, 120, 240, 360, 720, D, W, M
func (bc *BybitClient) GetKlines(ctx context.Context, pair model.TradingPair, interval string, limit int) ([]model.PriceInterval, error) {
	symbol := fmt.Sprintf("%s%s", pair.Base, pair.Quote)
	category := getCategoryParam(pair.Category)
	url := fmt.Sprintf("%s/v5/market/kline?category=%s&symbol=%s&interval=%s&limit=%d",
		bc.endpoint, category, symbol, interval, limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := bc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bybit api error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var raw struct {
		RetCode int    `json:"retCode"`
		RetMsg  string `json:"retMsg"`
		Result  struct {
			Category string     `json:"category"`
			Symbol   string     `json:"symbol"`
			List     [][]string `json:"list"` // [startTime, open, high, low, close, volume, turnover]
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("failed to decode klines: %w", err)
	}

	intervals := make([]model.PriceInterval, 0, limit)

	for _, kline := range raw.Result.List {
		if len(kline) < 6 {
			return nil, errBybitNoData
		}

		// Bybit response: [startTime, open, high, low, close, volume, turnover]
		// startTime is in milliseconds
		startTimestamp, err := decimal.NewFromString(kline[0])
		if err != nil {
			return nil, err
		}
		openTime := time.UnixMilli(startTimestamp.IntPart())

		openPrice, err := decimal.NewFromString(kline[1])
		if err != nil {
			return nil, err
		}
		highPrice, err := decimal.NewFromString(kline[2])
		if err != nil {
			return nil, err
		}
		lowPrice, err := decimal.NewFromString(kline[3])
		if err != nil {
			return nil, err
		}
		closePrice, err := decimal.NewFromString(kline[4])
		if err != nil {
			return nil, err
		}
		volume, err := decimal.NewFromString(kline[5])
		if err != nil {
			return nil, err
		}

		duration := parseInterval(interval)
		closeTime := openTime.Add(duration)

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

func (bc *BybitClient) GetOrderBook(ctx context.Context, pair model.TradingPair, limit int) (*model.OrderBook, error) {
	symbol := fmt.Sprintf("%s%s", pair.Base, pair.Quote)
	category := getCategoryParam(pair.Category)
	url := fmt.Sprintf("%s/v5/market/orderbook?category=%s&symbol=%s&limit=%d",
		bc.endpoint, category, symbol, limit)
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
		RetCode int `json:"retCode"`
		Result  struct {
			Bids [][]interface{} `json:"b"`
			Asks [][]interface{} `json:"a"`
			Ts   int64           `json:"ts"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	if raw.RetCode != 0 {
		return nil, errResponseFailed
	}

	bids, err := model.ParseOrderEntries[model.OrderBookBid](raw.Result.Bids)
	if err != nil {
		return nil, err
	}
	asks, err := model.ParseOrderEntries[model.OrderBookAsk](raw.Result.Asks)
	if err != nil {
		return nil, err
	}

	return &model.OrderBook{
		Symbol: pair.Symbol(),
		Time:   time.UnixMilli(raw.Result.Ts),
		Bids:   bids,
		Asks:   asks,
	}, nil
}

func parseInterval(interval string) time.Duration {
	switch interval {
	case "1":
		return 1 * time.Minute
	case "3":
		return 3 * time.Minute
	case "5":
		return 5 * time.Minute
	case "15":
		return 15 * time.Minute
	case "30":
		return 30 * time.Minute
	case "60":
		return 1 * time.Hour
	case "120":
		return 2 * time.Hour
	case "240":
		return 4 * time.Hour
	case "360":
		return 6 * time.Hour
	case "720":
		return 12 * time.Hour
	case "D":
		return 24 * time.Hour
	case "W":
		return 7 * 24 * time.Hour
	case "M":
		return 30 * 24 * time.Hour
	default:
		return 1 * time.Minute
	}
}
