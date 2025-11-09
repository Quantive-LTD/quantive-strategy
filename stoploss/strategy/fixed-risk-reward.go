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

package strategy

import (
	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/stoploss"
)

type RiskRewardRatio struct {
	stoploss.BaseResolver
	LastPrice   decimal.Decimal
	riskRatio   decimal.Decimal
	rewardRatio decimal.Decimal
	stopLoss    decimal.Decimal
	takeProfit  decimal.Decimal
}

func NewRiskRewardRatio(entryPrice, riskRatio, rewardRatio decimal.Decimal, callback stoploss.DefaultCallback) (stoploss.HybridWithoutTime, error) {
	if riskRatio.IsNegative() || riskRatio.GreaterThan(decimal.NewFromInt(1)) {
		return nil, errStopLossRateInvalid
	}
	if rewardRatio.IsNegative() || rewardRatio.GreaterThan(decimal.NewFromInt(1)) {
		return nil, errStopLossRateInvalid
	}
	s := &RiskRewardRatio{
		LastPrice:   entryPrice,
		riskRatio:   riskRatio,
		rewardRatio: rewardRatio,
		stopLoss:    entryPrice.Mul(decimal.NewFromInt(1).Sub(riskRatio)),
		takeProfit:  entryPrice.Mul(decimal.NewFromInt(1).Add(rewardRatio)),
		BaseResolver: stoploss.BaseResolver{
			Active:   true,
			Callback: callback,
		},
	}
	return s, nil
}

func (r *RiskRewardRatio) Calculate(currentPrice decimal.Decimal) (decimal.Decimal, decimal.Decimal, error) {
	if !r.Active {
		return decimal.Zero, decimal.Zero, stoploss.ErrStatusInvalid
	}
	r.LastPrice = currentPrice
	r.stopLoss = currentPrice.Mul(decimal.NewFromInt(1).Sub(r.riskRatio))
	r.takeProfit = currentPrice.Mul(decimal.NewFromInt(1).Add(r.rewardRatio))
	return r.stopLoss, r.takeProfit, nil
}

func (r *RiskRewardRatio) ShouldTriggerStopLoss(currentPrice decimal.Decimal) (bool, error) {
	if !r.Active {
		return false, stoploss.ErrStatusInvalid
	}
	if currentPrice.LessThanOrEqual(r.stopLoss) {
		err := r.Trigger(stoploss.TRIGGERED_REASON_HYBRID_RISK_REWARD_STOPLOSS)
		if err != nil {
			return true, stoploss.ErrCallBackFail
		}
		return true, nil
	}
	return false, nil
}

func (r *RiskRewardRatio) ShouldTriggerTakeProfit(currentPrice decimal.Decimal) (bool, error) {
	if !r.Active {
		return false, stoploss.ErrStatusInvalid
	}
	if currentPrice.GreaterThanOrEqual(r.takeProfit) {
		err := r.Trigger(stoploss.TRIGGERED_REASON_HYBRID_RISK_REWARD_TAKEPROFIT)
		if err != nil {
			return true, stoploss.ErrCallBackFail
		}
		return true, nil
	}
	return false, nil
}

func (r *RiskRewardRatio) GetStopLoss() (decimal.Decimal, error) {
	if !r.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return r.stopLoss, nil
}

func (r *RiskRewardRatio) GetTakeProfit() (decimal.Decimal, error) {
	if !r.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return r.takeProfit, nil
}

func (r *RiskRewardRatio) ReSet(newPrice decimal.Decimal) error {
	if !r.Active {
		return stoploss.ErrStatusInvalid
	}
	r.LastPrice = newPrice
	r.stopLoss = newPrice.Mul(decimal.NewFromInt(1).Sub(r.riskRatio))
	r.takeProfit = newPrice.Mul(decimal.NewFromInt(1).Add(r.rewardRatio))
	return nil
}
