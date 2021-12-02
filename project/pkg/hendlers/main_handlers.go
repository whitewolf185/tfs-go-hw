package hendlers

import (
	"context"
	"encoding/json"
	add_Conn2 "github.com/whitewolf185/fs-go-hw/project/pkg/hendlers/addConn"
	"github.com/whitewolf185/fs-go-hw/project/pkg/tg-bot/addTGbot"
	"github.com/whitewolf185/fs-go-hw/project/repository/DB/addDB"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/whitewolf185/fs-go-hw/project/cmd/addition"
	"github.com/whitewolf185/fs-go-hw/project/cmd/addition/MyErrors"
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

		log.Info("Waiting fot incoming options")
		option, ok := <-optionChan
		if !ok {
			return
		}
		log.Info("Handler caught options")
		canChan, err := api.GetCandles(wg, optionChan)
		if err != nil {
			MyErrors.GetCandlesErr()
		}

		api.OrderListener(wg)

		var stop addition.StopLossCh
		var take addition.TakeProfitCh
		var prevCan add_Conn2.EventMsg // возможно нужен будет, чтобы провернуть схему, которая нужна, чтобы убрать ошибочное срабатывание
		for {
			select {
			case <-ctx.Done():
				return

			case can := <-canChan:
				if prevCan.Candle.Time == 0 {
					prevCan = can
					continue
				}
				canOpen := add_Conn2.ConvertToFloat(can.Candle.Open)
				canClose := add_Conn2.ConvertToFloat(can.Candle.Close)

				if (canClose+canOpen)/2 <= stop.StopFl {
					api.SendOrder(addTGbot.SellOrder, option.Ticket[0], stop.Size)
				}

				if (canClose+canOpen)/2 >= take.TakeFl {
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
