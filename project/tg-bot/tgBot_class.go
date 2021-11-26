package tg_bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"main.go/project/MyErrors"
	"main.go/project/addition"
)

type TgBot struct {
	TgAPI  *tgbotapi.BotAPI
	orChan chan Orders

	ctx    context.Context
	cancel context.CancelFunc
}

func MakeTgBot(ctx context.Context, orChan chan Orders) TgBot {
	var (
		bot TgBot
		err error
	)
	token := addition.TakeTgBotToken()
	bot.TgAPI, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		MyErrors.TgBotErr(err)
	}
	bot.ctx, bot.cancel = context.WithCancel(ctx)
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

func (bot *TgBot) Close() error {
	bot.cancel()

	return nil
}
