package tg_bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type TgBot struct {
	TgAPI  *tgbotapi.BotAPI
	orChan chan Orders

	ctx    context.Context
	cancel context.CancelFunc

	errHandler MyErrors
}

func MakeTgBot(ctx context.Context, orChan chan Orders) TgBot {
	var (
		bot        TgBot
		err        error
		errHandler MyErrors
	)
	token := takeTgBotToken()
	bot.TgAPI, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		errHandler.TgBotErr(err)
	}
	bot.ctx, bot.cancel = context.WithCancel(ctx) // нужно ли?
	bot.orChan = make(chan Orders)
	bot.orChan = orChan

	return bot
}

func (bot *TgBot) SendOrder(orderType OrdersTypes, ticket string, price string) error {
	var order Orders

	switch orderType {
	case BuyOrder:
		order.Endpoint = "/api/v3/sendorder"

		dataP := make(map[string]string)
		dataP["orderType"] = "ioc"
		dataP["symbol"] = ticket
		dataP["side"] = "buy"
		dataP["size"] = "1"
		dataP["limitPrice"] = price

		order.PostData = dataP

		bot.orChan <- order
		log.Println("Order was sent")
	}

	return nil
}

func (bot TgBot) Close() error {
	bot.cancel()

	return nil
}
