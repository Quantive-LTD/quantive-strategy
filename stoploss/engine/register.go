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
	"sync"

	"github.com/wang900115/quant/stoploss"
)

// ! todo copy low level

type Portfolio struct {
	mutex                     sync.Mutex
	fixedStoplossStrategies   map[string]stoploss.FixedStopLoss
	timedStoplossStrategies   map[string]stoploss.TimeBasedStopLoss
	fixedTakeProfitStrategies map[string]stoploss.FixedTakeProfit
	timedTakeProfitStrategies map[string]stoploss.TimeBasedTakeProfit
	hybridFixedStrategies     map[string]stoploss.HybridWithoutTime
	hybridTimedStrategies     map[string]stoploss.HybridWithTime
	openGeneral               bool
	openHybrid                bool
	count                     int
}

func NewPortfolio() *Portfolio {
	return &Portfolio{
		fixedStoplossStrategies:   make(map[string]stoploss.FixedStopLoss),
		timedStoplossStrategies:   make(map[string]stoploss.TimeBasedStopLoss),
		fixedTakeProfitStrategies: make(map[string]stoploss.FixedTakeProfit),
		timedTakeProfitStrategies: make(map[string]stoploss.TimeBasedTakeProfit),
		hybridFixedStrategies:     make(map[string]stoploss.HybridWithoutTime),
		hybridTimedStrategies:     make(map[string]stoploss.HybridWithTime),
		openGeneral:               false,
		openHybrid:                false,
		count:                     0,
	}
}

func (p *Portfolio) RegistFixedStoplossStrategy(name string, strategy stoploss.FixedStopLoss) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.fixedStoplossStrategies[name] = strategy
	p.openGeneral = true
	p.count++
}

func (p *Portfolio) RegistTimedStoplossStrategy(name string, strategy stoploss.TimeBasedStopLoss) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.timedStoplossStrategies[name] = strategy
	p.openGeneral = true
	p.count++
}

func (p *Portfolio) RegistFixedTakeProfitStrategy(name string, strategy stoploss.FixedTakeProfit) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.fixedTakeProfitStrategies[name] = strategy
	p.openGeneral = true
	p.count++
}

func (p *Portfolio) RegistTimedTakeProfitStrategy(name string, strategy stoploss.TimeBasedTakeProfit) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.timedTakeProfitStrategies[name] = strategy
	p.openGeneral = true
	p.count++
}

func (p *Portfolio) RegistHybridStrategy(name string, strategy stoploss.HybridWithoutTime) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.hybridFixedStrategies[name] = strategy
	p.openHybrid = true
	p.count++
}

func (p *Portfolio) RegistHybridTimedStrategy(name string, strategy stoploss.HybridWithTime) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.hybridTimedStrategies[name] = strategy
	p.openHybrid = true
	p.count++
}

func (p *Portfolio) GetFixedStoplossStrategies() map[string]stoploss.FixedStopLoss {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Return a copy to avoid race conditions
	copyMap := make(map[string]stoploss.FixedStopLoss)
	for k, v := range p.fixedStoplossStrategies {
		copyMap[k] = v
	}
	return copyMap
}

func (p *Portfolio) GetTimedStoplossStrategies() map[string]stoploss.TimeBasedStopLoss {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Return a copy
	copyMap := make(map[string]stoploss.TimeBasedStopLoss)
	for k, v := range p.timedStoplossStrategies {
		copyMap[k] = v
	}
	return copyMap
}

func (p *Portfolio) GetFixedTakeProfitStrategies() map[string]stoploss.FixedTakeProfit {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Return a copy to avoid race conditions
	copyMap := make(map[string]stoploss.FixedTakeProfit)
	for k, v := range p.fixedTakeProfitStrategies {
		copyMap[k] = v
	}
	return copyMap
}

func (p *Portfolio) GetTimedTakeProfitStrategies() map[string]stoploss.TimeBasedTakeProfit {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Return a copy to avoid race conditions
	copyMap := make(map[string]stoploss.TimeBasedTakeProfit)
	for k, v := range p.timedTakeProfitStrategies {
		copyMap[k] = v
	}
	return copyMap
}

func (p *Portfolio) GetHybridStrategies() map[string]stoploss.HybridWithoutTime {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Return a copy to avoid race conditions
	copyMap := make(map[string]stoploss.HybridWithoutTime)
	for k, v := range p.hybridFixedStrategies {
		copyMap[k] = v
	}
	return copyMap
}

func (p *Portfolio) GetHybridTimedStrategies() map[string]stoploss.HybridWithTime {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Return a copy to avoid race conditions
	copyMap := make(map[string]stoploss.HybridWithTime)
	for k, v := range p.hybridTimedStrategies {
		copyMap[k] = v
	}
	return copyMap
}
