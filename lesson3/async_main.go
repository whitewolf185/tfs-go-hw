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

	C1m, C2m, C10m := addition.GetData(ctx, &wg, prices)

	fileC1m := addition.CreateCandle(&wg, C1m, domain.CandlePeriod1m)
	fileC2m := addition.CreateCandle(&wg, C2m, domain.CandlePeriod2m)
	fileC10m := addition.CreateCandle(&wg, C10m, domain.CandlePeriod10m)

	wg.Add(3)
	go addition.WriteToFile(&wg, fileC1m, domain.CandlePeriod1m)
	go addition.WriteToFile(&wg, fileC2m, domain.CandlePeriod2m)
	go addition.WriteToFile(&wg, fileC10m, domain.CandlePeriod10m)

	// end of working area
	<-sigs
	cancel()

	wg.Wait()
}
