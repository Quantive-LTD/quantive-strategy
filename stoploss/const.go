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

package stoploss

const (
	TRIGGERED_REASON_FIXED_TRAILING_STOPLOSS      = "Trailing Stop Loss Triggered"
	TRIGGERED_REASON_FIXED_TRAILING_TAKEPROFIT    = "Trailing Take Profit Triggered"
	TRIGGERED_REASON_FIXED_PERCENTCILE_STOPLOSS   = "Fixed Percentile Stop Loss Triggered"
	TRIGGERED_REASON_FIXED_PERCENTCILE_TAKEPROFIT = "Fixed Percentile Take Profit Triggered"
	TRIGGERED_REASON_FIXED_ATR_STOPLOSS           = "ATR Based Stop Loss Triggered"
	TRIGGERED_REASON_FIXED_ATR_TAKEPROFIT         = "ATR Based Take Profit Triggered"
	TRIGGERED_REASON_FIXED_MA_STOPLOSS            = "Moving Average Stop Loss Triggered"
	TRIGGERED_REASON_FIXED_MA_TAKEPROFIT          = "Moving Average Take Profit Triggered"
)

const (
	TRIGGERED_REASON_TIMED_TRAILING_STOPLOSS      = "Trailing Stop Loss Triggered with Time Delay"
	TRIGGERED_REASON_TIMED_TRAILING_TAKEPROFIT    = "Trailing Take Profit Triggered with Time Delay"
	TRIGGERED_REASON_TIMED_PERCENTCILE_STOPLOSS   = "Fixed Percentile Stop Loss Triggered with Time Delay"
	TRIGGERED_REASON_TIMED_PERCENTCILE_TAKEPROFIT = "Fixed Percentile Take Profit Triggered with Time Delay"
	TRIGGERED_REASON_TIMED_ATR_STOPLOSS           = "ATR Based Stop Loss Triggered with Time Delay"
	TRIGGERED_REASON_TIMED_ATR_TAKEPROFIT         = "ATR Based Take Profit Triggered with Time Delay"
)

const (
	TRIGGERED_REASON_HYBRID_RISK_REWARD_STOPLOSS   = "Hybrid Risk-Reward Stop Loss Triggered"
	TRIGGERED_REASON_HYBRID_RISK_REWARD_TAKEPROFIT = "Hybrid Risk-Reward Take Profit Triggered"
	TRIGGERED_REASON_STRUCTURE_SWING_STOPLOSS      = "Structure Swing Stop Loss Triggered"
	TRIGGERED_REASON_STRUCTURE_SWING_TAKEPROFIT    = "Structure Swing Take Profit Triggered"
)
