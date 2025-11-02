package okx

import (
	"context"
	"testing"

	"github.com/wang900115/quant/model"
	"github.com/wang900115/quant/model/currency"
	"github.com/wang900115/quant/model/trade"
)

func TestGetPrice(t *testing.T) {
	o := NewClient()
	t.Run("SPOT", func(t *testing.T) {
		pair := model.TradingPair{
			Base:     currency.BTCSymbol,
			Quote:    currency.USDTSymbol,
			Category: trade.SPOT,
		}
		price, err := o.GetPrice(context.Background(), pair)
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
		price, err := o.GetPrice(context.Background(), pair)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		t.Logf("FUTURES Price: %s", price.NewPrice.String())
	})
}

func TestGetKlines(t *testing.T) {
	o := NewClient()

	t.Run("SPOT", func(t *testing.T) {
		pair := model.TradingPair{
			Base:     currency.BTCSymbol,
			Quote:    currency.USDTSymbol,
			Category: trade.SPOT,
		}
		intervals, err := o.GetKlines(context.Background(), pair, "1h", 10)
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
		intervals, err := o.GetKlines(context.Background(), pair, "15m", 5)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(intervals) == 0 {
			t.Fatal("expected intervals, got empty slice")
		}
	})
}

func TestGetOrderBook(t *testing.T) {
	o := NewClient()
	t.Run("SPOT", func(t *testing.T) {
		pair := model.TradingPair{
			Base:     currency.BTCSymbol,
			Quote:    currency.USDTSymbol,
			Category: trade.SPOT,
		}
		orderBook, err := o.GetOrderBook(context.Background(), pair, 5)
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
		orderBook, err := o.GetOrderBook(context.Background(), pair, 5)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		t.Logf("Order Book: %+v", orderBook)
	})
}
