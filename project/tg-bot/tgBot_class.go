package tg_bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"main.go/project/addition"
	"main.go/project/addition/MyErrors"
	"main.go/project/addition/TG_bot"
	"main.go/project/addition/add_DB"
	"sync"
)

type TgBot struct {
	TgAPI  *tgbotapi.BotAPI
	chatID int64

	orChan  chan addition.Orders
	queChan chan add_DB.Query

	ctx    context.Context
	cancel context.CancelFunc
}

func MakeTgBot(ctx context.Context, orChan chan addition.Orders, queChan chan add_DB.Query) *TgBot {
	var (
		bot TgBot
		err error
	)
	token := TG_bot.TakeTgBotToken()
	bot.TgAPI, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		MyErrors.TgBotErr(err)
	}
	bot.ctx, bot.cancel = context.WithCancel(ctx)
	bot.orChan = orChan
	bot.queChan = queChan

	return &bot
}

func (bot *TgBot) SendOrder(orderType TG_bot.OrdersTypes, ticket string) {
	var order addition.Orders

	order.Endpoint = "/api/v3/sendorder"

	dataP := make(map[string]string)
	dataP["orderType"] = "mkt"
	dataP["symbol"] = ticket
	dataP["size"] = "1"

	switch orderType {
	case TG_bot.BuyOrder:
		dataP["side"] = "buy"

	case TG_bot.SellOrder:
		dataP["side"] = "sell"
	}

	order.PostData = dataP

	bot.orChan <- order
	log.Println("Order was sent")
}

func (bot *TgBot) SendMessage(message string) {
	msg := tgbotapi.NewMessage(bot.chatID, message)
	_, err := bot.TgAPI.Send(msg)
	if err != nil {
		MyErrors.SendMsgErr(err)
	}
}

func (bot *TgBot) SendMessageID(message string, ID int64) {
	msg := tgbotapi.NewMessage(ID, message)
	_, err := bot.TgAPI.Send(msg)
	if err != nil {
		MyErrors.SendMsgErr(err)
	}
}

func (bot *TgBot) OrderHandler(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-bot.ctx.Done():
				return

			case query := <-bot.queChan:
				msg := fmt.Sprintf("%s %f %d", query.Type, query.LimitPrice, query.Size)
				bot.SendMessage(msg)
			}
		}
	}()
}

func (bot *TgBot) checkOption(option addition.Options, started *bool) bool {
	if option.CanPer == "" {
		bot.SendMessage("Вы забыли ввести настройки. Введите команду /option")
		return true
	}
	if !*started {
		return true
	}
	return false
}

func (bot *TgBot) Close() {
	close(bot.queChan)
	close(bot.orChan)
	bot.cancel()
}
