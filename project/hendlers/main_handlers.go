package hendlers

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"main.go/project/MyErrors"
	"main.go/project/addition"
	"sync"
)

func HandStart(ctx context.Context, wg *sync.WaitGroup,
	orChan chan addition.Orders, optionChan chan addition.Options, DBQueChan chan addition.Query, TGQueChan chan addition.Query,
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
		option := <-optionChan
		log.Info("Handler caught options")
		canChan := api.connServ.GetCandles(api.Ws, wg, api.Ctx, option)

		api.OrderListener(wg)

		var stop addition.StopLossCh
		var take addition.TakeProfitCh
		var prevCan EventMsg // возможно нужен будет, чтобы провернуть схему, которая нужна, чтобы убрать ошибочное срабатывание
		for {
			select {
			case <-ctx.Done():
				return

			case can := <-canChan:
				if prevCan.Candle.Time == 0 {
					prevCan = can
					continue
				}
				canOpen := addition.ConvertToFloat(can.Candle.Open)
				canClose := addition.ConvertToFloat(can.Candle.Close)

				if (canClose+canOpen)/2 <= stop.StopFl {
					api.SendOrder(stop.Size)
				}

				if (canClose+canOpen)/2 >= take.TakeFl {
					api.SendOrder(take.Size)
				}

			case stop = <-TGStopChan:
				log.Info("New StopLoss is", stop)

			case take = <-TGTakeChan:
				log.Info("New TakeProfit is", take)
			}

		}
	}()
}
