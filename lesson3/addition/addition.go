package addition

import (
	"fmt"
	log "github.com/sirupsen/logrus"
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

func GetCanPerStr(period domain.CandlePeriod) (string, error) {
	switch period {
	case domain.CandlePeriod1m:
		return "1m", nil
	case domain.CandlePeriod2m:
		return "2m", nil
	case domain.CandlePeriod10m:
		return "10m", nil
	default:
		return "", domain.ErrUnknownPeriod
	}
}

func candleToStr(v domain.Candle, canPer domain.CandlePeriod) []string {
	var result []string
	// Ticker
	result = append(result, v.Ticker)
	// Period
	canPerStr, err := GetCanPerStr(canPer)
	if err != nil {
		panic(err)
	}
	result = append(result, canPerStr)
	// Open
	res := fmt.Sprintf("%f", v.Open)
	result = append(result, res)
	// High
	res = fmt.Sprintf("%f", v.High)
	result = append(result, res)
	// Low
	res = fmt.Sprintf("%f", v.Low)
	result = append(result, res)
	// Close
	res = fmt.Sprintf("%f", v.Close)
	result = append(result, res)
	// TS
	result = append(result, v.TS.String())

	return result
}

var Logger = log.New()
