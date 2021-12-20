package tgBot

import (
	"context"
	"strconv"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/whitewolf185/fs-go-hw/project/cmd/addition"
	"github.com/whitewolf185/fs-go-hw/project/cmd/addition/MyErrors"
	"github.com/whitewolf185/fs-go-hw/project/pkg/tg-bot/addTGbot"
	"github.com/whitewolf185/fs-go-hw/project/repository/DB/addDB"
)

func BotStart(ctx context.Context, wg *sync.WaitGroup) (chan addition.Orders, chan addition.Options, chan addDB.Query,
	chan addition.TakeProfitCh, chan addition.StopLossCh) {
	orChan := make(chan addition.Orders)
	optionChan := make(chan addition.Options)
	queryChan := make(chan addDB.Query)
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

				Type, err := addTGbot.MessageType(update.Message.Text)
				if err != nil {
					MyErrors.RegexpErr(err)
					continue
				}
				switch Type {
				case addTGbot.Start:
					tgBot.chatID = update.Message.Chat.ID
					tgBot.SendMessage("Приветики")
					started = true
					continue
				case addTGbot.OptionMsg:
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
					option, err = addTGbot.CreateOptions(ticket.Message.Text, canPer.Message.Text)
					if err != nil {
						tgBot.SendMessage("Вы что-то сделали не так. Посмотрите логи")
						MyErrors.TgBotMsgErr(err)
						continue
					}

					optionChan <- option
				case addTGbot.BuyNow:
					if tgBot.checkOption(option, &started) {
						continue
					}

					tgBot.SendOrder(addTGbot.BuyOrder, option.Ticket[0])

				case addTGbot.SellNow:
					if tgBot.checkOption(option, &started) {
						continue
					}

					tgBot.SendOrder(addTGbot.SellOrder, option.Ticket[0])

				case addTGbot.StopLoss:
					if tgBot.checkOption(option, &started) {
						continue
					}
					tgBot.SendMessage("Напишите цену срабатывания индикатора")
					msg := <-updates
					cost, err := addition.ConvertToFloat(msg.Message.Text)
					if err != nil {
						tgBot.SendMessage("Вы неправильно ввели значение. Введите команду /stoploss повторно и повторите попытку")
						continue
					}

					tgBot.SendMessage("Напишите количество продаж")
					msg = <-updates
					size, err := strconv.Atoi(msg.Message.Text)
					if err != nil {
						tgBot.SendMessage("Вы неправильно ввели значение. Введите команду /stoploss повторно и повторите попытку")
						continue
					}

					stop := addition.StopLossCh{
						StopFl: cost,
						Size:   size,
					}

					stopChan <- stop

				case addTGbot.TakeProfit:
					if tgBot.checkOption(option, &started) {
						continue
					}

					tgBot.SendMessage("Напишите цену срабатывания индикатора")
					msg := <-updates
					cost, err := addition.ConvertToFloat(msg.Message.Text)
					if err != nil {
						tgBot.SendMessage("Вы неправильно ввели значение. Введите команду /takeprofit повторно и повторите попытку")
						continue
					}

					tgBot.SendMessage("Напишите количество продаж")
					msg = <-updates
					size, err := strconv.Atoi(msg.Message.Text)
					if err != nil {
						tgBot.SendMessage("Вы неправильно ввели значение. Введите команду /takeprofit повторно и повторите попытку")
						continue
					}

					take := addition.TakeProfitCh{
						TakeFl: cost,
						Size:   size,
					}

					takeChan <- take
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
