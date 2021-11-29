package hendlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/whitewolf185/fs-go-hw/project/addition"
	"github.com/whitewolf185/fs-go-hw/project/addition/MyErrors"
	"github.com/whitewolf185/fs-go-hw/project/addition/add_Conn"
)

const (
	writeWait  = 2 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type SubEvent struct {
	Event string `json:"event"`
}

type WSMsg struct {
	Event   string   `json:"event"`
	Feed    string   `json:"feed"`
	Tickets []string `json:"product_ids"`
}

type Connection struct {
	SubMessage WSMsg
	Candle     add_Conn.EventMsg

	ws *websocket.Conn
}

// candleStream -- функция-обработчик отправки Orders
func (obj *Connection) candleStream(wg *sync.WaitGroup, ctx context.Context) (chan add_Conn.EventMsg, error) {
	canChan := make(chan add_Conn.EventMsg)

	wg.Add(1)
	go func() {
		prev := make(map[string]int64)
		for _, ticket := range obj.SubMessage.Tickets {
			prev[ticket] = 0
		}

		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				close(canChan)
				return
			default:
				var event add_Conn.EventMsg
				err := obj.ws.ReadJSON(&event)
				if err != nil {
					_ = MyErrors.WSReadMsgErr(errors.New("In cansleStream" + err.Error()))
					close(canChan)
					return
				}
				if event.Result == "error" { // обработчик ошибки отправки order
					log.Errorln(event)
				} else if prev[event.ProductId] != event.Candle.Time && event.Candle.Volume != 0 {
					canChan <- event
					prev[event.ProductId] = event.Candle.Time
				}
			}
		}
	}()

	return canChan, nil
}

func (obj Connection) prepareCandles(ws *websocket.Conn, wg *sync.WaitGroup, ctx context.Context, options addition.Options) (chan add_Conn.EventMsg, error) {
	obj.SubMessage.Event = "subscribe"
	obj.SubMessage.Feed = "candles_trade_" + string(options.CanPer)
	obj.SubMessage.Tickets = options.Ticket
	obj.ws = ws

	msg, err := json.Marshal(obj.SubMessage)
	if err != nil {
		MyErrors.MarshalErr(err)
		return nil, err
	}

	fmt.Println(string(msg))

	err = ws.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		MyErrors.WSWriteMsgErr(err)
		return nil, err
	}

	_, data, err := ws.ReadMessage()
	if err != nil {
		return nil, MyErrors.WSReadMsgErr(err)
	}

	jsonData := SubEvent{}
	if err = json.Unmarshal(data, &jsonData); err != nil {
		return nil, MyErrors.UnmarshalErr(err)
	}
	fmt.Println(jsonData)
	if jsonData.Event != "subscribed" {
		return nil, MyErrors.SubErr
	}

	return obj.candleStream(wg, ctx)
}

func (obj Connection) GetCandles(ws *websocket.Conn, wg *sync.WaitGroup, ctx context.Context, optionChan chan addition.Options) (chan add_Conn.EventMsg, error) {
	log.Info("Waiting fot incoming options")
	option, ok := <-optionChan
	if !ok {
		return nil, MyErrors.OptionChanErr
	}
	log.Info("Handler caught options")

	canChan, err := obj.prepareCandles(ws, wg, ctx, option)
	for i := 0; err != nil && i < 10; i++ {
		log.Errorln(err)
		log.Info("Some err was caught. Waiting fot another incoming options... Try", i)
		option, ok = <-optionChan
		if !ok {
			return nil, MyErrors.OptionChanErr
		}
		log.Info("Handler caught options")
		canChan, err = obj.prepareCandles(ws, wg, ctx, option)
	}

	return canChan, err
}

func (obj Connection) Unsubscribe(ws *websocket.Conn) error {
	obj.SubMessage.Event = "unsubscribe"
	err := ws.WriteJSON(obj.SubMessage)
	if err != nil {
		return err
	}

	time.Sleep(writeWait)

	return nil
}
