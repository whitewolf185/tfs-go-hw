package tg_bot

import (
	"context"
	"github.com/whitewolf185/fs-go-hw/project/pkg/tg-bot/TG_bot"
	"github.com/whitewolf185/fs-go-hw/project/repository/DB/add_DB"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/whitewolf185/fs-go-hw/project/cmd/addition"
	"github.com/whitewolf185/fs-go-hw/project/cmd/addition/MyErrors"
)

func BotStart(ctx context.Context, wg *sync.WaitGroup) (chan addition.Orders, chan addition.Options, chan add_DB.Query,
	chan addition.TakeProfitCh, chan addition.StopLossCh) {
	orChan := make(chan addition.Orders)
	optionChan := make(chan addition.Options)
	queryChan := make(chan add_DB.Query)
	takeChan := make(chan addition.TakeProfitCh)
	stopChan := make(chan addition.StopLossCh)

	wg.Add(1)
	go func() {
		tgBot := MakeTgBot(ctx, orChan, queryChan)
		tgBot.OrderHandler(wg)
		defer func() {
			tgBot.Close()
			wg.Done()
		}()

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates, err := tgBot.TgAPI.GetUpdatesChan(u)
		if err != nil {
			MyErrors.TgBotUpdateErr(err)
		}

		var option addition.Options
		started := false
		for {
			select {
			case update := <-updates:
				if update.Message == nil {
					continue
				}

				Type, err := TG_bot.MessageType(update.Message.Text)
				if err != nil {
					MyErrors.RegexpErr(err)
					continue
				}
				switch Type {
				case TG_bot.Start:
					tgBot.chatID = update.Message.Chat.ID
					tgBot.SendMessage("Приветики")
					started = true
					continue
				case TG_bot.OptionMsg:
					if !started {
						continue
					}
					if option.CanPer != "" {
						tgBot.SendMessage("Вы уже отправляли настройки")
						continue
					}
					tgBot.SendMessage("Теперь отправьте тикет, к которому вы хотите подключиться")
					ticket := <-updates
					tgBot.SendMessage("Теперь отправьте период свечки")
					canPer := <-updates
					option, err = TG_bot.CreateOptions(ticket.Message.Text, canPer.Message.Text)
					if err != nil {
						tgBot.SendMessage("Вы что-то сделали не так. Посмотрите логи")
						MyErrors.TgBotMsgErr(err)
						continue
					}

					optionChan <- option
				case TG_bot.BuyNow:
					if tgBot.checkOption(option, &started) {
						continue
					}

					tgBot.SendOrder(TG_bot.BuyOrder, option.Ticket[0])

				case TG_bot.SellNow:
					if tgBot.checkOption(option, &started) {
						continue
					}

					tgBot.SendOrder(TG_bot.SellOrder, option.Ticket[0])

				case TG_bot.StopLoss:
					if tgBot.checkOption(option, &started) {
						continue
					}

					// todo нужно запрашивать цену и сколько продавать
				case TG_bot.TakeProfit:
					if tgBot.checkOption(option, &started) {
						continue
					}

					// todo нужно запрашивать цену и сколько продавать
				}

			case <-ctx.Done():
				close(optionChan)
				close(takeChan)
				close(stopChan)
				return
			}
		}
	}()

	return orChan, optionChan, queryChan, takeChan, stopChan
}
