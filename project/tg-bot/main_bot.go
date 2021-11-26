package tg_bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"main.go/project/MyErrors"
	"sync"
)

func BotStart(ctx context.Context, wg *sync.WaitGroup) chan Orders {
	orChan := make(chan Orders)

	wg.Add(1)
	go func() {
		tgBot := MakeTgBot(ctx, orChan)
		defer func() {
			_ = tgBot.Close()
			wg.Done()
		}()

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates, err := tgBot.TgAPI.GetUpdatesChan(u)
		if err != nil {
			MyErrors.TgBotUpdateErr(err)
		}

		for {
			select {
			case update := <-updates:
				if update.Message == nil {
					continue
				}

				// TODO сделать конфигурацию ws через бота.
				// В ней должны быть:
				// • настройка отслеживаемого тикета
				// • какие свечки слушать

				log.Println(update.Message.Text)
				Ticket := "PI_XBTUSD"
				_ = tgBot.SendOrder(BuyOrder, Ticket, "60000")

			case <-ctx.Done():
				return
			}
		}
	}()

	return orChan
}
