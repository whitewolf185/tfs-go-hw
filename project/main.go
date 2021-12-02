package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/whitewolf185/fs-go-hw/project/pkg/hendlers"
	"github.com/whitewolf185/fs-go-hw/project/pkg/tg-bot"
	"github.com/whitewolf185/fs-go-hw/project/repository/DB"
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
