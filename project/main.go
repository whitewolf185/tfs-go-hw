package main

import (
	"context"
	"main.go/project/hendlers"
	tg_bot "main.go/project/tg-bot"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)

	orChan := tg_bot.BotStart(ctx, &wg)
	go hendlers.HandStart(ctx, &wg, orChan)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs
	cancel()
	wg.Wait()
}
