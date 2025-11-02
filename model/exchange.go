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
