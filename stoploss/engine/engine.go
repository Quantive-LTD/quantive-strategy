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
	"runtime/debug"
	"sync"
	"time"

	"github.com/wang900115/quant/model"
	"github.com/wang900115/quant/model/result"
	"github.com/wang900115/quant/stoploss"
)

type strategyEngineError struct{ msg string }

func (e *strategyEngineError) Error() string {
	return e.msg
}

var (
	errNoStrategies = &strategyEngineError{"no strategies registered"}
	errNonsupported = &strategyEngineError{"unsupported strategy type"}
)

type Config struct {
	// Buffer size for internal channels
	BufferSize int
	// Timeout for reading data from channels
	ReadTimeout time.Duration
	// Heartbeat interval for strategy processing goroutines
	CheckInterval time.Duration
	// Heartbeat interval from internal channels
	HeartbeatInterval time.Duration
	// Interval between for failed goroutine retry mechanism
	RetryInterval time.Duration
}

func DefaultConfig() Config {
	return Config{
		BufferSize:        2048,
		ReadTimeout:       3 * time.Second,
		CheckInterval:     5 * time.Second,
		HeartbeatInterval: 5 * time.Second,
		RetryInterval:     1 * time.Second,
	}
}

type StrategyEngine struct {
	engine    *Engine
	portfolio *Portfolio
	execution *Execution
	Reporter  *Report
	Metrics   *Metrics
	Config    Config
}

func New(config Config) *StrategyEngine {
	return &StrategyEngine{
		engine:    NewEngine(config.RetryInterval, config.CheckInterval),
		portfolio: NewPortfolio(),
		execution: NewExecutionManager(config.BufferSize),
		Reporter:  NewReport(),
		Metrics:   NewMetrics(),
		Config:    config,
	}
}

func (csm *StrategyEngine) RegisterStrategy(name string, strategy interface{}) error {
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
		return errNonsupported
	}
	return nil
}

func (csm *StrategyEngine) Start() error {

	if csm.portfolio.count == 0 {
		return errNoStrategies
	}
	goroutineCount := 0
	if len(csm.portfolio.fixedStoplossStrategies) > 0 {
		goroutineCount++
		csm.engine.SafeGo(func(ctx context.Context) {
			csm.handleFixedStopLoss(ctx)
		}, nil)
	}
	if len(csm.portfolio.timedStoplossStrategies) > 0 {
		goroutineCount++
		csm.engine.SafeGo(func(ctx context.Context) {
			csm.handleTimedStopLoss(ctx)
		}, nil)
	}
	if len(csm.portfolio.fixedTakeProfitStrategies) > 0 {
		goroutineCount++
		csm.engine.SafeGo(func(ctx context.Context) {
			csm.handleFixedProfit(ctx)
		}, nil)
	}
	if len(csm.portfolio.timedTakeProfitStrategies) > 0 {
		goroutineCount++
		csm.engine.SafeGo(func(ctx context.Context) {
			csm.handleTimedProfit(ctx)
		}, nil)
	}
	if len(csm.portfolio.hybridFixedStrategies) > 0 {
		goroutineCount++
		csm.engine.SafeGo(func(ctx context.Context) {
			csm.handleFixedHybrid(ctx)
		}, nil)
	}
	if len(csm.portfolio.hybridTimedStrategies) > 0 {
		goroutineCount++
		csm.engine.SafeGo(func(ctx context.Context) {
			csm.handleTimedHybrid(ctx)
		}, nil)
	}

	generalResult, hybridResult := csm.execution.getResult()
	if csm.portfolio.openGeneral {
		csm.engine.SafeGo(func(ctx context.Context) {
			csm.Reporter.ProcessGeneralResult(generalResult, ctx)
		}, nil)
	}

	if csm.portfolio.openHybrid {
		csm.engine.SafeGo(func(ctx context.Context) {
			csm.Reporter.ProcessHybridResult(hybridResult, ctx)
		}, nil)
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

func (csm *StrategyEngine) Snapshot() map[string]int64 {
	return csm.Reporter.Stats()
}

func (csm *StrategyEngine) GetMetrics() map[string]interface{} {
	return csm.Metrics.Stats()
}

func (csm *StrategyEngine) handleFixedStopLoss(ctx context.Context) {
	ticker := time.NewTicker(csm.Config.HeartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("[handleFixedStopLoss] stopped")
			return
		case update := <-csm.execution.fixedStoplossChannel:
			csm.processFixedStopStrategies(update, ctx)
		case <-ticker.C:
			log.Println("[handleFixedStopLoss] heartbeat")
			continue
		}

	}
}

func (csm *StrategyEngine) handleTimedStopLoss(ctx context.Context) {
	ticker := time.NewTicker(csm.Config.HeartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("[handleTimedStopLoss] stopped")
			return
		case update := <-csm.execution.timedStoplossChannel:
			csm.processTimedStopStrategies(update, ctx)
		case <-ticker.C:
			log.Println("[handleTimedStopLoss] heartbeat")
			continue
		}
	}
}

func (csm *StrategyEngine) handleFixedProfit(ctx context.Context) {
	ticker := time.NewTicker(csm.Config.HeartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("[handleFixedProfit] stopped")
			return
		case update := <-csm.execution.fixedTakeProfitChannel:
			csm.processFixedProfitStrategies(update, ctx)
		case <-ticker.C:
			log.Println("[handleFixedProfit] heartbeat")
			continue
		}
	}
}

func (csm *StrategyEngine) handleTimedProfit(ctx context.Context) {
	ticker := time.NewTicker(csm.Config.HeartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("[handleTimedProfit] stopped")
			return
		case update := <-csm.execution.timedTakeProfitChannel:
			csm.processTimedProfitStrategies(update, ctx)
		case <-ticker.C:
			log.Println("[handleTimedProfit] heartbeat")
			continue
		}
	}
}

func (csm *StrategyEngine) handleFixedHybrid(ctx context.Context) {
	ticker := time.NewTicker(csm.Config.HeartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("[handleFixedHybrid] stopped")
			return
		case update := <-csm.execution.hybridFixedChannel:
			csm.processHybridFixedStrategies(update, ctx)
		case <-ticker.C:
			log.Println("[handleFixedHybrid] heartbeat")
			continue
		}
	}
}

func (csm *StrategyEngine) handleTimedHybrid(ctx context.Context) {
	ticker := time.NewTicker(csm.Config.HeartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("[handleTimedHybrid] stopped")
			return
		case update := <-csm.execution.hybridTimedChannel:
			csm.processHybridTimedStrategies(update, ctx)
		case <-ticker.C:
			log.Println("[handleTimedHybrid] heartbeat")
			continue
		}
	}
}

func (csm *StrategyEngine) processFixedStopStrategies(update model.PricePoint, ctx context.Context) {
	strategies := csm.portfolio.GetFixedStoplossStrategies()
	for name, strategy := range strategies {
		shouldTrigger, err := strategy.ShouldTriggerStopLoss(update.NewPrice)
		newThreshold, calcErr := strategy.CalculateStopLoss(update.NewPrice)
		if calcErr == nil {
			result := result.NewGeneral(name, model.FIXED, model.STOP_LOSS, update.NewPrice, newThreshold, update.UpdatedAt, time.Duration(0))
			if err == nil {
				result.SetTriggered(shouldTrigger)
			} else {
				result.SetError(err)
			}
			select {
			case csm.execution.generalResults <- *result:
				// successfully sent
			case <-time.After(csm.Config.ReadTimeout):
				csm.Metrics.RecordChannelTimeout(model.FIXED, model.STOP_LOSS)
			case <-ctx.Done():
				return
			}
		}
	}
}

func (csm *StrategyEngine) processTimedStopStrategies(update model.PricePoint, ctx context.Context) {
	strategies := csm.portfolio.GetTimedStoplossStrategies()
	for name, strategy := range strategies {
		timeThreshold, _ := strategy.GetTimeThreshold()
		shouldTrigger, err := strategy.ShouldTriggerStopLoss(update.NewPrice, update.UpdatedAt.UnixMilli())
		newThreshold, calcErr := strategy.CalculateStopLoss(update.NewPrice)
		if calcErr == nil {
			result := result.NewGeneral(name, model.TIMED, model.STOP_LOSS, update.NewPrice, newThreshold, update.UpdatedAt, time.Duration(timeThreshold))
			if err == nil {
				result.SetTriggered(shouldTrigger)
			} else {
				result.SetError(err)
			}
			select {
			case csm.execution.generalResults <- *result:
				// successfully sent
			case <-time.After(csm.Config.ReadTimeout):
				csm.Metrics.RecordChannelTimeout(model.TIMED, model.STOP_LOSS)
			case <-ctx.Done():
				return
			}
		}
	}
}

func (csm *StrategyEngine) processFixedProfitStrategies(update model.PricePoint, ctx context.Context) {
	strategies := csm.portfolio.GetFixedTakeProfitStrategies()
	for name, strategy := range strategies {
		shouldTrigger, err := strategy.ShouldTriggerTakeProfit(update.NewPrice)
		newThreshold, calcErr := strategy.CalculateTakeProfit(update.NewPrice)
		if calcErr == nil {
			result := result.NewGeneral(name, model.FIXED, model.TAKE_PROFIT, update.NewPrice, newThreshold, update.UpdatedAt, time.Duration(0))
			if err == nil {
				result.SetTriggered(shouldTrigger)
			} else {
				result.SetError(err)
			}
			select {
			case csm.execution.generalResults <- *result:
				// successfully sent
			case <-time.After(csm.Config.ReadTimeout):
				csm.Metrics.RecordChannelTimeout(model.FIXED, model.TAKE_PROFIT)
			case <-ctx.Done():
				return
			}
		}
	}
}

func (csm *StrategyEngine) processTimedProfitStrategies(update model.PricePoint, ctx context.Context) {
	strategies := csm.portfolio.GetTimedTakeProfitStrategies()
	for name, strategy := range strategies {
		timeThreshold, _ := strategy.GetTimeThreshold()
		shouldTrigger, err := strategy.ShouldTriggerTakeProfit(update.NewPrice, update.UpdatedAt.UnixMilli())
		newThreshold, calcErr := strategy.CalculateTakeProfit(update.NewPrice)
		if calcErr == nil {
			result := result.NewGeneral(name, model.TIMED, model.TAKE_PROFIT, update.NewPrice, newThreshold, update.UpdatedAt, time.Duration(timeThreshold))
			if err == nil {
				result.SetTriggered(shouldTrigger)
			} else {
				result.SetError(err)
			}
			select {
			case csm.execution.generalResults <- *result:
				// successfully sent
			case <-time.After(csm.Config.ReadTimeout):
				csm.Metrics.RecordChannelTimeout(model.TIMED, model.TAKE_PROFIT)
			case <-ctx.Done():
				return
			}
		}
	}
}

func (csm *StrategyEngine) processHybridFixedStrategies(update model.PricePoint, ctx context.Context) {
	strategies := csm.portfolio.GetHybridStrategies()
	for name, strategy := range strategies {
		shouldTriggerSL, errSL := strategy.ShouldTriggerStopLoss(update.NewPrice)
		shouldTriggerTP, errTP := strategy.ShouldTriggerTakeProfit(update.NewPrice)
		newStop, newProfit, calcErr := strategy.Calculate(update.NewPrice)
		if calcErr == nil {
			result := result.NewHybrid(name, model.HYBRID_FIXED, update.NewPrice, newStop, newProfit, update.UpdatedAt, time.Duration(0))
			if errSL == nil && shouldTriggerSL {
				result.SetTriggered(true, model.STOP_LOSS)
			} else if errTP == nil && shouldTriggerTP {
				result.SetTriggered(true, model.TAKE_PROFIT)
			} else if errSL != nil {
				result.SetError(errSL)
			} else if errTP != nil {
				result.SetError(errTP)
			}

			select {
			case csm.execution.hybridResults <- *result:
				// successfully sent
			case <-time.After(csm.Config.ReadTimeout):
				csm.Metrics.RecordChannelTimeout(model.HYBRID_FIXED, "")
			case <-ctx.Done():
				return
			}
		}
	}
}

func (csm *StrategyEngine) processHybridTimedStrategies(update model.PricePoint, ctx context.Context) {
	strategies := csm.portfolio.GetHybridStrategies()
	for name, strategy := range strategies {
		shouldTriggerSL, errSL := strategy.ShouldTriggerStopLoss(update.NewPrice)
		shouldTriggerTP, errTP := strategy.ShouldTriggerTakeProfit(update.NewPrice)
		newStop, newProfit, calcErr := strategy.Calculate(update.NewPrice)
		if calcErr == nil {
			result := result.NewHybrid(name, model.HYBRID_TIMED, update.NewPrice, newStop, newProfit, update.UpdatedAt, time.Duration(0))
			if errSL == nil && shouldTriggerSL {
				result.SetTriggered(true, model.STOP_LOSS)
			} else if errTP == nil && shouldTriggerTP {
				result.SetTriggered(true, model.TAKE_PROFIT)
			} else if errSL != nil {
				result.SetError(errSL)
			} else if errTP != nil {
				result.SetError(errTP)
			}

			select {
			case csm.execution.hybridResults <- *result:
				// successfully sent
			case <-time.After(csm.Config.ReadTimeout):
				csm.Metrics.RecordChannelTimeout(model.HYBRID_TIMED, "")
			case <-ctx.Done():
				return
			}
		}
	}
}

func (csm *StrategyEngine) Collect(pricePoint model.PricePoint, callback func()) {
	csm.Metrics.RecordReceived()

	if len(csm.portfolio.fixedStoplossStrategies) > 0 {
		dataFeedWithMetrics(pricePoint, csm.execution.fixedStoplossChannel, model.FIXED, model.STOP_LOSS, csm.Metrics, callback)
	}
	if len(csm.portfolio.timedStoplossStrategies) > 0 {
		dataFeedWithMetrics(pricePoint, csm.execution.timedStoplossChannel, model.TIMED, model.STOP_LOSS, csm.Metrics, callback)
	}
	if len(csm.portfolio.fixedTakeProfitStrategies) > 0 {
		dataFeedWithMetrics(pricePoint, csm.execution.fixedTakeProfitChannel, model.FIXED, model.TAKE_PROFIT, csm.Metrics, callback)
	}
	if len(csm.portfolio.timedTakeProfitStrategies) > 0 {
		dataFeedWithMetrics(pricePoint, csm.execution.timedTakeProfitChannel, model.TIMED, model.TAKE_PROFIT, csm.Metrics, callback)
	}
	if len(csm.portfolio.hybridFixedStrategies) > 0 {
		dataFeedWithMetrics(pricePoint, csm.execution.hybridFixedChannel, model.HYBRID_FIXED, "", csm.Metrics, callback)
	}
	if len(csm.portfolio.hybridTimedStrategies) > 0 {
		dataFeedWithMetrics(pricePoint, csm.execution.hybridTimedChannel, model.HYBRID_TIMED, "", csm.Metrics, callback)
	}
}

func (csm *StrategyEngine) Stop() {
	csm.engine.Stop()
	csm.execution.closeChannels()
}

type Engine struct {
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	HealthCheck   time.Duration
	RetryInterval time.Duration
}

func NewEngine(retry time.Duration, health time.Duration) *Engine {
	ctx, cancel := context.WithCancel(context.Background())
	return &Engine{
		ctx:           ctx,
		cancel:        cancel,
		RetryInterval: retry,
		HealthCheck:   health,
	}
}

func (e *Engine) Go(fn func(ctx context.Context), recoverFunc func(r any)) {
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()

		defer func() {
			if r := recover(); r != nil {
				if recoverFunc != nil {
					recoverFunc(r)
				} else {
					log.Printf("Recovered from panic in engine goroutine: %v", r)
				}
			}
		}()
		fn(e.ctx)
	}()
}

func (e *Engine) SafeGo(fn func(ctx context.Context), restartFunc func()) {
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()

		for {
			done := make(chan struct{})
			panicCh := make(chan any, 1)

			go func() {
				defer func() {
					if r := recover(); r != nil {
						select {
						case panicCh <- r:
						default:
						}
					}
				}()
				fn(e.ctx)
				close(done)
			}()

			select {
			case <-done:
				return
			case r := <-panicCh:
				log.Printf("[Engine.SafeGo] goroutine panic recovered: %v\n%s", r, debug.Stack())
				if restartFunc != nil {
					restartFunc()
				}
			case <-time.After(e.HealthCheck):
				log.Println("[Engine.SafeGo] goroutine health check timeout, restarting...")
				if restartFunc != nil {
					restartFunc()
				}
			case <-e.ctx.Done():
				return
			}

			select {
			case <-e.ctx.Done():
				return
			case <-time.After(e.RetryInterval):
			}
		}
	}()
}

func (e *Engine) Stop() {
	e.cancel()
	e.wg.Wait()
}

func dataFeedWithMetrics(pricePoint model.PricePoint, channel chan model.PricePoint, typ model.StrategyType, category model.StrategyCategory, metrics *Metrics, callback func()) {
	select {
	case channel <- pricePoint:
		metrics.RecordChannelSend(typ, category)
	default:
		// Channel is full
		metrics.RecordChannelDrop(typ, category)
		metrics.RecordDropped()
		callback()
	}
}
