package model

type StrategyType string

const (
	FIXED        StrategyType = "Fixed"
	TIMED        StrategyType = "Timed"
	HYBRID_TIMED StrategyType = "Hybrid-Timed"
	HYBRID_FIXED StrategyType = "Hybrid-Fixed"
)

func (st StrategyType) String() string {
	return string(st)
}

type StrategyCategory string

const (
	STOP_LOSS   StrategyCategory = "stop_loss"
	TAKE_PROFIT StrategyCategory = "take_profit"
)

func (sc StrategyCategory) String() string {
	return string(sc)
}

type StrategySide string

const (
	LONG  StrategySide = "long"
	SHORT StrategySide = "short"
)

func (ss StrategySide) String() string {
	return string(ss)
}

type StrategySignal string

const (
	BUY  StrategySignal = "buy"
	SELL StrategySignal = "sell"
)

func (ss StrategySignal) String() string {
	return string(ss)
}
