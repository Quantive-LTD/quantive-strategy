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

package engine

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/wang900115/quant/model"
	"github.com/wang900115/quant/model/result"
	"github.com/wang900115/quant/stoploss"
)

type engineError struct{ msg string }

func (e *engineError) Error() string {
	return e.msg
}

var (
	errNoStrategies = &engineError{"no strategies registered"}
	errNonsupported = &engineError{"unsupported strategy type"}
)

type Config struct {
	BufferSize    int
	ReadTimeout   time.Duration
	CheckInterval time.Duration
}

type Engine struct {
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	shutdown chan struct{}

	fss bool
	tss bool
	fts bool
	tts bool
	hfs bool
	hts bool

	portfolio *Portfolio
	execution *Execution
	Reporter  *Report
	Config    Config
}

func New(config Config) *Engine {
	ctx, cancel := context.WithCancel(context.Background())
	return &Engine{
		ctx:       ctx,
		cancel:    cancel,
		shutdown:  make(chan struct{}),
		portfolio: NewPortfolio(),
		execution: NewExecutionManager(config.BufferSize),
		Reporter:  NewReport(),
		Config:    config,

		fss: false,
		tss: false,
		fts: false,
		tts: false,
		hfs: false,
		hts: false,
	}
}

func (csm *Engine) RegisterStrategy(name string, strategy interface{}) error {
	switch s := strategy.(type) {
	case stoploss.FixedStopLoss:
		csm.portfolio.RegistFixedStoplossStrategy(name, s)
		csm.fss = true
	case stoploss.TimeBasedStopLoss:
		csm.portfolio.RegistTimedStoplossStrategy(name, s)
		csm.tss = true
	case stoploss.FixedTakeProfit:
		csm.portfolio.RegistFixedTakeProfitStrategy(name, s)
		csm.fts = true
	case stoploss.TimeBasedTakeProfit:
		csm.portfolio.RegistTimedTakeProfitStrategy(name, s)
		csm.tts = true
	case stoploss.HybridWithoutTime:
		csm.portfolio.RegistHybridStrategy(name, s)
		csm.hfs = true
	case stoploss.HybridWithTime:
		csm.portfolio.RegistHybridTimedStrategy(name, s)
		csm.hts = true
	default:
		return errNonsupported
	}
	return nil
}

func (csm *Engine) Start() error {
	count := 0

	if csm.fss {
		count++
	}
	if csm.tss {
		count++
	}
	if csm.fts {
		count++
	}
	if csm.tts {
		count++
	}
	if csm.hfs {
		count++
	}
	if csm.hts {
		count++
	}

	if count == 0 {
		return errNoStrategies
	}

	csm.wg.Add(count)
	if len(csm.portfolio.fixedStoplossStrategies) > 0 {
		go csm.handleFixedStopLoss()
	}
	if len(csm.portfolio.timedStoplossStrategies) > 0 {
		go csm.handleTimedStopLoss()
	}
	if len(csm.portfolio.fixedTakeProfitStrategies) > 0 {
		go csm.handleFixedProfit()
	}
	if len(csm.portfolio.timedTakeProfitStrategies) > 0 {
		go csm.handleTimedProfit()
	}
	if len(csm.portfolio.hybridFixedStrategies) > 0 {
		go csm.handleFixedHybrid()
	}
	if len(csm.portfolio.hybridTimedStrategies) > 0 {
		go csm.handleTimedHybrid()
	}
	generalResult, hybridResult := csm.execution.getResult()
	if len(csm.portfolio.fixedStoplossStrategies) > 0 || len(csm.portfolio.timedStoplossStrategies) > 0 ||
		len(csm.portfolio.fixedTakeProfitStrategies) > 0 ||
		len(csm.portfolio.timedTakeProfitStrategies) > 0 {
		go csm.Reporter.ProcessGeneralResult(generalResult)
	}

	if len(csm.portfolio.hybridFixedStrategies) > 0 || len(csm.portfolio.hybridTimedStrategies) > 0 {
		go csm.Reporter.ProcessHybridResult(hybridResult)
	}

	return nil
}

func (csm *Engine) Snapshot() map[string]int64 {
	return csm.Reporter.Stats()
}

func (csm *Engine) handleFixedStopLoss() {
	defer csm.wg.Done()
	log.Println("[handleFixedStopLoss] started")

	ticker := time.NewTicker(csm.Config.CheckInterval)
	defer ticker.Stop()
	for {
		select {
		case <-csm.ctx.Done():
			log.Println("[handleFixedStopLoss] stopped")
			return
		case update := <-csm.execution.fixedStoplossChannel:
			log.Println("[handleFixedStopLoss] received update price:", update.NewPrice)
			csm.processFixedStopStrategies(update)
		case <-ticker.C:
			// Periodic health check
			log.Println("[handleFixedStopLoss] heartbeat")
			continue
		}

	}
}

func (csm *Engine) handleTimedStopLoss() {
	defer csm.wg.Done()
	log.Println("[handleTimedStopLoss] started")

	ticker := time.NewTicker(csm.Config.CheckInterval)
	defer ticker.Stop()
	for {
		select {
		case <-csm.ctx.Done():
			log.Println("[handleTimedStopLoss] stopped")
			return
		case update := <-csm.execution.timedStoplossChannel:
			log.Println("[handleTimedStopLoss] received update price:", update.NewPrice)
			csm.processTimedStopStrategies(update)
		case <-ticker.C:
			// Periodic health check
			log.Println("[handleTimedStopLoss] heartbeat")
			continue
		}
	}
}

func (csm *Engine) handleFixedProfit() {
	defer csm.wg.Done()
	log.Println("[handleFixedProfit] started")

	ticker := time.NewTicker(csm.Config.CheckInterval)
	defer ticker.Stop()
	for {
		select {
		case <-csm.ctx.Done():
			log.Println("[handleFixedProfit] stopped")
			return
		case update := <-csm.execution.fixedTakeProfitChannel:
			log.Println("[handleFixedProfit] received update price:", update.NewPrice)
			csm.processFixedProfitStrategies(update)
		case <-ticker.C:
			// Periodic health check
			log.Println("[handleFixedProfit] heartbeat")
			continue
		}
	}
}

func (csm *Engine) handleTimedProfit() {
	defer csm.wg.Done()
	log.Println("[handleTimedProfit] started")

	ticker := time.NewTicker(csm.Config.CheckInterval)
	defer ticker.Stop()
	for {
		select {
		case <-csm.ctx.Done():
			log.Println("[handleTimedProfit] stopped")
			return
		case update := <-csm.execution.timedTakeProfitChannel:
			log.Println("[handleTimedProfit] received update price:", update.NewPrice)
			csm.processTimedProfitStrategies(update)
		case <-ticker.C:
			// Periodic health check
			log.Println("[handleTimedProfit] heartbeat")
			continue
		}
	}
}

func (csm *Engine) handleFixedHybrid() {
	defer csm.wg.Done()
	log.Println("[handleFixedHybrid] started")

	ticker := time.NewTicker(csm.Config.CheckInterval)
	defer ticker.Stop()
	for {
		select {
		case <-csm.ctx.Done():
			log.Println("[handleFixedHybrid] stopped")
			return
		case update := <-csm.execution.hybridFixedChannel:
			log.Println("[handleFixedHybrid] received update price:", update.NewPrice)
			csm.processHybridFixedStrategies(update)
		case <-ticker.C:
			// Periodic health check
			log.Println("[handleFixedHybrid] heartbeat")
			continue
		}
	}
}

func (csm *Engine) handleTimedHybrid() {
	defer csm.wg.Done()
	log.Println("[handleTimedHybrid] started")

	ticker := time.NewTicker(csm.Config.CheckInterval)
	defer ticker.Stop()
	for {
		select {
		case <-csm.ctx.Done():
			log.Println("[handleTimedHybrid] stopped")
			return
		case update := <-csm.execution.hybridTimedChannel:
			log.Println("[handleTimedHybrid] received update price:", update.NewPrice)
			csm.processHybridTimedStrategies(update)
		case <-ticker.C:
			// Periodic health check
			log.Println("[handleTimedHybrid] heartbeat")
			continue
		}
	}
}

func (csm *Engine) processFixedStopStrategies(update model.PricePoint) {
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

func (csm *Engine) processTimedStopStrategies(update model.PricePoint) {
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

func (csm *Engine) processFixedProfitStrategies(update model.PricePoint) {
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

func (csm *Engine) processTimedProfitStrategies(update model.PricePoint) {
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

func (csm *Engine) processHybridFixedStrategies(update model.PricePoint) {
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

func (csm *Engine) processHybridTimedStrategies(update model.PricePoint) {
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

func (csm *Engine) Collect(pricePoint model.PricePoint, callback func()) {
	channels := []struct {
		ch   chan model.PricePoint
		name string
	}{
		{csm.execution.fixedStoplossChannel, "fixed stoploss"},
		{csm.execution.timedStoplossChannel, "timed stoploss"},
		{csm.execution.fixedTakeProfitChannel, "fixed take profit"},
		{csm.execution.timedTakeProfitChannel, "timed take profit"},
		{csm.execution.hybridFixedChannel, "hybrid fixed"},
		{csm.execution.hybridTimedChannel, "hybrid timed"},
	}

	for _, c := range channels {
		select {
		case c.ch <- pricePoint:
			log.Printf("Sent to %s channel \n", c.name)
		default:
			callback()
		}
	}
}

func (csm *Engine) Stop() {
	csm.cancel()
	csm.execution.closeChannels()
	csm.wg.Wait()
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
