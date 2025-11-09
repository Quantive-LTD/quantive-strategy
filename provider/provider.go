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

package provider

import (
	"context"
	"errors"
	"sync"

	"github.com/wang900115/quant/model"
)

var (
	errMissingProvider = errors.New("provider not found")
)

type Provider interface {
	GetPrice(ctx context.Context, pair model.TradingPair) (*model.PricePoint, error)
	GetKlines(ctx context.Context, pair model.TradingPair, interval string, limit int) ([]model.PriceInterval, error)
	GetOrderBook(ctx context.Context, pair model.TradingPair, limit int) (*model.OrderBook, error)
}

type Providers struct {
	mu       sync.RWMutex
	registry map[model.ExchangeId]Provider
}

func New() Providers {
	return Providers{
		registry: make(map[model.ExchangeId]Provider),
	}
}

func (p *Providers) Register(exchangeID model.ExchangeId, provider Provider) {
	p.mu.Lock()
	p.registry[exchangeID] = provider
	p.mu.Unlock()
}

func (p *Providers) Unregister(exchangeID model.ExchangeId) {
	p.mu.Lock()
	delete(p.registry, exchangeID)
	p.mu.Unlock()
}

func (p *Providers) GetPrice(ctx context.Context, pair model.TradingPair) (*model.PricePoint, error) {
	p.mu.RLock()
	provider, ok := p.registry[pair.ExchangeID]
	p.mu.RUnlock()
	if !ok {
		return nil, errMissingProvider
	}
	return provider.GetPrice(ctx, pair)
}

func (p *Providers) GetKlines(ctx context.Context, pair model.TradingPair, interval string, limit int) ([]model.PriceInterval, error) {
	p.mu.RLock()
	provider, ok := p.registry[pair.ExchangeID]
	p.mu.RUnlock()
	if !ok {
		return nil, errMissingProvider
	}
	return provider.GetKlines(ctx, pair, interval, limit)
}

func (p *Providers) GetOrderBook(ctx context.Context, pair model.TradingPair, limit int) (*model.OrderBook, error) {
	p.mu.RLock()
	provider, ok := p.registry[pair.ExchangeID]
	p.mu.RUnlock()
	if !ok {
		return nil, errMissingProvider
	}
	return provider.GetOrderBook(ctx, pair, limit)
}

func (p *Providers) ListProviders() []model.Exchange {
	p.mu.RLock()
	defer p.mu.RUnlock()
	exchanges := make([]model.Exchange, 0, len(p.registry))
	for id := range p.registry {
		ex := model.GetExchange(id)
		exchanges = append(exchanges, ex)
	}
	return exchanges
}
