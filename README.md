
![Quant Logo](./assets/assets.png)

## Overview

**Quant** is a quantitative analysis platform designed for financial data tracking, strategy modeling, and performance evaluation. It provides modular tools for algorithmic trading, data analytics, and strategy backtesting.

## Features

- 📊 **Data Integration** – Connect to multiple exchanges and data providers.  
- ⚙️ **Strategy Engine** – Design and test your own trading strategies.  
- 🧠 **Analytics Module** – Visualize market trends, price indicators, and risk metrics.  
- 💾 **Persistence Layer** – Store and query market data with efficient caching.  
- 🧩 **Modular Design** – Clean architecture suitable for research and production.

## Tech Stack

- **Backend:** Go / Python  
- **Database:** MySQL / Redis  
- **Visualization:** React + TypeScript  
- **Deployment:** Docker + Kubernetes  

## Getting Started

```bash
# Clone the repository
git clone https://github.com/yourusername/quant.git
cd quant

# Install dependencies
go mod tidy

# Run the service
go run main.go
```

## External Source 
Quant integrates multiple external market data sources for real-time and historical analysis:

- Binance
 – Spot, Futures, and Options trading data.

- Coinbase
 – Cryptocurrency spot market data and account info.

- OKX
 – Spot, Futures, and Perpetual contracts.

- Bybit
 – Inverse and USDT-margined derivatives data.


## Architecture 
```kotlin
┌─────────────────────┐
│   Data Collector    │  ← Market data feeds
└────────┬────────────┘
         ↓
┌─────────────────────┐
│  Strategy Engine    │  ← Backtesting / Live trading
└────────┬────────────┘
         ↓
┌─────────────────────┐
│  Analytics Module   │  ← Reports & visualization
└─────────────────────┘
```
