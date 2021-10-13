package addition

import (
	"errors"
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

func showWriterErr(period domain.CandlePeriod) error {
	return errors.New(fmt.Sprintf("writer err in %s", period))
}

func badCloseErr(period domain.CandlePeriod) error {
	return errors.New(fmt.Sprintf("Bad close file %s", period))
}

func badOpenErr(period domain.CandlePeriod) error {
	return errors.New(fmt.Sprintf("Bad open file %s", period))
}

var Logger = log.New()
