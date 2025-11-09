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
	"github.com/wang900115/quant/model"
	"github.com/wang900115/quant/model/result"
)

type Execution struct {
	fixedStoplossChannel   chan model.PricePoint
	timedStoplossChannel   chan model.PricePoint
	fixedTakeProfitChannel chan model.PricePoint
	timedTakeProfitChannel chan model.PricePoint
	hybridFixedChannel     chan model.PricePoint
	hybridTimedChannel     chan model.PricePoint

	generalResults chan result.StrategyGeneralResult
	hybridResults  chan result.StrategyHybridResult
}

func NewExecutionManager(bufferSize int) *Execution {
	return &Execution{
		fixedStoplossChannel:   make(chan model.PricePoint, bufferSize),
		timedStoplossChannel:   make(chan model.PricePoint, bufferSize),
		fixedTakeProfitChannel: make(chan model.PricePoint, bufferSize),
		timedTakeProfitChannel: make(chan model.PricePoint, bufferSize),
		hybridFixedChannel:     make(chan model.PricePoint, bufferSize),
		hybridTimedChannel:     make(chan model.PricePoint, bufferSize),
		generalResults:         make(chan result.StrategyGeneralResult, bufferSize),
		hybridResults:          make(chan result.StrategyHybridResult, bufferSize),
	}
}

func (e *Execution) getResult() (<-chan result.StrategyGeneralResult, <-chan result.StrategyHybridResult) {
	return e.generalResults, e.hybridResults
}

func (e *Execution) closeChannels() {
	close(e.fixedStoplossChannel)
	close(e.timedStoplossChannel)
	close(e.fixedTakeProfitChannel)
	close(e.timedTakeProfitChannel)
	close(e.hybridFixedChannel)
	close(e.hybridTimedChannel)
	close(e.generalResults)
	close(e.hybridResults)
}
