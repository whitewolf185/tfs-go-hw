package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/whitewolf185/fs-go-hw/project/DB"
	"github.com/whitewolf185/fs-go-hw/project/hendlers"
	tg_bot "github.com/whitewolf185/fs-go-hw/project/tg-bot"
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
