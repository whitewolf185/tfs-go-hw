package addition

import (
	"hw-async/domain"
)

func CreateNewCandle(val domain.Price, canPer domain.CandlePeriod) domain.Candle {
	var (
		candle domain.Candle
		err    error
	)

	candle.Low = val.Value
	candle.High = val.Value
	candle.Ticker = val.Ticker
	candle.TS, err = domain.PeriodTS(domain.CandlePeriod1m, val.TS)
	if err != nil {
		panic(err)
	}
	candle.Open = val.Value
	candle.Close = val.Value
	candle.Period = canPer

	return candle
}

func MinVal(lhs, rhs float64) float64 {
	if lhs < rhs {
		return lhs
	}
	return rhs
}

func MaxVal(lhs, rhs float64) float64 {
	if lhs > rhs {
		return lhs
	}
	return rhs
}
