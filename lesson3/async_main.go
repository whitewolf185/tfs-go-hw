package main

import (
	"context"
	"fmt"
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

	dataLogger := log.New()
	for {
		select {
		case <-ctx.Done():
			dataLogger.Info("done")
			return
		case v := <-prices:
			fmt.Println(v)
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

	wg.Add(1)
	go getData(prices, ctx, &wg)
	<-sigs
	cancel()

	wg.Wait()
}
