// Copyright 2025 Perry. All rights reserved.

// Licensed MIT License

// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// https://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package manager

import (
	"context"
	"sync"
	"time"

	"github.com/wang900115/quant/model"
	"github.com/wang900115/quant/model/result"
	"github.com/wang900115/quant/stoploss"
)

type ManagerConfig struct {
	BufferSize    int
	ReadTimeout   time.Duration
	CheckInterval time.Duration
}

type Manager struct {
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	shutdown chan struct{}

	portfolio *Portfolio
	execution *Execution
	Reporter  *Report
	ManagerConfig
}

func New(config ManagerConfig) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		ctx:           ctx,
		cancel:        cancel,
		shutdown:      make(chan struct{}),
		portfolio:     NewPortfolio(),
		execution:     NewExecutionManager(config.BufferSize),
		Reporter:      NewReport(),
		ManagerConfig: config,
	}
}

func (csm *Manager) RegisterStrategy(name string, strategy interface{}) {
	switch s := strategy.(type) {
	case stoploss.FixedStopLoss:
		csm.portfolio.RegistFixedStoplossStrategy(name, s)
	case stoploss.TimeBasedStopLoss:
		csm.portfolio.RegistTimedStoplossStrategy(name, s)
	case stoploss.FixedTakeProfit:
		csm.portfolio.RegistFixedTakeProfitStrategy(name, s)
	case stoploss.TimeBasedTakeProfit:
		csm.portfolio.RegistTimedTakeProfitStrategy(name, s)
	case stoploss.HybridWithoutTime:
		csm.portfolio.RegistHybridStrategy(name, s)
	case stoploss.HybridWithTime:
		csm.portfolio.RegistHybridTimedStrategy(name, s)
	default:
		panic("unsupported strategy type")
	}
}

func (csm *Manager) Start() {
	csm.wg.Add(6)
	go csm.handleFixedStopLoss()
	go csm.handleTimedStopLoss()
	go csm.handleFixedProfit()
	go csm.handleTimedProfit()
	go csm.handleFixedHybrid()
	go csm.handleTimedHybrid()

	generalResult, hybridResult := csm.execution.getResult()
	go csm.Reporter.ProcessGeneralResult(generalResult)
	go csm.Reporter.ProcessHybridResult(hybridResult)
}

func (csm *Manager) Snapshot() map[string]int64 {
	return csm.Reporter.Stats()
}

func (csm *Manager) handleFixedStopLoss() {
	defer csm.wg.Done()
	for {
		select {
		case <-csm.ctx.Done():
			return
		case update := <-csm.execution.fixedStoplossChannel:
			csm.processFixedStopStrategies(update)
		case <-time.After(csm.CheckInterval):
			// Periodic health check
			continue
		}

	}
}

func (csm *Manager) handleTimedStopLoss() {
	defer csm.wg.Done()
	for {
		select {
		case <-csm.ctx.Done():
			return
		case update := <-csm.execution.timedStoplossChannel:
			csm.processTimedStopStrategies(update)
		case <-time.After(csm.CheckInterval):
			// Periodic health check
			continue
		}
	}
}

func (csm *Manager) handleFixedProfit() {
	defer csm.wg.Done()
	for {
		select {
		case <-csm.ctx.Done():
			return
		case update := <-csm.execution.fixedTakeProfitChannel:
			csm.processFixedProfitStrategies(update)
		case <-time.After(csm.CheckInterval):
			// Periodic health check
			continue
		}
	}
}

func (csm *Manager) handleTimedProfit() {
	defer csm.wg.Done()
	for {
		select {
		case <-csm.ctx.Done():
			return
		case update := <-csm.execution.timedTakeProfitChannel:
			csm.processTimedProfitStrategies(update)
		case <-time.After(csm.CheckInterval):
			// Periodic health check
			continue
		}
	}
}

func (csm *Manager) handleFixedHybrid() {
	defer csm.wg.Done()
	for {
		select {
		case <-csm.ctx.Done():
			return
		case update := <-csm.execution.hybridFixedChannel:
			csm.processHybridFixedStrategies(update)
		case <-time.After(csm.CheckInterval):
			// Periodic health check
			continue
		}
	}
}

func (csm *Manager) handleTimedHybrid() {
	defer csm.wg.Done()
	for {
		select {
		case <-csm.ctx.Done():
			return
		case update := <-csm.execution.hybridTimedChannel:
			csm.processHybridTimedStrategies(update)
		case <-time.After(csm.CheckInterval):
			// Periodic health check
			continue
		}
	}
}

func (csm *Manager) processFixedStopStrategies(update model.PricePoint) {
	strategies := csm.portfolio.GetFixedStoplossStrategies()
	for name, strategy := range strategies {
		newThreshold, err := strategy.CalculateStopLoss(update.NewPrice)
		if err == nil {
			result := result.NewGeneral(name, "Fixed", "StopLoss", update.NewPrice, newThreshold, update.UpdatedAt, time.Duration(0))
			shouldTrigger, err := strategy.ShouldTriggerStopLoss(update.NewPrice)
			if err == nil {
				result.SetTriggered(shouldTrigger)
			} else {
				result.SetError(err)
			}
			select {
			case csm.execution.generalResults <- *result:
			case <-csm.ctx.Done():
				return
			}
		}
	}
}

func (csm *Manager) processTimedStopStrategies(update model.PricePoint) {
	strategies := csm.portfolio.GetTimedStoplossStrategies()
	for name, strategy := range strategies {
		timeThreshold, err := strategy.GetTimeThreshold()
		newThreshold, err := strategy.CalculateStopLoss(update.NewPrice)
		if err == nil {
			result := result.NewGeneral(name, "Timed", "StopLoss", update.NewPrice, newThreshold, update.UpdatedAt, time.Duration(timeThreshold))
			shouldTrigger, err := strategy.ShouldTriggerStopLoss(update.NewPrice, update.UpdatedAt.UnixMilli())
			if err == nil && shouldTrigger {
				result.SetTriggered(shouldTrigger)
			} else {
				result.SetError(err)
			}
			select {
			case csm.execution.generalResults <- *result:
			case <-csm.ctx.Done():
				return
			}
		}
	}
}

func (csm *Manager) processFixedProfitStrategies(update model.PricePoint) {
	strategies := csm.portfolio.GetFixedTakeProfitStrategies()
	for name, strategy := range strategies {
		newThreshold, err := strategy.CalculateTakeProfit(update.NewPrice)
		if err == nil {
			result := result.NewGeneral(name, "Fixed", "TakeProfit", update.NewPrice, newThreshold, update.UpdatedAt, time.Duration(0))
			shouldTrigger, err := strategy.ShouldTriggerTakeProfit(update.NewPrice)
			if err == nil {
				result.SetTriggered(shouldTrigger)
			} else {
				result.SetError(err)
			}
			select {
			case csm.execution.generalResults <- *result:
			case <-csm.ctx.Done():
				return
			}
		}
	}
}

func (csm *Manager) processTimedProfitStrategies(update model.PricePoint) {
	strategies := csm.portfolio.GetTimedTakeProfitStrategies()
	for name, strategy := range strategies {
		timeThreshold, err := strategy.GetTimeThreshold()
		newThreshold, err := strategy.CalculateTakeProfit(update.NewPrice)
		if err == nil {
			result := result.NewGeneral(name, "Timed", "TakeProfit", update.NewPrice, newThreshold, update.UpdatedAt, time.Duration(timeThreshold))
			shouldTrigger, err := strategy.ShouldTriggerTakeProfit(update.NewPrice, update.UpdatedAt.UnixMilli())
			if err == nil {
				result.SetTriggered(shouldTrigger)
			} else {
				result.SetError(err)
			}
			select {
			case csm.execution.generalResults <- *result:
			case <-csm.ctx.Done():
				return
			}
		}
	}
}

func (csm *Manager) processHybridFixedStrategies(update model.PricePoint) {
	strategies := csm.portfolio.GetHybridStrategies()
	for name, strategy := range strategies {
		newStop, newProfit, err := strategy.Calculate(update.NewPrice)
		if err == nil {
			result := result.NewHybrid(name, "Fixed", "Hybrid", newStop, newProfit, update.NewPrice, update.UpdatedAt, time.Duration(0))
			shouldTrigger, err := strategy.ShouldTriggerStopLoss(update.NewPrice)
			if err == nil {
				result.SetTriggered(shouldTrigger)
			} else {
				result.SetError(err)
			}
			select {
			case csm.execution.hybridResults <- *result:
			case <-csm.ctx.Done():
				return
			}
		}
	}
}

func (csm *Manager) processHybridTimedStrategies(update model.PricePoint) {
	strategies := csm.portfolio.GetHybridStrategies()
	for name, strategy := range strategies {
		newStop, newProfit, err := strategy.Calculate(update.NewPrice)
		if err == nil {
			result := result.NewHybrid(name, "Timed", "Hybrid", newStop, newProfit, update.NewPrice, update.UpdatedAt, time.Duration(0))
			shouldTrigger, err := strategy.ShouldTriggerStopLoss(update.NewPrice)
			if err == nil {
				result.SetTriggered(shouldTrigger)
			} else {
				result.SetError(err)
			}
			select {
			case csm.execution.hybridResults <- *result:
			case <-csm.ctx.Done():
				return
			}
		}
	}
}

func (csm *Manager) Collect(pricePoint model.PricePoint, callback func()) {
	select {
	case csm.execution.fixedStoplossChannel <- pricePoint:
	case <-csm.ctx.Done():
		return
	default:
		callback()
	}

	select {
	case csm.execution.timedStoplossChannel <- pricePoint:
	case <-csm.ctx.Done():
		return
	default:
		callback()
	}

	select {
	case csm.execution.timedTakeProfitChannel <- pricePoint:
	case <-csm.ctx.Done():
		return
	default:
		callback()
	}

	select {
	case csm.execution.hybridFixedChannel <- pricePoint:
	case <-csm.ctx.Done():
		return
	default:
		callback()
	}

	select {
	case csm.execution.hybridTimedChannel <- pricePoint:
	case <-csm.ctx.Done():
		return
	default:
		callback()
	}
}

func (csm *Manager) Stop() {
	csm.cancel()
	csm.wg.Wait()
	csm.execution.closeChannels()
}

// func (csm *Manager) Stats() map[string]interface{} {
// 	stats := make(map[string]interface{})
// 	stats["FixedStopLoss"] = csm.portfolio.fixedStoplossStrategies
// 	stats["TimedStopLoss"] = csm.portfolio.timedStoplossStrategies
// 	stats["FixedTakeProfit"] = csm.portfolio.fixedTakeProfitStrategies
// 	stats["TimedTakeProfit"] = csm.portfolio.timedTakeProfitStrategies
// 	stats["HybridWithoutTime"] = csm.portfolio.hybridStrategies
// 	stats["HybridWithTime"] = csm.portfolio.hybridTimedStrategies
// 	return stats
// }
