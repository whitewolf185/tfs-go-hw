package addition

import (
	"context"
	"fmt"
	"hw-async/domain"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func GetData(ctx context.Context, prices <-chan domain.Price, wg *sync.WaitGroup) {
	defer wg.Done()

	C1m := make(chan domain.Price)
	C2m := make(chan domain.Price)
	C10m := make(chan domain.Price)

	wg.Add(3)
	go createCandle(ctx, wg, C1m, domain.CandlePeriod1m)
	go createCandle(ctx, wg, C2m, domain.CandlePeriod2m)
	go createCandle(ctx, wg, C10m, domain.CandlePeriod10m)

	dataLogger := log.New()
	for {
		select {
		case <-ctx.Done():
			dataLogger.Info("done")
			return
		case v := <-prices:
			dataLogger.Info(v)
			C1m <- v
			C2m <- v
			C10m <- v
		}
	}
}

func createCandle(ctx context.Context, wg *sync.WaitGroup,
	price <-chan domain.Price, canPer domain.CandlePeriod) {
	defer wg.Done()
	dataLogger := log.New()
	var (
		// TODO думаю можно будет без слайса делать, и после закрытия сразу в файл записывать
		candles []domain.Candle
		err     error
		timer   time.Time
	)
	candleMap := make(map[string]domain.Candle)
	start := false

	for {
		select {
		case <-ctx.Done():
			// TODO реализация закрытия свечки
			if start {
				for _, candle := range candleMap {
					candles = append(candles, candle)
				}
			}
			fmt.Println(canPer, " ", len(candles)/4)
			return
		case val := <-price:
			candle, ok := candleMap[val.Ticker]
			start = true
			if !ok {
				// только создал свечку
				candle = CreateNewCandle(val, canPer)

				candleMap[val.Ticker] = candle
			} else {
				timer, err = domain.PeriodTS(canPer, val.TS)
				if err != nil {
					panic(err)
				}
				comparison, errInt := domain.PeriodTSInt(canPer)
				if errInt != nil {
					panic(errInt)
				}

				if timer.Minute()-candle.TS.Minute() >= comparison { // закрытие свечи
					// TODO сделать отправку свечи в файл
					candles = append(candles, candle)
					dataLogger.Info("created new candle ", canPer)
					candle = CreateNewCandle(val, canPer)
				} else {
					candle.Low = MinVal(candle.Low, val.Value)
					candle.High = MaxVal(candle.High, val.Value)
				}

				candleMap[val.Ticker] = candle
			}
		}
	}
}
