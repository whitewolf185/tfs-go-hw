package hendlers

import (
	"context"
	"encoding/json"
	"math"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/whitewolf185/fs-go-hw/project/cmd/addition"
	"github.com/whitewolf185/fs-go-hw/project/cmd/addition/MyErrors"
	"github.com/whitewolf185/fs-go-hw/project/pkg/tg-bot/addTGbot"
	"github.com/whitewolf185/fs-go-hw/project/repository/DB/addDB"
)

func HandStart(ctx context.Context, wg *sync.WaitGroup,
	orChan chan addition.Orders, optionChan chan addition.Options, dBQueChan chan addDB.Query, tGQueChan chan addDB.Query,
	tGTakeChan chan addition.TakeProfitCh, tGStopChan chan addition.StopLossCh) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		api := MakeAPI(ctx, Connection{}, orChan, dBQueChan, tGQueChan)
		err := api.WebsocketConnect(wg)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := api.Close(); err != nil {
				MyErrors.BadAPIClose(err)
			}
		}()

		_, data, err := api.Ws.ReadMessage()
		if err != nil {
			_ = MyErrors.WSReadMsgErr(err)
		}
		jsonData := make(map[string]interface{})
		_ = json.Unmarshal(data, &jsonData)
		log.Info(jsonData)

		canChan, option, err := api.GetCandles(wg, optionChan)
		if err != nil {
			MyErrors.GetCandlesErr()
		}

		api.OrderListener(wg)

		var stop addition.StopLossCh
		take := addition.TakeProfitCh{
			TakeFl: math.MaxFloat32,
		}
		for {
			select {
			case <-ctx.Done():
				return

			case can := <-canChan:
				canOpen, err := addition.ConvertToFloat(can.Candle.Open)
				if err != nil {
					log.Error(err)
					continue
				}
				canClose, err := addition.ConvertToFloat(can.Candle.Close)
				if err != nil {
					log.Error(err)
					continue
				}

				if (canClose+canOpen)/2 <= stop.StopFl {
					stop.StopFl = 0
					log.Info("Indicator stoploss has gone off")
					api.SendOrder(addTGbot.SellOrder, option.Ticket[0], stop.Size)
				}

				if (canClose+canOpen)/2 >= take.TakeFl {
					take.TakeFl = math.MaxFloat32
					log.Info("Indicator takeprofit has gone off")
					api.SendOrder(addTGbot.BuyOrder, option.Ticket[0], take.Size)
				}

			case stop = <-tGStopChan:
				log.Info("New StopLoss is", stop)

			case take = <-tGTakeChan:
				log.Info("New TakeProfit is", take)
			}
		}
	}()
}
