/*
 The entire transaction model consists of three phases: opening, closing, and reporting.
 Each phase is represented by a distinct struct to encapsulate relevant details.
 This modular approach allows for clear tracking of trades from initiation to completion.
*/

package trade

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/model/currency"
)

type Category string

const (
	SPOT    Category = "SPOT"
	FUTURES Category = "FUTURES"
)

type TradeInfo struct {
	// unique trade identifier of trading system
	ID uuid.UUID
	// trading pair involved in the trade
	Base  currency.Currency
	Quote currency.Currency
}

// TradeOpen represents the opening details of a trade
type TradeOpen struct {
	//reference to TradeInfo.ID
	TradeID uuid.UUID
	// reference quote price to base currency's price
	EntryPrice decimal.Decimal
	// quantity in base currency
	Quantity decimal.Decimal
	// total cost in quote currency
	TotalCost decimal.Decimal
	CreatedAt time.Time
}

// TradeClose represents the closing details of a trade
type TradeClose struct {
	// reference to TradeInfo.ID
	TradeID uuid.UUID
	// quote price to base currency's price
	WithdrawalPrice decimal.Decimal
	// quantity in base currency
	Quantity decimal.Decimal
	// total payout in quote currency
	TotalPayout decimal.Decimal
	CreatedAt   time.Time
}

// TradeReport summarizes the result of a completed trade
type TradeReport struct {
	// reference to TradeInfo.ID
	TradeID uuid.UUID
	// profit or loss in quote currency
	ProfitLoss decimal.Decimal
	// percentage return on the trade in base currency
	ReturnBasePct decimal.Decimal
	// percentage return on the trade in quote currency
	ReturnQuotePct decimal.Decimal
	// duration of the trade
	Duration  time.Duration
	CreatedAt time.Time
}
