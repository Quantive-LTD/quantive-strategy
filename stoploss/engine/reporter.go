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

package engine

import (
	"fmt"

	"github.com/wang900115/quant/model/result"
)

type Report struct {
	generalCount int64
	hybridCount  int64
	triggerCount int64
	errorCount   int64
}

func NewReport() *Report {
	return &Report{}
}

func (rp *Report) ProcessGeneralResult(res <-chan result.StrategyGeneralResult) {
	rp.generalCount++

	for r := range res {
		if r.Error != nil {
			rp.errorCount++
			fmt.Printf("🔴 ERROR in %s (%s): %v\n",
				r.StrategyName, r.StrategyType, r.Error)
			continue
		}

		if r.Triggered {
			rp.triggerCount++
			fmt.Printf("🔔 TRIGGER: %s (%s) - %s at price %s\n",
				r.StrategyName, r.StrategyType, r.TriggerType,
				r.LastPrice.String())
		} else {
			fmt.Printf("📊 UPDATE: %s (%s) - threshold: %s, price: %s\n",
				r.StrategyName, r.StrategyType,
				r.Stat.PriceThreshold.String(), r.LastPrice.String())
		}
	}
}

func (rp *Report) ProcessHybridResult(res <-chan result.StrategyHybridResult) {
	rp.hybridCount++

	for r := range res {
		if r.Error != nil {
			rp.errorCount++
			fmt.Printf("🔴 HYBRID ERROR in %s: %v\n", r.StrategyName, r.Error)
			continue
		}

		if r.Triggered {
			rp.triggerCount++
			fmt.Printf("🔔 HYBRID TRIGGER: %s at price %s stoploss at %s take profit at %s\n",
				r.StrategyName, r.LastPrice.String(), r.StopStat.PriceThreshold.String(), r.ProfitStat.PriceThreshold.String())
		} else {
			fmt.Printf("📊 HYBRID UPDATE: %s at price %s stoploss at %s take profit at %s\n",
				r.StrategyName, r.LastPrice.String(), r.StopStat.PriceThreshold.String(), r.ProfitStat.PriceThreshold.String())
		}
	}
}

func (rp *Report) Stats() map[string]int64 {
	return map[string]int64{
		"general_results": rp.generalCount,
		"hybrid_results":  rp.hybridCount,
		"triggers":        rp.triggerCount,
		"errors":          rp.errorCount,
	}
}
