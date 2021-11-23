package hendlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

const (
	writeWait  = 2 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type GetterCan interface {
	GetCandles(ws *websocket.Conn, wg *sync.WaitGroup, ctx context.Context, options Options) (chan Candles, error)
}

type Unsubscriber interface {
	Unsubscribe(ws *websocket.Conn) error
}

type WSMsg struct {
	Event   string   `json:"event"`
	Feed    string   `json:"feed"`
	Tickets []string `json:"product_ids"`
}

type Connection struct {
	SubMessage WSMsg
	Candle     Candles
	errHandler MyErrors

	ws *websocket.Conn
}

func (obj *Connection) candleStream(wg *sync.WaitGroup, ctx context.Context) (chan Candles, error) {
	canChan := make(chan Candles)

	wg.Add(1)
	go func() {
		var (
			prev int64
		)
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				close(canChan)
				return
			default:
				var can Candles
				err := obj.ws.ReadJSON(&can)
				if err != nil {
					obj.errHandler.WSReadMsgErr(errors.New("In cansleStream" + err.Error()))
					close(canChan)
					return
				}
				if prev != can.Candle.Time && can.Candle.Volume != 0 {
					canChan <- can
					prev = can.Candle.Time
				}
			}
		}
	}()

	return canChan, nil
}

func (obj Connection) GetCandles(ws *websocket.Conn, wg *sync.WaitGroup, ctx context.Context, options Options) (chan Candles, error) {
	obj.SubMessage.Event = "subscribe"
	obj.SubMessage.Feed = "candles_trade_" + string(options.canPer)
	obj.SubMessage.Tickets = options.ticket
	obj.ws = ws

	msg, err := json.Marshal(obj.SubMessage)
	if err != nil {
		err = obj.errHandler.MarshalErr(err)
		return nil, err
	}

	fmt.Println(string(msg))

	err = ws.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return nil, err
	}

	_, data, err := ws.ReadMessage()
	if err != nil {
		obj.errHandler.WSReadMsgErr(err)
	}

	jsonData := make(map[string]interface{})

	_ = json.Unmarshal(data, &jsonData)

	fmt.Println(jsonData)

	if jsonData["event"] != "subscribed" {
		return nil, SubErr
	}

	return obj.candleStream(wg, ctx)
}

func (obj Connection) Unsubscribe(ws *websocket.Conn) error {
	obj.SubMessage.Event = "unsubscribe"
	err := ws.WriteJSON(obj.SubMessage)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(writeWait)
	<-ticker.C
	ticker.Stop()

	return nil
}
