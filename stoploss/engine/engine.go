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

	if csm.portfolio.count == 0 {
		return errNoStrategies
	}
	goroutineCount := 0
	if len(csm.portfolio.fixedStoplossStrategies) > 0 {
		goroutineCount++
		csm.wg.Add(1)
		go csm.handleFixedStopLoss()
	}
	if len(csm.portfolio.timedStoplossStrategies) > 0 {
		goroutineCount++
		csm.wg.Add(1)
		go csm.handleTimedStopLoss()
	}
	if len(csm.portfolio.fixedTakeProfitStrategies) > 0 {
		goroutineCount++
		csm.wg.Add(1)
		go csm.handleFixedProfit()
	}
	if len(csm.portfolio.timedTakeProfitStrategies) > 0 {
		goroutineCount++
		csm.wg.Add(1)
		go csm.handleTimedProfit()
	}
	if len(csm.portfolio.hybridFixedStrategies) > 0 {
		goroutineCount++
		csm.wg.Add(1)
		go csm.handleFixedHybrid()
	}
	if len(csm.portfolio.hybridTimedStrategies) > 0 {
		goroutineCount++
		csm.wg.Add(1)
		go csm.handleTimedHybrid()
	}

	generalResult, hybridResult := csm.execution.getResult()
	if csm.portfolio.openGeneral {
		csm.wg.Add(1)
		go csm.Reporter.ProcessGeneralResult(generalResult, csm.ctx, &csm.wg)
	}

	if csm.portfolio.openHybrid {
		csm.wg.Add(1)
		go csm.Reporter.ProcessHybridResult(hybridResult, csm.ctx, &csm.wg)
	}

	// Log the number of started goroutines
	log.Printf("Started %d strategy goroutines + %d reporter goroutines", goroutineCount, func() int {
		count := 0
		if csm.portfolio.openGeneral {
			count++
		}
		if csm.portfolio.openHybrid {
			count++
		}
		return count
	}())

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
			csm.processFixedStopStrategies(update)
		case <-ticker.C:
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
			csm.processTimedStopStrategies(update)
		case <-ticker.C:
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
			csm.processFixedProfitStrategies(update)
		case <-ticker.C:
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
			csm.processTimedProfitStrategies(update)
		case <-ticker.C:
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
			csm.processHybridFixedStrategies(update)
		case <-ticker.C:
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
			csm.processHybridTimedStrategies(update)
		case <-ticker.C:
			log.Println("[handleTimedHybrid] heartbeat")
			continue
		}
	}
}

func (csm *Engine) processFixedStopStrategies(update model.PricePoint) {
	strategies := csm.portfolio.GetFixedStoplossStrategies()
	for name, strategy := range strategies {
		shouldTrigger, err := strategy.ShouldTriggerStopLoss(update.NewPrice)
		newThreshold, calcErr := strategy.CalculateStopLoss(update.NewPrice)
		if calcErr == nil {
			result := result.NewGeneral(name, "Fixed", "StopLoss", update.NewPrice, newThreshold, update.UpdatedAt, time.Duration(0))
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
		timeThreshold, _ := strategy.GetTimeThreshold()
		shouldTrigger, err := strategy.ShouldTriggerStopLoss(update.NewPrice, update.UpdatedAt.UnixMilli())
		newThreshold, calcErr := strategy.CalculateStopLoss(update.NewPrice)
		if calcErr == nil {
			result := result.NewGeneral(name, "Timed", "StopLoss", update.NewPrice, newThreshold, update.UpdatedAt, time.Duration(timeThreshold))
			if err == nil && shouldTrigger {
				result.SetTriggered(shouldTrigger)
			} else if err != nil {
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
		shouldTrigger, err := strategy.ShouldTriggerTakeProfit(update.NewPrice)
		newThreshold, calcErr := strategy.CalculateTakeProfit(update.NewPrice)
		if calcErr == nil {
			result := result.NewGeneral(name, "Fixed", "TakeProfit", update.NewPrice, newThreshold, update.UpdatedAt, time.Duration(0))
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
		timeThreshold, _ := strategy.GetTimeThreshold()
		shouldTrigger, err := strategy.ShouldTriggerTakeProfit(update.NewPrice, update.UpdatedAt.UnixMilli())
		newThreshold, calcErr := strategy.CalculateTakeProfit(update.NewPrice)
		if calcErr == nil {
			result := result.NewGeneral(name, "Timed", "TakeProfit", update.NewPrice, newThreshold, update.UpdatedAt, time.Duration(timeThreshold))
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
		shouldTriggerSL, errSL := strategy.ShouldTriggerStopLoss(update.NewPrice)
		shouldTriggerTP, errTP := strategy.ShouldTriggerTakeProfit(update.NewPrice)
		newStop, newProfit, calcErr := strategy.Calculate(update.NewPrice)
		if calcErr == nil {
			result := result.NewHybrid(name, "Fixed", update.NewPrice, newStop, newProfit, update.UpdatedAt, time.Duration(0))
			if errSL == nil && shouldTriggerSL {
				result.SetTriggered(true, "StopLoss")
			} else if errTP == nil && shouldTriggerTP {
				result.SetTriggered(true, "TakeProfit")
			} else if errSL != nil {
				result.SetError(errSL)
			} else if errTP != nil {
				result.SetError(errTP)
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
		shouldTriggerSL, errSL := strategy.ShouldTriggerStopLoss(update.NewPrice)
		shouldTriggerTP, errTP := strategy.ShouldTriggerTakeProfit(update.NewPrice)
		newStop, newProfit, calcErr := strategy.Calculate(update.NewPrice)
		if calcErr == nil {
			result := result.NewHybrid(name, "Timed", update.NewPrice, newStop, newProfit, update.UpdatedAt, time.Duration(0))
			if errSL == nil && shouldTriggerSL {
				result.SetTriggered(true, "StopLoss")
			} else if errTP == nil && shouldTriggerTP {
				result.SetTriggered(true, "TakeProfit")
			} else if errSL != nil {
				result.SetError(errSL)
			} else if errTP != nil {
				result.SetError(errTP)
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
	if csm.fss {
		dataFeed(pricePoint, csm.execution.fixedStoplossChannel, callback)
	}
	if csm.tss {
		dataFeed(pricePoint, csm.execution.timedStoplossChannel, callback)
	}
	if csm.fts {
		dataFeed(pricePoint, csm.execution.fixedTakeProfitChannel, callback)
	}
	if csm.tts {
		dataFeed(pricePoint, csm.execution.timedTakeProfitChannel, callback)
	}
	if csm.hfs {
		dataFeed(pricePoint, csm.execution.hybridFixedChannel, callback)
	}
	if csm.hts {
		dataFeed(pricePoint, csm.execution.hybridTimedChannel, callback)
	}
}

func dataFeed(pricePoint model.PricePoint, channel chan model.PricePoint, callback func()) {
	select {
	case channel <- pricePoint:
		// Successfully sent to channel
	default:
		// Channel is full, handle accordingly
		callback()
	}
}

func (csm *Engine) Stop() {
	// 1. Cancel context to signal all goroutines to stop
	csm.cancel()
	// 2. Wait for all goroutines to finish gracefully
	csm.wg.Wait()
	// 3. Close channels after goroutines have stopped
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
