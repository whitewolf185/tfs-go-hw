package main

import (
	"context"
	"hw-async/addition"
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
	wg.Add(1)
	go addition.GetData(ctx, prices, &wg)

	// end of working area
	<-sigs
	cancel()

	wg.Wait()
}
