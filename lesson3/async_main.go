package main

import (
	"context"
	"hw-async/addition"
	"hw-async/domain"
	"hw-async/generator"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

var tickers = []string{"AAPL", "SBER", "NVDA", "TSLA"}

func getData(prices <-chan domain.Price, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	C1m := make(chan domain.Price)
	C2m := make(chan domain.Price)
	C10m := make(chan domain.Price)

	dataLogger := log.New()
	for {
		select {
		case <-ctx.Done():
			dataLogger.Info("done")
			return
		case v := <-prices:
			C1m <- v
			C2m <- v
			C10m <- v
		}
	}
}

func createCandle(ctx context.Context, price <-chan domain.Price, canPer domain.CandlePeriod) {
	var (
		// TODO думаю можно будет без слайса делать, и после закрытия сразу в файл записывать
		candles []domain.Candle
		err     error
		timer   time.Time
	)
	candleMap := make(map[string]domain.Candle, 0)

	for {
		select {
		case <-ctx.Done():
			// TODO реализация закрытия свечки
		case val := <-price:
			candle, ok := candleMap[val.Ticker]
			if !ok {
				// только создал свечку
				candle = addition.CreateNewCandle(val, canPer)

				candleMap[val.Ticker] = candle
			} else {
				timer, err = domain.PeriodTS(canPer, val.TS)
				if err != nil {
					panic(err)
				}

				if timer.Minute()-candle.TS.Minute() >= 1 {
					// TODO сделать отправку свечи в файл
					candle.Close = val.Value
					candles = append(candles, candle)
					candle = addition.CreateNewCandle(val, canPer)
				} else {
					candle.Low = addition.MinVal(candle.Low, val.Value)
					candle.High = addition.MaxVal(candle.High, val.Value)

				}
				candleMap[val.Ticker] = candle
			}
		}
	}
}

func main() {
	logger := log.New()
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	// signals set up
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	pg := generator.NewPricesGenerator(generator.Config{
		Factor:  10,
		Delay:   time.Millisecond * 500,
		Tickers: tickers,
	})

	logger.Info("start prices generator...")
	prices := pg.Prices(ctx)

	// working area
	wg.Add(1)
	go getData(prices, ctx, &wg)

	// end of working area
	<-sigs
	cancel()

	wg.Wait()
}
