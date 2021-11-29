package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	tg_bot "main.go/project/tg-bot"

	"main.go/project/DB"
	"main.go/project/hendlers"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	DBQueChan := DB.StartDB(ctx, &wg)
	orChan, optionChan, TGQueChan, TGTakeChan, TGStopChan := tg_bot.BotStart(ctx, &wg)
	hendlers.HandStart(ctx, &wg, orChan, optionChan, DBQueChan, TGQueChan, TGTakeChan, TGStopChan)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs
	cancel()
	wg.Wait()
}
