package hendlers

import (
	"context"
	"encoding/json"
	"sync"

	log "github.com/sirupsen/logrus"

	"main.go/project/addition"
	"main.go/project/addition/MyErrors"
	"main.go/project/addition/TG_bot"
	"main.go/project/addition/add_Conn"
	"main.go/project/addition/add_DB"
)

func HandStart(ctx context.Context, wg *sync.WaitGroup,
	orChan chan addition.Orders, optionChan chan addition.Options, DBQueChan chan add_DB.Query, TGQueChan chan add_DB.Query,
	TGTakeChan chan addition.TakeProfitCh, TGStopChan chan addition.StopLossCh) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		api := MakeAPI(Connection{}, ctx, orChan, DBQueChan, TGQueChan)
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
		log.Info(jsonData)

		log.Info("Waiting fot incoming options")
		option, ok := <-optionChan
		if !ok {
			return
		}
		log.Info("Handler caught options")
		canChan := api.connServ.GetCandles(api.Ws, wg, api.Ctx, option)

		api.OrderListener(wg)

		var stop addition.StopLossCh
		var take addition.TakeProfitCh
		var prevCan add_Conn.EventMsg // возможно нужен будет, чтобы провернуть схему, которая нужна, чтобы убрать ошибочное срабатывание
		for {
			select {
			case <-ctx.Done():
				return

			case can := <-canChan:
				if prevCan.Candle.Time == 0 {
					prevCan = can
					continue
				}
				canOpen := add_Conn.ConvertToFloat(can.Candle.Open)
				canClose := add_Conn.ConvertToFloat(can.Candle.Close)

				if (canClose+canOpen)/2 <= stop.StopFl {
					api.SendOrder(TG_bot.SellOrder, option.Ticket[0], stop.Size)
				}

				if (canClose+canOpen)/2 >= take.TakeFl {
					api.SendOrder(TG_bot.BuyOrder, option.Ticket[0], take.Size)
				}

			case stop = <-TGStopChan:
				log.Info("New StopLoss is", stop)

			case take = <-TGTakeChan:
				log.Info("New TakeProfit is", take)
			}

		}
	}()
}
