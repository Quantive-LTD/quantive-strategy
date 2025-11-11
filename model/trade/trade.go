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
	INVERSE Category = "INVERSE"
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
