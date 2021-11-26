package hendlers

import (
	"context"
	"encoding/json"
	"fmt"
	"main.go/project/MyErrors"
	"main.go/project/addition"
	tg_bot "main.go/project/tg-bot"
	"sync"
)

func HandStart(ctx context.Context, wg *sync.WaitGroup, orChan chan tg_bot.Orders) {
	defer wg.Done()
	api := MakeAPI(Connection{}, ctx, orChan)
	api.WebsocketConnect(wg)
	defer func() {
		if err := api.Close(); err != nil {
			MyErrors.BadApiClose(err)
		}
	}()

	_, data, err := api.Ws.ReadMessage()
	if err != nil {
		MyErrors.WSReadMsgErr(err)
	}

	jsonData := make(map[string]interface{})

	_ = json.Unmarshal(data, &jsonData)

	fmt.Println(jsonData)

	// TODO сюда нужно написать функцию, которая бы забирала настройки из тг бота
	option, err := addition.CreateOptions([]string{"PI_XBTUSD"}, "1m")
	if err != nil {
		MyErrors.UnknownPeriod(err)
		// TODO не совершаю коннект, а боту отправляю информацию, что нужно бы период переписать
	}

	canChan := api.connServ.GetCandles(api.Ws, wg, api.Ctx, option)

	api.OrderListener(wg)
	for {
		select {
		case <-ctx.Done():
			return

		case can := <-canChan:
			fmt.Println(can)
		}
	}
}
