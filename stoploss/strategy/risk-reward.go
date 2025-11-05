package strategy

import (
	"github.com/shopspring/decimal"
	"github.com/wang900115/quant/stoploss"
)

type riskRewardRatio struct {
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
	s := &riskRewardRatio{
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

func (r *riskRewardRatio) CalculateStopLoss(currentPrice decimal.Decimal) (decimal.Decimal, error) {
	if !r.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return r.stopLoss, nil
}

func (r *riskRewardRatio) CalculateTakeProfit(currentPrice decimal.Decimal) (decimal.Decimal, error) {
	if !r.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return r.takeProfit, nil
}

func (r *riskRewardRatio) ShouldTriggerStopLoss(currentPrice decimal.Decimal) (bool, error) {
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

func (r *riskRewardRatio) ShouldTriggerTakeProfit(currentPrice decimal.Decimal) (bool, error) {
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

func (r *riskRewardRatio) GetStopLoss() (decimal.Decimal, error) {
	if !r.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return r.stopLoss, nil
}

func (r *riskRewardRatio) GetTakeProfit() (decimal.Decimal, error) {
	if !r.Active {
		return decimal.Zero, stoploss.ErrStatusInvalid
	}
	return r.takeProfit, nil
}

func (r *riskRewardRatio) ReSet(newPrice decimal.Decimal) error {
	if !r.Active {
		return stoploss.ErrStatusInvalid
	}
	r.LastPrice = newPrice
	r.stopLoss = newPrice.Mul(decimal.NewFromInt(1).Sub(r.riskRatio))
	r.takeProfit = newPrice.Mul(decimal.NewFromInt(1).Add(r.rewardRatio))
	return nil
}
